package plugin

import (
	"context"
	"seall/internal/api/mediaapi"
	"seall/internal/extension"
	"seall/internal/library/anime"

	"github.com/dop251/goja"
	"github.com/rs/zerolog"
)

type Simkl struct {
	ctx    *AppContextImpl
	ext    *extension.Extension
	logger *zerolog.Logger
}

// BindSimkl binds the simkl API to the Goja runtime.
// Permissions need to be checked by the caller.
// Permissions needed: simkl
func (a *AppContextImpl) BindSimkl(vm *goja.Runtime, logger *zerolog.Logger, ext *extension.Extension) {
	al := &Simkl{
		ctx:    a,
		ext:    ext,
		logger: new(logger.With().Str("id", ext.ID).Logger()),
	}
	simklObj := vm.NewObject()
	_ = simklObj.Set("refreshAnimeCollection", al.RefreshAnimeCollection)
	_ = simklObj.Set("refreshMangaCollection", al.RefreshMangaCollection)

	// Bind simkl platform
	mediaPlatformRef, ok := a.mediaPlatformRef.Get()
	if ok {
		_ = simklObj.Set("updateEntry", func(mediaID int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput) error {
			return mediaPlatformRef.Get().UpdateEntry(context.Background(), mediaID, status, scoreRaw, progress, startedAt, completedAt)
		})
		_ = simklObj.Set("updateEntryProgress", func(mediaID int, progress int, totalEpisodes *int) error {
			return mediaPlatformRef.Get().UpdateEntryProgress(context.Background(), mediaID, progress, totalEpisodes)
		})
		_ = simklObj.Set("updateEntryRepeat", func(mediaID int, repeat int) error {
			return mediaPlatformRef.Get().UpdateEntryRepeat(context.Background(), mediaID, repeat)
		})
		_ = simklObj.Set("deleteEntry", func(mediaID int, entryId int) error {
			return mediaPlatformRef.Get().DeleteEntry(context.Background(), mediaID, entryId)
		})
		_ = simklObj.Set("GetMediaCollection", func(bypassCache bool) (*mediaapi.AnimeCollection, error) {
			return mediaPlatformRef.Get().GetMediaCollection(context.Background(), bypassCache)
		})
		_ = simklObj.Set("GetRawMediaCollection", func(bypassCache bool) (*mediaapi.AnimeCollection, error) {
			return mediaPlatformRef.Get().GetRawMediaCollection(context.Background(), bypassCache)
		})
		_ = simklObj.Set("getMangaCollection", func(bypassCache bool) (*mediaapi.MangaCollection, error) {
			return mediaPlatformRef.Get().GetMangaCollection(context.Background(), bypassCache)
		})
		_ = simklObj.Set("getRawMangaCollection", func(bypassCache bool) (*mediaapi.MangaCollection, error) {
			return mediaPlatformRef.Get().GetRawMangaCollection(context.Background(), bypassCache)
		})
		_ = simklObj.Set("getAnime", func(mediaID int) (*mediaapi.BaseAnime, error) {
			return mediaPlatformRef.Get().GetAnime(context.Background(), mediaID)
		})
		_ = simklObj.Set("getManga", func(mediaID int) (*mediaapi.BaseManga, error) {
			return mediaPlatformRef.Get().GetManga(context.Background(), mediaID)
		})
		_ = simklObj.Set("getAnimeDetails", func(mediaID int) (*mediaapi.AnimeDetailsById_Media, error) {
			return mediaPlatformRef.Get().GetAnimeDetails(context.Background(), mediaID)
		})
		_ = simklObj.Set("getMangaDetails", func(mediaID int) (*mediaapi.MangaDetailsById_Media, error) {
			return mediaPlatformRef.Get().GetMangaDetails(context.Background(), mediaID)
		})
		_ = simklObj.Set("GetMediaCollectionWithRelations", func() (*mediaapi.AnimeCollectionWithRelations, error) {
			return mediaPlatformRef.Get().GetMediaCollectionWithRelations(context.Background())
		})
		_ = simklObj.Set("addMediaToCollection", func(mIds []int) error {
			return mediaPlatformRef.Get().AddMediaToCollection(context.Background(), mIds)
		})
		_ = simklObj.Set("getStudioDetails", func(studioID int) (*mediaapi.StudioDetails, error) {
			return mediaPlatformRef.Get().GetStudioDetails(context.Background(), studioID)
		})
		_ = simklObj.Set("listAnime", func(page *int, search *string, perPage *int, sort []*mediaapi.MediaSort, status []*mediaapi.MediaStatus, genres []*string, tags []*string, averageScoreGreater *int, season *mediaapi.MediaSeason, seasonYear *int, format *mediaapi.MediaFormat, isAdult *bool) (*mediaapi.ListAnime, error) {
			return mediaPlatformRef.Get().GetMediaApiClient().ListAnime(context.Background(), page, search, perPage, sort, status, genres, tags, averageScoreGreater, season, seasonYear, format, isAdult)
		})
		_ = simklObj.Set("listManga", func(page *int, search *string, perPage *int, sort []*mediaapi.MediaSort, status []*mediaapi.MediaStatus, genres []*string, tags []*string, averageScoreGreater *int, startDateGreater *string, startDateLesser *string, format *mediaapi.MediaFormat, countryOfOrigin *string, isAdult *bool) (*mediaapi.ListManga, error) {
			return mediaPlatformRef.Get().GetMediaApiClient().ListManga(context.Background(), page, search, perPage, sort, status, genres, tags, averageScoreGreater, startDateGreater, startDateLesser, format, countryOfOrigin, isAdult)
		})
		_ = simklObj.Set("listRecentAnime", func(page *int, perPage *int, airingAtGreater *int, airingAtLesser *int, notYetAired *bool) (*mediaapi.ListRecentAnime, error) {
			return mediaPlatformRef.Get().GetMediaApiClient().ListRecentAnime(context.Background(), page, perPage, airingAtGreater, airingAtLesser, notYetAired)
		})
		_ = simklObj.Set("clearCache", func() {
			mediaPlatformRef.Get().ClearCache()
			anime.ClearEpisodeCollectionCache()
			anime.ClearMissingEpisodesCache()
			anime.ClearScheduleCache()
		})
		_ = simklObj.Set("customQuery", func(body map[string]interface{}, token string) (interface{}, error) {
			return mediaapi.CustomQuery(body, a.logger, token)
		})

	}

	_ = vm.Set("$simkl", simklObj)
}

func (a *Simkl) RefreshAnimeCollection() {
	a.logger.Trace().Msg("plugin: Refreshing anime collection")
	onRefreshMediaCollection, ok := a.ctx.onRefreshMediaCollection.Get()
	if !ok {
		return
	}

	onRefreshMediaCollection()
}

func (a *Simkl) RefreshMangaCollection() {
	a.logger.Trace().Msg("plugin: Refreshing manga collection")
	onRefreshMangaCollection, ok := a.ctx.onRefreshMangaCollection.Get()
	if !ok {
		return
	}

	onRefreshMangaCollection()
}
