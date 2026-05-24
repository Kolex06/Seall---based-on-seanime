package simulated_platform

import (
	"context"
	"encoding/json"
	"errors"
	"seall/internal/api/mediaapi"
	"seall/internal/customsource"
	"seall/internal/database/db"
	"seall/internal/extension"
	"seall/internal/hook"
	"seall/internal/local"
	"seall/internal/platforms/platform"
	"seall/internal/platforms/shared_platform"
	"seall/internal/util"
	"seall/internal/util/limiter"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	// ErrMediaNotFound means the media wasn't found in the local collection
	ErrMediaNotFound = errors.New("media not found")
)

// SimulatedPlatform used when the user is not authenticated to SIMKL.
// It acts as a dummy account using simulated collections stored locally.
type SimulatedPlatform struct {
	logger       *zerolog.Logger
	localManager local.Manager
	client       mediaapi.MediaApiClient // should only receive an unauthenticated client

	// Cache for collections
	animeCollection                *mediaapi.AnimeCollection
	mangaCollection                *mediaapi.MangaCollection
	mu                             sync.RWMutex
	collectionMu                   sync.RWMutex // used to protect access to collections
	lastAnimeCollectionRefetchTime time.Time    // used to prevent refetching too many times
	lastMangaCollectionRefetchTime time.Time    // used to prevent refetching too many times
	helper                         *shared_platform.PlatformHelper
	db                             *db.Database
	refreshRateLimit               *limiter.Limiter
	refreshAnimeMetadataCancelFunc context.CancelFunc
	refreshMangaMetadataCancelFunc context.CancelFunc
}

func NewSimulatedPlatform(localManager local.Manager, client *util.Ref[mediaapi.MediaApiClient], extensionBankRef *util.Ref[*extension.UnifiedBank], logger *zerolog.Logger, db *db.Database) (platform.Platform, error) {
	sp := &SimulatedPlatform{
		logger:           logger,
		localManager:     localManager,
		client:           shared_platform.NewCacheLayer(client),
		refreshRateLimit: limiter.NewLimiter(2*time.Second, 1),
		helper:           shared_platform.NewPlatformHelper(extensionBankRef, db, logger),
		db:               db,
	}

	return sp, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Implementation
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (sp *SimulatedPlatform) SetUsername(username string) {
	// no-op
}

func (sp *SimulatedPlatform) Close() {
	sp.helper.Close()
}

func (sp *SimulatedPlatform) ClearCache() {
	sp.helper.ClearCache()
}

// UpdateEntry updates the entry for the given media ID.
// If the entry doesn't exist, it will be added automatically after determining the media type.
func (sp *SimulatedPlatform) UpdateEntry(ctx context.Context, mediaID int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput) error {
	sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Updating entry")

	return sp.helper.TriggerUpdateEntryHooks(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt, func(event *platform.PreUpdateEntryEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := sp.helper.HandleCustomSourceUpdateEntry(ctx, mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt); handled {
			return err
		}

		sp.mu.Lock()
		defer sp.mu.Unlock()

		// Try anime first
		animeWrapper := sp.GetMediaCollectionWrapper()
		if _, err := animeWrapper.FindEntry(mediaID); err == nil {
			return animeWrapper.UpdateEntry(mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt)
		}

		// Try manga
		mangaWrapper := sp.GetMangaCollectionWrapper()
		if _, err := mangaWrapper.FindEntry(mediaID); err == nil {
			return mangaWrapper.UpdateEntry(mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt)
		}

		// Entry doesn't exist, determine media type and add it
		defaultStatus := mediaapi.MediaListStatusPlanning
		if event.Status != nil {
			defaultStatus = *event.Status
		}

		// Try to fetch as anime first
		if _, err := sp.client.BaseAnimeByID(ctx, &mediaID); err == nil {
			// It's an anime, add it to anime collection
			sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Adding new anime entry")
			if err := animeWrapper.AddEntry(mediaID, defaultStatus); err != nil {
				return err
			}
			// Update with provided values if there are additional updates needed
			if event.Status != &defaultStatus || event.ScoreRaw != nil || event.Progress != nil || event.StartedAt != nil || event.CompletedAt != nil {
				return animeWrapper.UpdateEntry(mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt)
			}
			return nil
		}

		// Try to fetch as manga
		if _, err := sp.client.BaseMangaByID(ctx, &mediaID); err == nil {
			// It's a manga, add it to manga collection
			sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Adding new manga entry")
			if err := mangaWrapper.AddEntry(mediaID, defaultStatus); err != nil {
				return err
			}
			// Update with provided values if there are additional updates needed
			if event.Status != &defaultStatus || event.ScoreRaw != nil || event.Progress != nil || event.StartedAt != nil || event.CompletedAt != nil {
				return mangaWrapper.UpdateEntry(mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt)
			}
			return nil
		}

		// Media not found in either anime or manga
		return errors.New("media not found on SIMKL")
	})
}

func (sp *SimulatedPlatform) UpdateEntryProgress(ctx context.Context, mediaID int, progress int, totalEpisodes *int) error {
	sp.logger.Trace().Int("mediaID", mediaID).Int("progress", progress).Msg("simulated platform: Updating entry progress")

	return sp.helper.TriggerUpdateEntryProgressHooks(ctx, mediaID, progress, totalEpisodes, func(event *platform.PreUpdateEntryProgressEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := sp.helper.HandleCustomSourceUpdateEntryProgress(ctx, mediaID, *event.Progress, event.TotalCount); handled {
			return err
		}

		sp.mu.Lock()
		defer sp.mu.Unlock()

		status := mediaapi.MediaListStatusCurrent
		if event.TotalCount != nil && *event.Progress >= *event.TotalCount {
			status = mediaapi.MediaListStatusCompleted
			*event.Status = status
		}

		// Try anime first
		animeWrapper := sp.GetMediaCollectionWrapper()
		if _, err := animeWrapper.FindEntry(mediaID); err == nil {
			return animeWrapper.UpdateEntryProgress(mediaID, *event.Progress, event.TotalCount)
		}

		// Try manga
		mangaWrapper := sp.GetMangaCollectionWrapper()
		if _, err := mangaWrapper.FindEntry(mediaID); err == nil {
			return mangaWrapper.UpdateEntryProgress(mediaID, *event.Progress, event.TotalCount)
		}

		// Entry doesn't exist, determine media type and add it
		// Try to fetch as anime first
		if _, err := sp.client.BaseAnimeByID(ctx, &mediaID); err == nil {
			// It's an anime, add it to anime collection
			sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Adding new anime entry for progress update")
			if err := animeWrapper.AddEntry(mediaID, status); err != nil {
				return err
			}
			return animeWrapper.UpdateEntryProgress(mediaID, *event.Progress, event.TotalCount)
		}

		// Try to fetch as manga
		if _, err := sp.client.BaseMangaByID(ctx, &mediaID); err == nil {
			// It's a manga, add it to manga collection
			sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Adding new manga entry for progress update")
			if err := mangaWrapper.AddEntry(mediaID, status); err != nil {
				return err
			}
			return mangaWrapper.UpdateEntryProgress(mediaID, *event.Progress, event.TotalCount)
		}

		// Media not found in either anime or manga
		return errors.New("media not found on SIMKL")
	})
}

func (sp *SimulatedPlatform) UpdateEntryRepeat(ctx context.Context, mediaID int, repeat int) error {
	sp.logger.Trace().Int("mediaID", mediaID).Int("repeat", repeat).Msg("simulated platform: Updating entry repeat")

	return sp.helper.TriggerUpdateEntryRepeatHooks(ctx, mediaID, repeat, func(event *platform.PreUpdateEntryRepeatEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := sp.helper.HandleCustomSourceUpdateEntryRepeat(ctx, mediaID, *event.Repeat); handled {
			return err
		}

		sp.mu.Lock()
		defer sp.mu.Unlock()

		// Try anime first
		wrapper := sp.GetMediaCollectionWrapper()
		if entry, err := wrapper.FindEntry(mediaID); err == nil {
			if animeEntry, ok := entry.(*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries); ok {
				animeEntry.Repeat = event.Repeat
				sp.localManager.SaveSimulatedAnimeCollection(sp.animeCollection)
				return nil
			}
		}

		// Try manga
		wrapper = sp.GetMangaCollectionWrapper()
		if entry, err := wrapper.FindEntry(mediaID); err == nil {
			if mangaEntry, ok := entry.(*mediaapi.MangaCollection_MediaListCollection_Lists_Entries); ok {
				mangaEntry.Repeat = event.Repeat
				sp.localManager.SaveSimulatedMangaCollection(sp.mangaCollection)
				return nil
			}
		}

		return ErrMediaNotFound
	})
}

func (sp *SimulatedPlatform) DeleteEntry(ctx context.Context, mediaId, entryId int) error {
	sp.logger.Trace().Int("entryId", entryId).Int("mediaId", mediaId).Msg("simulated platform: Deleting entry")

	return sp.helper.TriggerDeleteEntryHooks(ctx, mediaId, entryId, func(event *platform.PreDeleteEntryEvent) error {
		if handled, err := sp.helper.HandleCustomSourceDeleteEntry(ctx, *event.MediaID, *event.EntryID); handled {
			return err
		}

		sp.mu.Lock()
		defer sp.mu.Unlock()

		// Try anime first
		wrapper := sp.GetMediaCollectionWrapper()
		if _, err := wrapper.FindEntry(*event.EntryID, true); err == nil {
			return wrapper.DeleteEntry(*event.EntryID, true)
		}

		// Try manga
		wrapper = sp.GetMangaCollectionWrapper()
		if _, err := wrapper.FindEntry(*event.EntryID, true); err == nil {
			return wrapper.DeleteEntry(*event.EntryID, true)
		}

		return ErrMediaNotFound
	})
}

func (sp *SimulatedPlatform) GetAnime(ctx context.Context, mediaID int) (*mediaapi.BaseAnime, error) {
	sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Getting anime")

	if cachedAnime, ok := sp.helper.GetCachedBaseAnime(mediaID); ok {
		sp.logger.Trace().Msg("simulated platform: Returning anime from cache")
		return sp.helper.TriggerGetAnimeEvent(cachedAnime)
	}

	// Check if this is a custom source entry
	if media, isCustom, err := sp.helper.HandleCustomSourceAnime(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}

		triggeredMedia, err := sp.helper.TriggerGetAnimeEvent(media)
		if err != nil {
			return nil, err
		}

		sp.helper.SetCachedBaseAnime(mediaID, triggeredMedia)

		// Update media data in collection if it exists (simulated platform specific)
		sp.mu.Lock()
		wrapper := sp.GetMediaCollectionWrapper()
		if _, err := wrapper.FindEntry(mediaID); err == nil {
			_ = wrapper.UpdateMediaData(mediaID, triggeredMedia)
		}
		sp.mu.Unlock()

		return triggeredMedia, nil
	}

	// Get anime from simkl
	resp, err := sp.client.BaseAnimeByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}
	media := resp.GetMedia()

	triggeredMedia, err := sp.helper.TriggerGetAnimeEvent(media)
	if err != nil {
		return nil, err
	}

	sp.helper.SetCachedBaseAnime(mediaID, triggeredMedia)

	// Update media data in collection if it exists (simulated platform specific)
	sp.mu.Lock()
	wrapper := sp.GetMediaCollectionWrapper()
	if _, err := wrapper.FindEntry(mediaID); err == nil {
		_ = wrapper.UpdateMediaData(mediaID, triggeredMedia)
	}
	sp.mu.Unlock()

	return triggeredMedia, nil
}

func (sp *SimulatedPlatform) GetAnimeByMalID(ctx context.Context, malID int) (*mediaapi.BaseAnime, error) {
	sp.logger.Trace().Int("malID", malID).Msg("simulated platform: Getting anime by MAL ID")

	resp, err := sp.client.BaseAnimeByMalID(ctx, &malID)
	if err != nil {
		return nil, err
	}

	media := resp.GetMedia()
	triggeredMedia, err := sp.helper.TriggerGetAnimeEvent(media)
	if err != nil {
		return nil, err
	}

	// Update media data in collection if it exists (simulated platform specific)
	if triggeredMedia != nil {
		sp.mu.Lock()
		wrapper := sp.GetMediaCollectionWrapper()
		if _, err := wrapper.FindEntry(triggeredMedia.GetID()); err == nil {
			_ = wrapper.UpdateMediaData(triggeredMedia.GetID(), triggeredMedia)
		}
		sp.mu.Unlock()
	}

	return triggeredMedia, nil
}

func (sp *SimulatedPlatform) GetAnimeDetails(ctx context.Context, mediaID int) (*mediaapi.AnimeDetailsById_Media, error) {
	sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Getting anime details")

	// Check if this is a custom source entry
	if media, isCustom, err := sp.helper.HandleCustomSourceAnimeDetails(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		return sp.helper.TriggerGetAnimeDetailsEvent(media)
	}

	// Get from SIMKL
	resp, err := sp.client.AnimeDetailsByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}
	media := resp.GetMedia()

	return sp.helper.TriggerGetAnimeDetailsEvent(media)
}

func (sp *SimulatedPlatform) GetAnimeWithRelations(ctx context.Context, mediaID int) (*mediaapi.CompleteAnime, error) {
	sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Getting anime with relations")

	if cachedAnime, ok := sp.helper.GetCachedCompleteAnime(mediaID); ok {
		sp.logger.Trace().Msg("simulated platform: Cache HIT for anime with relations")
		return cachedAnime, nil
	}

	// Check if this is a custom source entry
	if media, isCustom, err := sp.helper.HandleCustomSourceAnimeWithRelations(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		sp.helper.SetCachedCompleteAnime(mediaID, media)
		return media, nil
	}

	// Get from SIMKL
	resp, err := sp.client.CompleteAnimeByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}
	media := resp.GetMedia()

	sp.helper.SetCachedCompleteAnime(mediaID, media)
	return media, nil
}

func (sp *SimulatedPlatform) GetManga(ctx context.Context, mediaID int) (*mediaapi.BaseManga, error) {
	sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Getting manga")

	if cachedManga, ok := sp.helper.GetCachedBaseManga(mediaID); ok {
		sp.logger.Trace().Msg("simulated platform: Returning manga from cache")
		return sp.helper.TriggerGetMangaEvent(cachedManga)
	}

	// Check if this is a custom source entry
	if media, isCustom, err := sp.helper.HandleCustomSourceManga(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}

		triggeredMedia, err := sp.helper.TriggerGetMangaEvent(media)
		if err != nil {
			return nil, err
		}

		sp.helper.SetCachedBaseManga(mediaID, triggeredMedia)

		// Update media data in collection if it exists (simulated platform specific)
		sp.mu.Lock()
		wrapper := sp.GetMangaCollectionWrapper()
		if _, err := wrapper.FindEntry(mediaID); err == nil {
			_ = wrapper.UpdateMediaData(mediaID, triggeredMedia)
		}
		sp.mu.Unlock()

		return triggeredMedia, nil
	}

	// Get manga from simkl
	resp, err := sp.client.BaseMangaByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}
	media := resp.GetMedia()

	triggeredMedia, err := sp.helper.TriggerGetMangaEvent(media)
	if err != nil {
		return nil, err
	}

	sp.helper.SetCachedBaseManga(mediaID, triggeredMedia)

	// Update media data in collection if it exists (simulated platform specific)
	sp.mu.Lock()
	wrapper := sp.GetMangaCollectionWrapper()
	if _, err := wrapper.FindEntry(mediaID); err == nil {
		_ = wrapper.UpdateMediaData(mediaID, triggeredMedia)
	}
	sp.mu.Unlock()

	return triggeredMedia, nil
}

func (sp *SimulatedPlatform) GetMangaDetails(ctx context.Context, mediaID int) (*mediaapi.MangaDetailsById_Media, error) {
	sp.logger.Trace().Int("mediaID", mediaID).Msg("simulated platform: Getting manga details")

	// Check if this is a custom source entry
	if media, isCustom, err := sp.helper.HandleCustomSourceMangaDetails(ctx, mediaID); isCustom {
		return media, err
	}

	// Get from SIMKL
	resp, err := sp.client.MangaDetailsByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}

	return resp.GetMedia(), nil
}

func (sp *SimulatedPlatform) GetMediaCollection(ctx context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error) {
	sp.logger.Trace().Bool("bypassCache", bypassCache).Msg("simulated platform: Getting anime collection")

	if !bypassCache && sp.animeCollection != nil {
		event := new(platform.GetCachedAnimeCollectionEvent)
		event.AnimeCollection = sp.animeCollection
		err := hook.GlobalHookManager.OnGetCachedAnimeCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.AnimeCollection, nil
	}

	if bypassCache {
		sp.invalidateAnimeCollectionCache()
	}

	collection, err := sp.getOrCreateAnimeCollection()
	if err != nil {
		return nil, err
	}

	// Merge custom source entries if available
	sp.helper.MergeCustomSourceAnimeEntries(collection)

	event := new(platform.GetMediaCollectionEvent)
	event.AnimeCollection = collection

	err = hook.GlobalHookManager.OnGetMediaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (sp *SimulatedPlatform) GetRawMediaCollection(ctx context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error) {
	sp.logger.Trace().Bool("bypassCache", bypassCache).Msg("simulated platform: Getting raw anime collection")

	if !bypassCache && sp.animeCollection != nil {
		event := new(platform.GetCachedRawAnimeCollectionEvent)
		event.AnimeCollection = sp.animeCollection
		err := hook.GlobalHookManager.OnGetCachedRawAnimeCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.AnimeCollection, nil
	}

	if bypassCache {
		sp.invalidateAnimeCollectionCache()
	}

	collection, err := sp.getOrCreateAnimeCollection()
	if err != nil {
		return nil, err
	}

	// Merge custom source entries if available
	sp.helper.MergeCustomSourceAnimeEntries(collection)

	event := new(platform.GetRawMediaCollectionEvent)
	event.AnimeCollection = collection

	err = hook.GlobalHookManager.OnGetRawMediaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (sp *SimulatedPlatform) RefreshAnimeCollection(ctx context.Context) (*mediaapi.AnimeCollection, error) {
	sp.logger.Trace().Msg("simulated platform: Refreshing anime collection")

	sp.invalidateAnimeCollectionCache()
	collection, err := sp.getOrCreateAnimeCollection()
	if err != nil {
		return nil, err
	}
	if err := sp.refreshAnimeCollectionMetadata(ctx, collection); err != nil {
		return nil, err
	}

	// Merge custom source entries if available
	sp.helper.MergeCustomSourceAnimeEntries(collection)

	event := new(platform.GetMediaCollectionEvent)
	event.AnimeCollection = collection

	err = hook.GlobalHookManager.OnGetMediaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	event2 := new(platform.GetRawMediaCollectionEvent)
	event2.AnimeCollection = collection

	err = hook.GlobalHookManager.OnGetRawMediaCollection().Trigger(event2)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

// GetMediaCollectionWithRelations returns the anime collection (without relations)
func (sp *SimulatedPlatform) GetMediaCollectionWithRelations(ctx context.Context) (*mediaapi.AnimeCollectionWithRelations, error) {
	sp.logger.Trace().Msg("simulated platform: Getting anime collection with relations")

	// Use JSON to convert the collection structs
	collection, err := sp.getOrCreateAnimeCollection()
	if err != nil {
		return nil, err
	}

	collectionWithRelations := &mediaapi.AnimeCollectionWithRelations{}

	marshaled, err := json.Marshal(collection)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaled, collectionWithRelations)
	if err != nil {
		return nil, err
	}

	// For simulated platform, the anime collection will not have relations
	return collectionWithRelations, nil
}

func (sp *SimulatedPlatform) GetMangaCollection(ctx context.Context, bypassCache bool) (*mediaapi.MangaCollection, error) {
	sp.logger.Trace().Bool("bypassCache", bypassCache).Msg("simulated platform: Getting manga collection")

	if !bypassCache && sp.mangaCollection != nil {
		event := new(platform.GetCachedMangaCollectionEvent)
		event.MangaCollection = sp.mangaCollection
		err := hook.GlobalHookManager.OnGetCachedMangaCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.MangaCollection, nil
	}

	if bypassCache {
		sp.invalidateMangaCollectionCache()
	}

	collection, err := sp.getOrCreateMangaCollection()
	if err != nil {
		return nil, err
	}

	// Merge custom source entries if available
	sp.helper.MergeCustomSourceMangaEntries(collection)

	event := new(platform.GetMangaCollectionEvent)
	event.MangaCollection = collection

	err = hook.GlobalHookManager.OnGetMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (sp *SimulatedPlatform) GetRawMangaCollection(ctx context.Context, bypassCache bool) (*mediaapi.MangaCollection, error) {
	sp.logger.Trace().Bool("bypassCache", bypassCache).Msg("simulated platform: Getting raw manga collection")

	if !bypassCache && sp.mangaCollection != nil {
		event := new(platform.GetCachedRawMangaCollectionEvent)
		event.MangaCollection = sp.mangaCollection
		err := hook.GlobalHookManager.OnGetCachedRawMangaCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.MangaCollection, nil
	}

	if bypassCache {
		sp.invalidateMangaCollectionCache()
	}

	collection, err := sp.getOrCreateMangaCollection()
	if err != nil {
		return nil, err
	}

	// Merge custom source entries if available
	sp.helper.MergeCustomSourceMangaEntries(collection)

	event := new(platform.GetRawMangaCollectionEvent)
	event.MangaCollection = collection

	err = hook.GlobalHookManager.OnGetRawMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (sp *SimulatedPlatform) RefreshMangaCollection(ctx context.Context) (*mediaapi.MangaCollection, error) {
	sp.logger.Trace().Msg("simulated platform: Refreshing manga collection")

	sp.invalidateMangaCollectionCache()
	collection, err := sp.getOrCreateMangaCollection()
	if err != nil {
		return nil, err
	}
	if err := sp.refreshMangaCollectionMetadata(ctx, collection); err != nil {
		return nil, err
	}

	// Merge custom source entries if available
	sp.helper.MergeCustomSourceMangaEntries(collection)

	event := new(platform.GetMangaCollectionEvent)
	event.MangaCollection = collection

	err = hook.GlobalHookManager.OnGetMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	event2 := new(platform.GetRawMangaCollectionEvent)
	event2.MangaCollection = collection

	err = hook.GlobalHookManager.OnGetRawMangaCollection().Trigger(event2)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (sp *SimulatedPlatform) AddMediaToCollection(ctx context.Context, mIds []int) error {
	sp.logger.Trace().Interface("mediaIDs", mIds).Msg("simulated platform: Adding media to collection")

	sp.mu.Lock()
	defer sp.mu.Unlock()

	// DEVNOTE: We assume it's anime for now since it's only been used for anime
	wrapper := sp.GetMediaCollectionWrapper()
	for _, mediaID := range mIds {
		// Try to add as anime first, if it fails, ignore
		_ = wrapper.AddEntry(mediaID, mediaapi.MediaListStatusPlanning)
	}

	return nil
}

func (sp *SimulatedPlatform) GetStudioDetails(ctx context.Context, studioID int) (*mediaapi.StudioDetails, error) {
	sp.logger.Trace().Int("studioID", studioID).Msg("simulated platform: Getting studio details")

	ret, err := sp.client.StudioDetails(ctx, &studioID)
	if err != nil {
		return nil, err
	}

	return sp.helper.TriggerGetStudioDetailsEvent(ret)
}

func (sp *SimulatedPlatform) GetMediaApiClient() mediaapi.MediaApiClient {
	return sp.client
}

func (sp *SimulatedPlatform) GetViewerStats(ctx context.Context) (*mediaapi.ViewerStats, error) {
	return nil, errors.New("use a real account to get stats")
}

func (sp *SimulatedPlatform) GetAnimeAiringSchedule(ctx context.Context) (*mediaapi.AnimeAiringSchedule, error) {
	collection, err := sp.GetMediaCollection(ctx, false)
	if err != nil {
		return nil, err
	}

	return sp.helper.BuildAnimeAiringSchedule(ctx, collection, sp.client)
}

func (sp *SimulatedPlatform) refreshAnimeCollectionMetadata(ctx context.Context, collection *mediaapi.AnimeCollection) error {
	if sp.refreshAnimeMetadataCancelFunc != nil {
		sp.refreshAnimeMetadataCancelFunc()
	}
	var fctx context.Context
	fctx, sp.refreshAnimeMetadataCancelFunc = context.WithCancel(ctx)
	defer sp.refreshAnimeMetadataCancelFunc()

	mediaIDs := collectRefreshableAnimeIDs(collection)
	if len(mediaIDs) == 0 {
		return nil
	}

	sp.logger.Trace().Int("count", len(mediaIDs)).Msg("simulated platform: Refreshing mutable anime metadata")
	wrapper := sp.GetMediaCollectionWrapper()

	for _, mediaID := range mediaIDs {
		if err := fctx.Err(); err != nil {
			return err
		}

		sp.refreshRateLimit.Wait()

		resp, err := sp.client.BaseAnimeByID(ctx, &mediaID)
		if err != nil {
			sp.logger.Warn().Err(err).Int("mediaID", mediaID).Msg("simulated platform: Failed to refresh anime metadata")
			continue
		}

		triggeredMedia, err := sp.helper.TriggerGetAnimeEvent(resp.GetMedia())
		if err != nil {
			sp.logger.Warn().Err(err).Int("mediaID", mediaID).Msg("simulated platform: Failed to process refreshed anime metadata")
			continue
		}
		if triggeredMedia == nil {
			continue
		}

		sp.helper.SetCachedBaseAnime(mediaID, triggeredMedia)

		sp.mu.Lock()
		err = wrapper.UpdateMediaData(mediaID, triggeredMedia)
		sp.mu.Unlock()
		if err != nil && !errors.Is(err, ErrMediaNotFound) {
			sp.logger.Warn().Err(err).Int("mediaID", mediaID).Msg("simulated platform: Failed to save refreshed anime metadata")
		}
	}

	return nil
}

func (sp *SimulatedPlatform) refreshMangaCollectionMetadata(ctx context.Context, collection *mediaapi.MangaCollection) error {
	if sp.refreshMangaMetadataCancelFunc != nil {
		sp.refreshMangaMetadataCancelFunc()
	}
	var fctx context.Context
	fctx, sp.refreshMangaMetadataCancelFunc = context.WithCancel(ctx)
	defer sp.refreshMangaMetadataCancelFunc()

	mediaIDs := collectRefreshableMangaIDs(collection)
	if len(mediaIDs) == 0 {
		return nil
	}

	sp.logger.Trace().Int("count", len(mediaIDs)).Msg("simulated platform: Refreshing mutable manga metadata")
	wrapper := sp.GetMangaCollectionWrapper()

	for _, mediaID := range mediaIDs {
		if err := fctx.Err(); err != nil {
			return err
		}

		sp.refreshRateLimit.Wait()

		resp, err := sp.client.BaseMangaByID(ctx, &mediaID)
		if err != nil {
			sp.logger.Warn().Err(err).Int("mediaID", mediaID).Msg("simulated platform: Failed to refresh manga metadata")
			continue
		}

		triggeredMedia, err := sp.helper.TriggerGetMangaEvent(resp.GetMedia())
		if err != nil {
			sp.logger.Warn().Err(err).Int("mediaID", mediaID).Msg("simulated platform: Failed to process refreshed manga metadata")
			continue
		}
		if triggeredMedia == nil {
			continue
		}

		sp.helper.SetCachedBaseManga(mediaID, triggeredMedia)

		sp.mu.Lock()
		err = wrapper.UpdateMediaData(mediaID, triggeredMedia)
		sp.mu.Unlock()
		if err != nil && !errors.Is(err, ErrMediaNotFound) {
			sp.logger.Warn().Err(err).Int("mediaID", mediaID).Msg("simulated platform: Failed to save refreshed manga metadata")
		}
	}

	return nil
}

func collectRefreshableAnimeIDs(collection *mediaapi.AnimeCollection) []int {
	if collection == nil || collection.GetMediaListCollection() == nil {
		return nil
	}

	ret := make([]int, 0)
	seen := make(map[int]struct{})
	for _, list := range collection.GetMediaListCollection().GetLists() {
		for _, entry := range list.GetEntries() {
			if entry == nil || entry.GetMedia() == nil {
				continue
			}

			mediaID := entry.GetMedia().GetID()
			if _, ok := seen[mediaID]; ok || !shouldRefreshSimulatedMedia(entry.GetStatus(), entry.GetMedia().GetStatus(), mediaID) {
				continue
			}

			seen[mediaID] = struct{}{}
			ret = append(ret, mediaID)
		}
	}

	return ret
}

func collectRefreshableMangaIDs(collection *mediaapi.MangaCollection) []int {
	if collection == nil || collection.GetMediaListCollection() == nil {
		return nil
	}

	ret := make([]int, 0)
	seen := make(map[int]struct{})
	for _, list := range collection.GetMediaListCollection().GetLists() {
		for _, entry := range list.GetEntries() {
			if entry == nil || entry.GetMedia() == nil {
				continue
			}

			mediaID := entry.GetMedia().GetID()
			if _, ok := seen[mediaID]; ok || !shouldRefreshSimulatedMedia(entry.GetStatus(), entry.GetMedia().GetStatus(), mediaID) {
				continue
			}

			seen[mediaID] = struct{}{}
			ret = append(ret, mediaID)
		}
	}

	return ret
}

func shouldRefreshSimulatedMedia(entryStatus *mediaapi.MediaListStatus, mediaStatus *mediaapi.MediaStatus, mediaID int) bool {
	if entryStatus == nil || mediaStatus == nil || customsource.IsExtensionId(mediaID) {
		return false
	}

	// todo: expand when simkl rate limits are less dogshit
	switch *entryStatus {
	case mediaapi.MediaListStatusCurrent:
	default:
		return false
	}

	switch *mediaStatus {
	case mediaapi.MediaStatusReleasing:
		return true
	default:
		return false
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper Methods
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (sp *SimulatedPlatform) getOrCreateAnimeCollection() (*mediaapi.AnimeCollection, error) {
	sp.collectionMu.RLock()
	if sp.animeCollection != nil {
		defer sp.collectionMu.RUnlock()
		return sp.animeCollection, nil
	}
	sp.collectionMu.RUnlock()

	sp.collectionMu.Lock()
	defer sp.collectionMu.Unlock()

	// Double-check after acquiring write lock
	if sp.animeCollection != nil {
		return sp.animeCollection, nil
	}

	// Try to load from database
	if collection := sp.localManager.GetSimulatedAnimeCollection(); collection.IsPresent() {
		sp.animeCollection = collection.MustGet()
		return sp.animeCollection, nil
	}

	// Create empty collection
	sp.animeCollection = &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{},
		},
	}

	// Save empty collection
	sp.localManager.SaveSimulatedAnimeCollection(sp.animeCollection)

	return sp.animeCollection, nil
}

func (sp *SimulatedPlatform) getOrCreateMangaCollection() (*mediaapi.MangaCollection, error) {
	sp.collectionMu.RLock()
	if sp.mangaCollection != nil {
		defer sp.collectionMu.RUnlock()
		return sp.mangaCollection, nil
	}
	sp.collectionMu.RUnlock()

	sp.collectionMu.Lock()
	defer sp.collectionMu.Unlock()

	// Double-check after acquiring write lock
	if sp.mangaCollection != nil {
		return sp.mangaCollection, nil
	}

	// Try to load from database
	if collection := sp.localManager.GetSimulatedMangaCollection(); collection.IsPresent() {
		sp.mangaCollection = collection.MustGet()
		return sp.mangaCollection, nil
	}

	// Create empty collection
	sp.mangaCollection = &mediaapi.MangaCollection{
		MediaListCollection: &mediaapi.MangaCollection_MediaListCollection{
			Lists: []*mediaapi.MangaCollection_MediaListCollection_Lists{},
		},
	}

	// Save empty collection
	sp.localManager.SaveSimulatedMangaCollection(sp.mangaCollection)

	return sp.mangaCollection, nil
}

func (sp *SimulatedPlatform) invalidateAnimeCollectionCache() {
	sp.collectionMu.Lock()
	defer sp.collectionMu.Unlock()
	sp.animeCollection = nil
}

func (sp *SimulatedPlatform) invalidateMangaCollectionCache() {
	sp.collectionMu.Lock()
	defer sp.collectionMu.Unlock()
	sp.mangaCollection = nil
}
