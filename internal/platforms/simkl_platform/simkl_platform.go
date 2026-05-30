package simkl_platform

import (
	"context"
	"errors"
	"net/url"
	"seall/internal/api/mediaapi"
	"seall/internal/api/simkl"
	"seall/internal/customsource"
	"seall/internal/database/db"
	"seall/internal/extension"
	"seall/internal/platforms/platform"
	"seall/internal/platforms/shared_platform"
	"seall/internal/util"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

type SimklPlatform struct {
	logger         *zerolog.Logger
	username       string
	simklClientRef *util.Ref[*simkl.Client]
	mediaClientRef *util.Ref[mediaapi.MediaApiClient]
	helper         *shared_platform.PlatformHelper
	db             *db.Database

	mu                 sync.RWMutex
	animeCollection    *mediaapi.AnimeCollection
	rawAnimeCollection *mediaapi.AnimeCollection
	mangaCollection    *mediaapi.MangaCollection
	rawMangaCollection *mediaapi.MangaCollection
	mediaKinds         map[int]simkl.MediaType
}

func NewSimklPlatform(simklClientRef *util.Ref[*simkl.Client], mediaClientRef *util.Ref[mediaapi.MediaApiClient], extensionBankRef *util.Ref[*extension.UnifiedBank], logger *zerolog.Logger, database *db.Database) platform.Platform {
	return &SimklPlatform{
		logger:         logger,
		simklClientRef: simklClientRef,
		mediaClientRef: mediaClientRef,
		helper:         shared_platform.NewPlatformHelper(extensionBankRef, database, logger),
		db:             database,
		mediaKinds:     make(map[int]simkl.MediaType),
	}
}

func (sp *SimklPlatform) SetUsername(username string) {
	sp.username = username
}

func (sp *SimklPlatform) ClearCache() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.animeCollection = nil
	sp.rawAnimeCollection = nil
	sp.mangaCollection = nil
	sp.rawMangaCollection = nil
	sp.mediaKinds = make(map[int]simkl.MediaType)
	if sp.helper != nil {
		sp.helper.ClearCache()
	}
}

func (sp *SimklPlatform) Close() {
	if sp.helper != nil {
		sp.helper.Close()
	}
}

func (sp *SimklPlatform) UpdateEntry(ctx context.Context, mediaID int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput) error {
	if customsource.IsExtensionId(mediaID) {
		if handled, err := sp.helper.HandleCustomSourceUpdateEntry(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt); handled {
			return err
		}
	}

	kind := sp.kindForMediaID(mediaID)
	media := simkl.ToSimklMedia(kind, mediaID)
	media.To = simkl.FromSIMKLStatus(status)
	media.Rating = simkl.ScoreToSIMKL(scoreRaw)
	req := simkl.AddItemsRequest{}
	if kind == simkl.MediaTypeMovies {
		req.Movies = []simkl.StandardMedia{media}
	} else {
		req.Shows = []simkl.StandardMedia{media}
	}
	_, err := sp.simklClient().AddToList(ctx, req)
	return err
}

func (sp *SimklPlatform) UpdateEntryProgress(ctx context.Context, mediaID int, progress int, totalEpisodes *int) error {
	status := mediaapi.MediaListStatusCurrent
	if totalEpisodes != nil && *totalEpisodes > 0 && progress >= *totalEpisodes {
		status = mediaapi.MediaListStatusCompleted
	}
	return sp.UpdateEntry(ctx, mediaID, &status, nil, &progress, nil, nil)
}

func (sp *SimklPlatform) UpdateEntryRepeat(ctx context.Context, mediaID int, repeat int) error {
	return nil
}

func (sp *SimklPlatform) DeleteEntry(ctx context.Context, mediaID int, entryID int) error {
	if customsource.IsExtensionId(mediaID) {
		if handled, err := sp.helper.HandleCustomSourceDeleteEntry(ctx, mediaID, entryID); handled {
			return err
		}
	}

	kind := sp.kindForMediaID(mediaID)
	media := simkl.ToSimklMedia(kind, mediaID)
	req := simkl.AddItemsRequest{}
	if kind == simkl.MediaTypeMovies {
		req.Movies = []simkl.StandardMedia{media}
	} else {
		req.Shows = []simkl.StandardMedia{media}
	}
	return sp.simklClient().RemoveItems(ctx, req)
}

func (sp *SimklPlatform) GetAnime(ctx context.Context, mediaID int) (*mediaapi.BaseAnime, error) {
	if cachedAnime, ok := sp.helper.GetCachedBaseAnime(mediaID); ok {
		return sp.helper.TriggerGetAnimeEvent(cachedAnime)
	}

	if media, isCustom, err := sp.helper.HandleCustomSourceAnime(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		sp.helper.SetCachedBaseAnime(mediaID, media)
		return sp.helper.TriggerGetAnimeEvent(media)
	}

	kind := sp.kindForMediaID(mediaID)
	if cached, ok := simkl.CachedDiscoveryMedia(mediaID); ok {
		ret := simkl.ToBaseAnime(cached.Kind, &cached.Media)
		sp.helper.SetCachedBaseAnime(mediaID, ret)
		sp.setMediaKind(mediaID, cached.Kind)
		return sp.helper.TriggerGetAnimeEvent(ret)
	}
	media, err := sp.fetchMediaDetails(ctx, kind, mediaID)
	if err != nil {
		return nil, err
	}
	ret := simkl.ToBaseAnime(kind, media)
	sp.helper.SetCachedBaseAnime(mediaID, ret)
	return sp.helper.TriggerGetAnimeEvent(ret)
}

func (sp *SimklPlatform) GetAnimeByMalID(ctx context.Context, malID int) (*mediaapi.BaseAnime, error) {
	return nil, errors.New("simkl: lookup by MAL id is not implemented yet")
}

func (sp *SimklPlatform) GetAnimeWithRelations(ctx context.Context, mediaID int) (*mediaapi.CompleteAnime, error) {
	if cachedAnime, ok := sp.helper.GetCachedCompleteAnime(mediaID); ok {
		return cachedAnime, nil
	}
	if media, isCustom, err := sp.helper.HandleCustomSourceAnimeWithRelations(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		sp.helper.SetCachedCompleteAnime(mediaID, media)
		return media, nil
	}

	kind := sp.kindForMediaID(mediaID)
	if cached, ok := simkl.CachedDiscoveryMedia(mediaID); ok {
		ret := simkl.ToCompleteAnime(cached.Kind, &cached.Media)
		sp.helper.SetCachedCompleteAnime(mediaID, ret)
		sp.setMediaKind(mediaID, cached.Kind)
		return ret, nil
	}
	media, err := sp.fetchMediaDetails(ctx, kind, mediaID)
	if err != nil {
		return nil, err
	}
	ret := simkl.ToCompleteAnime(kind, media)
	sp.helper.SetCachedCompleteAnime(mediaID, ret)
	return ret, nil
}

func (sp *SimklPlatform) GetAnimeDetails(ctx context.Context, mediaID int) (*mediaapi.AnimeDetailsById_Media, error) {
	if media, isCustom, err := sp.helper.HandleCustomSourceAnimeDetails(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		return sp.helper.TriggerGetAnimeDetailsEvent(media)
	}

	kind := sp.kindForMediaID(mediaID)
	if cached, ok := simkl.CachedDiscoveryMedia(mediaID); ok && len(cached.Media.Genres) > 0 {
		ret := simkl.ToAnimeDetails(cached.Kind, &cached.Media)
		sp.setMediaKind(mediaID, cached.Kind)
		return sp.helper.TriggerGetAnimeDetailsEvent(ret)
	}
	media, err := sp.fetchMediaDetails(ctx, kind, mediaID)
	if err != nil {
		return nil, err
	}
	ret := simkl.ToAnimeDetails(kind, media)
	return sp.helper.TriggerGetAnimeDetailsEvent(ret)
}

func (sp *SimklPlatform) GetManga(ctx context.Context, mediaID int) (*mediaapi.BaseManga, error) {
	return nil, errors.New("simkl: manga is not supported by SIMKL")
}

func (sp *SimklPlatform) GetMediaCollection(ctx context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error) {
	sp.mu.RLock()
	if !bypassCache && sp.animeCollection != nil {
		ret := sp.animeCollection
		sp.mu.RUnlock()
		return ret, nil
	}
	sp.mu.RUnlock()

	if err := sp.refreshAnimeCollection(ctx); err != nil {
		return nil, err
	}

	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.animeCollection, nil
}

func (sp *SimklPlatform) GetRawMediaCollection(ctx context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error) {
	sp.mu.RLock()
	if !bypassCache && sp.rawAnimeCollection != nil {
		ret := sp.rawAnimeCollection
		sp.mu.RUnlock()
		return ret, nil
	}
	sp.mu.RUnlock()

	if err := sp.refreshAnimeCollection(ctx); err != nil {
		return nil, err
	}

	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.rawAnimeCollection, nil
}

func (sp *SimklPlatform) RefreshAnimeCollection(ctx context.Context) (*mediaapi.AnimeCollection, error) {
	if err := sp.refreshAnimeCollection(ctx); err != nil {
		return nil, err
	}
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.animeCollection, nil
}

func (sp *SimklPlatform) refreshAnimeCollection(ctx context.Context) error {
	if sp.logger != nil {
		sp.logger.Debug().Msg("simkl platform: fetching all media collection")
	}

	// SIMKL asks clients to check activities before fetching all-items.
	_, _ = sp.simklClient().Activities(ctx)

	q := url.Values{}
	q.Set("extended", "full")
	items, err := sp.simklClient().AllItems(ctx, simkl.MediaTypeAll, "", q)
	if err != nil {
		return err
	}

	collection := simkl.ToSIMKLAnimeCollection(items)
	sp.helper.MergeCustomSourceAnimeEntries(collection)
	raw := cloneAnimeCollection(collection)
	collection.MediaListCollection.Lists = sp.helper.FilterOutCustomAnimeLists(collection.MediaListCollection.Lists)

	kinds := make(map[int]simkl.MediaType)
	for _, item := range items.Movies {
		if media := item.Media(); media != nil {
			kinds[media.IDs.PrimarySimklID()] = simkl.MediaTypeMovies
		}
	}
	for _, item := range items.Shows {
		if media := item.Media(); media != nil {
			kinds[media.IDs.PrimarySimklID()] = simkl.MediaTypeShows
		}
	}
	for _, item := range items.Anime {
		if media := item.Media(); media != nil {
			kinds[media.IDs.PrimarySimklID()] = simkl.MediaTypeAnime
		}
	}

	sp.mu.Lock()
	sp.animeCollection = collection
	sp.rawAnimeCollection = raw
	sp.mediaKinds = kinds
	sp.mu.Unlock()
	return nil
}

func (sp *SimklPlatform) GetMediaCollectionWithRelations(ctx context.Context) (*mediaapi.AnimeCollectionWithRelations, error) {
	return &mediaapi.AnimeCollectionWithRelations{
		MediaListCollection: &mediaapi.AnimeCollectionWithRelations_MediaListCollection{
			Lists: []*mediaapi.AnimeCollectionWithRelations_MediaListCollection_Lists{},
		},
	}, nil
}

func (sp *SimklPlatform) GetMangaCollection(ctx context.Context, bypassCache bool) (*mediaapi.MangaCollection, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	if sp.mangaCollection == nil {
		sp.mangaCollection = emptyMangaCollection()
	}
	return sp.mangaCollection, nil
}

func (sp *SimklPlatform) GetRawMangaCollection(ctx context.Context, bypassCache bool) (*mediaapi.MangaCollection, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	if sp.rawMangaCollection == nil {
		sp.rawMangaCollection = emptyMangaCollection()
	}
	return sp.rawMangaCollection, nil
}

func (sp *SimklPlatform) RefreshMangaCollection(ctx context.Context) (*mediaapi.MangaCollection, error) {
	return sp.GetMangaCollection(ctx, true)
}

func (sp *SimklPlatform) GetMangaDetails(ctx context.Context, mediaID int) (*mediaapi.MangaDetailsById_Media, error) {
	return nil, errors.New("simkl: manga is not supported by SIMKL")
}

func (sp *SimklPlatform) AddMediaToCollection(ctx context.Context, ids []int) error {
	for _, id := range ids {
		status := mediaapi.MediaListStatusPlanning
		if err := sp.UpdateEntry(ctx, id, &status, nil, nil, nil, nil); err != nil {
			return err
		}
	}
	return nil
}

func (sp *SimklPlatform) GetStudioDetails(ctx context.Context, studioID int) (*mediaapi.StudioDetails, error) {
	return nil, errors.New("simkl: studio details are not supported by SIMKL")
}

func (sp *SimklPlatform) GetMediaApiClient() mediaapi.MediaApiClient {
	if sp.mediaClientRef == nil {
		return nil
	}
	return sp.mediaClientRef.Get()
}

func (sp *SimklPlatform) GetViewerStats(ctx context.Context) (*mediaapi.ViewerStats, error) {
	collection, err := sp.GetRawMediaCollection(ctx, false)
	if err != nil {
		return nil, err
	}
	return simklViewerStatsFromCollection(collection), nil
}

func (sp *SimklPlatform) GetAnimeAiringSchedule(ctx context.Context) (*mediaapi.AnimeAiringSchedule, error) {
	return &mediaapi.AnimeAiringSchedule{}, nil
}

func (sp *SimklPlatform) simklClient() *simkl.Client {
	return sp.simklClientRef.Get()
}

func (sp *SimklPlatform) kindForMediaID(mediaID int) simkl.MediaType {
	sp.mu.RLock()
	kind, ok := sp.mediaKinds[mediaID]
	sp.mu.RUnlock()
	if ok {
		return kind
	}
	return simkl.MediaTypeAll
}

func (sp *SimklPlatform) setMediaKind(mediaID int, kind simkl.MediaType) {
	if kind == simkl.MediaTypeAll || kind == "" {
		return
	}
	sp.mu.Lock()
	sp.mediaKinds[mediaID] = kind
	sp.mu.Unlock()
}

func (sp *SimklPlatform) fetchMediaDetails(ctx context.Context, kind simkl.MediaType, mediaID int) (*simkl.StandardMedia, error) {
	kinds := simklMediaKindCandidates(kind)
	var lastErr error
	for _, k := range kinds {
		media, err := sp.simklClient().MediaDetails(ctx, k, strconv.Itoa(mediaID), "full")
		if err == nil && media != nil {
			if media.IDs.Simkl == 0 {
				media.IDs.Simkl = mediaID
			}
			if media.Type == "" {
				media.Type = string(k)
			}
			sp.setMediaKind(mediaID, simkl.KindFromStandardMedia(k, media))
			return media, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func simklMediaKindCandidates(kind simkl.MediaType) []simkl.MediaType {
	ret := make([]simkl.MediaType, 0, 3)
	add := func(candidate simkl.MediaType) {
		if candidate != simkl.MediaTypeMovies && candidate != simkl.MediaTypeShows && candidate != simkl.MediaTypeAnime {
			return
		}
		for _, existing := range ret {
			if existing == candidate {
				return
			}
		}
		ret = append(ret, candidate)
	}

	add(kind)
	add(simkl.MediaTypeMovies)
	add(simkl.MediaTypeShows)
	add(simkl.MediaTypeAnime)
	return ret
}

func cloneAnimeCollection(collection *mediaapi.AnimeCollection) *mediaapi.AnimeCollection {
	if collection == nil || collection.MediaListCollection == nil {
		return &mediaapi.AnimeCollection{MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{}}
	}
	listsCopy := make([]*mediaapi.AnimeCollection_MediaListCollection_Lists, len(collection.MediaListCollection.Lists))
	copy(listsCopy, collection.MediaListCollection.Lists)
	return &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{Lists: listsCopy},
	}
}

func emptyMangaCollection() *mediaapi.MangaCollection {
	return &mediaapi.MangaCollection{
		MediaListCollection: &mediaapi.MangaCollection_MediaListCollection{
			Lists: []*mediaapi.MangaCollection_MediaListCollection_Lists{},
		},
	}
}
