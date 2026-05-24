package media_platform

import (
	"context"
	"errors"
	"seall/internal/api/mediaapi"
	"seall/internal/customsource"
	"seall/internal/database/db"
	"seall/internal/extension"
	"seall/internal/hook"
	"seall/internal/platforms/platform"
	"seall/internal/platforms/shared_platform"
	"seall/internal/util"
	"seall/internal/util/limiter"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/mo"
)

type (
	MediaPlatform struct {
		logger                 *zerolog.Logger
		username               mo.Option[string]
		simklClient            mediaapi.MediaApiClient
		useFixtureCollections  bool
		animeCollection        mo.Option[*mediaapi.AnimeCollection]
		rawAnimeCollection     mo.Option[*mediaapi.AnimeCollection]
		mangaCollection        mo.Option[*mediaapi.MangaCollection]
		rawMangaCollection     mo.Option[*mediaapi.MangaCollection]
		isOffline              bool
		offlinePlatformEnabled bool
		helper                 *shared_platform.PlatformHelper
		db                     *db.Database
		extensionBankRef       *util.Ref[*extension.UnifiedBank]
	}
)

func NewMediaPlatform(simklClientRef *util.Ref[mediaapi.MediaApiClient], extensionBankRef *util.Ref[*extension.UnifiedBank], logger *zerolog.Logger, db *db.Database, logoutFunc ...func()) platform.Platform {
	_, useFixtureCollections := simklClientRef.Get().(*mediaapi.FixtureMediaApiClient)

	ap := &MediaPlatform{
		simklClient:           shared_platform.NewCacheLayer(simklClientRef, logoutFunc...),
		logger:                logger,
		username:              mo.None[string](),
		useFixtureCollections: useFixtureCollections,
		animeCollection:       mo.None[*mediaapi.AnimeCollection](),
		rawAnimeCollection:    mo.None[*mediaapi.AnimeCollection](),
		mangaCollection:       mo.None[*mediaapi.MangaCollection](),
		rawMangaCollection:    mo.None[*mediaapi.MangaCollection](),
		extensionBankRef:      extensionBankRef,
		helper:                shared_platform.NewPlatformHelper(extensionBankRef, db, logger),
		db:                    db,
	}

	return ap
}

func (ap *MediaPlatform) ClearCache() {
	ap.helper.ClearCache()
}

func (ap *MediaPlatform) Close() {
	ap.helper.Close()
}

func (ap *MediaPlatform) SetUsername(username string) {
	// Set the username for the MediaPlatform
	if username == "" {
		ap.username = mo.Some[string]("")
		return
	}

	ap.username = mo.Some(username)
}

func (ap *MediaPlatform) SetMediaApiClient(client mediaapi.MediaApiClient) {
	// Set the MediaApiClient for the MediaPlatform
	ap.simklClient = client
}

func (ap *MediaPlatform) getUsername() (*string, bool) {
	if ap.username.IsPresent() {
		return ap.username.ToPointer(), true
	}
	if ap.useFixtureCollections {
		return nil, true
	}

	return nil, false
}

func (ap *MediaPlatform) UpdateEntry(ctx context.Context, mediaID int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput) error {
	ap.logger.Trace().Msg("simkl platform: Updating entry")

	// Use shared hook handling
	return ap.helper.TriggerUpdateEntryHooks(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt, func(event *platform.PreUpdateEntryEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := ap.helper.HandleCustomSourceUpdateEntry(ctx, mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt); handled {
			return err
		}

		_, err := ap.simklClient.UpdateMediaListEntry(ctx, event.MediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt)
		return err
	})
}

func (ap *MediaPlatform) UpdateEntryProgress(ctx context.Context, mediaID int, progress int, totalCount *int) error {
	ap.logger.Trace().Msg("simkl platform: Updating entry progress")

	// Use shared hook handling
	return ap.helper.TriggerUpdateEntryProgressHooks(ctx, mediaID, progress, totalCount, func(event *platform.PreUpdateEntryProgressEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := ap.helper.HandleCustomSourceUpdateEntryProgress(ctx, mediaID, *event.Progress, event.TotalCount); handled {
			return err
		}

		realTotalCount := 0
		if totalCount != nil && *totalCount > 0 {
			realTotalCount = *totalCount
		}

		// Check if the anime is in the repeating list
		// If it is, set the status to repeating
		if ap.rawAnimeCollection.IsPresent() {
			for _, list := range ap.rawAnimeCollection.MustGet().MediaListCollection.Lists {
				if list.Status != nil && *list.Status == mediaapi.MediaListStatusRepeating {
					if list.Entries != nil {
						for _, entry := range list.Entries {
							if entry.GetMedia().GetID() == mediaID {
								*event.Status = mediaapi.MediaListStatusRepeating
								break
							}
						}
					}
				}
			}
		}
		if realTotalCount > 0 && *event.Progress >= realTotalCount {
			*event.Status = mediaapi.MediaListStatusCompleted
		}

		if realTotalCount > 0 && *event.Progress > realTotalCount {
			*event.Progress = realTotalCount
		}

		_, err := ap.simklClient.UpdateMediaListEntryProgress(
			ctx,
			event.MediaID,
			event.Progress,
			event.Status,
		)
		return err
	})
}

func (ap *MediaPlatform) UpdateEntryRepeat(ctx context.Context, mediaID int, repeat int) error {
	ap.logger.Trace().Msg("simkl platform: Updating entry repeat")

	// Use shared hook handling
	return ap.helper.TriggerUpdateEntryRepeatHooks(ctx, mediaID, repeat, func(event *platform.PreUpdateEntryRepeatEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := ap.helper.HandleCustomSourceUpdateEntryRepeat(ctx, mediaID, *event.Repeat); handled {
			return err
		}

		_, err := ap.simklClient.UpdateMediaListEntryRepeat(ctx, event.MediaID, event.Repeat)
		return err
	})
}

func (ap *MediaPlatform) DeleteEntry(ctx context.Context, mediaID, entryId int) error {
	ap.logger.Trace().Msg("simkl platform: Deleting entry")

	return ap.helper.TriggerDeleteEntryHooks(ctx, mediaID, entryId, func(event *platform.PreDeleteEntryEvent) error {
		if handled, err := ap.helper.HandleCustomSourceDeleteEntry(ctx, *event.MediaID, *event.EntryID); handled {
			return err
		}

		_, err := ap.simklClient.DeleteEntry(ctx, event.EntryID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (ap *MediaPlatform) GetAnime(ctx context.Context, mediaID int) (*mediaapi.BaseAnime, error) {
	ap.logger.Trace().Int("mediaId", mediaID).Msg("simkl platform: Fetching anime")

	if cachedAnime, ok := ap.helper.GetCachedBaseAnime(mediaID); ok {
		ap.logger.Trace().Msg("simkl platform: Returning anime from cache")
		return ap.helper.TriggerGetAnimeEvent(cachedAnime)
	}

	// Check if this is a custom source entry
	if media, isCustom, err := ap.helper.HandleCustomSourceAnime(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}

		triggeredMedia, err := ap.helper.TriggerGetAnimeEvent(media)
		if err != nil {
			return nil, err
		}

		ap.helper.SetCachedBaseAnime(mediaID, triggeredMedia)
		return triggeredMedia, nil
	}

	// Get from SIMKL
	ret, err := ap.simklClient.BaseAnimeByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}

	media := ret.GetMedia()
	triggeredMedia, err := ap.helper.TriggerGetAnimeEvent(media)
	if err != nil {
		return nil, err
	}

	ap.helper.SetCachedBaseAnime(mediaID, triggeredMedia)
	return triggeredMedia, nil
}

func (ap *MediaPlatform) GetAnimeByMalID(ctx context.Context, malID int) (*mediaapi.BaseAnime, error) {
	ap.logger.Trace().Msg("simkl platform: Fetching anime by MAL ID")
	ret, err := ap.simklClient.BaseAnimeByMalID(ctx, &malID)
	if err != nil {
		return nil, err
	}

	media := ret.GetMedia()
	return ap.helper.TriggerGetAnimeEvent(media)
}

func (ap *MediaPlatform) GetAnimeDetails(ctx context.Context, mediaID int) (*mediaapi.AnimeDetailsById_Media, error) {
	ap.logger.Trace().Int("mediaId", mediaID).Msg("simkl platform: Fetching anime details")

	// Check if this is a custom source entry
	if media, isCustom, err := ap.helper.HandleCustomSourceAnimeDetails(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		return ap.helper.TriggerGetAnimeDetailsEvent(media)
	}

	// Get from SIMKL
	ret, err := ap.simklClient.AnimeDetailsByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}

	media := ret.GetMedia()
	return ap.helper.TriggerGetAnimeDetailsEvent(media)
}

func (ap *MediaPlatform) GetAnimeWithRelations(ctx context.Context, mediaID int) (*mediaapi.CompleteAnime, error) {
	ap.logger.Trace().Int("mediaId", mediaID).Msg("simkl platform: Fetching anime with relations")

	if cachedAnime, ok := ap.helper.GetCachedCompleteAnime(mediaID); ok {
		ap.logger.Trace().Msg("simkl platform: Cache HIT for anime with relations")
		return cachedAnime, nil
	}

	// Check if this is a custom source entry
	if media, isCustom, err := ap.helper.HandleCustomSourceAnimeWithRelations(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		ap.helper.SetCachedCompleteAnime(mediaID, media)
		return media, nil
	}

	// Get from SIMKL
	ret, err := ap.simklClient.CompleteAnimeByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}
	media := ret.GetMedia()

	ap.helper.SetCachedCompleteAnime(mediaID, media)
	return media, nil
}

func (ap *MediaPlatform) GetManga(ctx context.Context, mediaID int) (*mediaapi.BaseManga, error) {
	ap.logger.Trace().Int("mediaId", mediaID).Msg("simkl platform: Fetching manga")

	if cachedManga, ok := ap.helper.GetCachedBaseManga(mediaID); ok {
		ap.logger.Trace().Msg("simkl platform: Returning manga from cache")
		return ap.helper.TriggerGetMangaEvent(cachedManga)
	}

	// Check if this is a custom source entry
	if media, isCustom, err := ap.helper.HandleCustomSourceManga(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}

		triggeredMedia, err := ap.helper.TriggerGetMangaEvent(media)
		if err != nil {
			return nil, err
		}

		ap.helper.SetCachedBaseManga(mediaID, triggeredMedia)
		return triggeredMedia, nil
	}

	// Get from SIMKL
	ret, err := ap.simklClient.BaseMangaByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}

	media := ret.GetMedia()
	triggeredMedia, err := ap.helper.TriggerGetMangaEvent(media)
	if err != nil {
		return nil, err
	}

	ap.helper.SetCachedBaseManga(mediaID, triggeredMedia)
	return triggeredMedia, nil
}

func (ap *MediaPlatform) GetMangaDetails(ctx context.Context, mediaID int) (*mediaapi.MangaDetailsById_Media, error) {
	ap.logger.Trace().Msg("simkl platform: Fetching manga details")

	// Check if this is a custom source entry
	if media, isCustom, err := ap.helper.HandleCustomSourceMangaDetails(ctx, mediaID); isCustom {
		return media, err
	}

	// Get from SIMKL
	ret, err := ap.simklClient.MangaDetailsByID(ctx, &mediaID)
	if err != nil {
		return nil, err
	}

	return ret.GetMedia(), nil
}

func (ap *MediaPlatform) GetMediaCollection(ctx context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error) {
	if !bypassCache && ap.animeCollection.IsPresent() {
		event := new(platform.GetCachedAnimeCollectionEvent)
		event.AnimeCollection = ap.animeCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedAnimeCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.AnimeCollection, nil
	}

	if _, ok := ap.getUsername(); !ok {
		return nil, nil
	}

	err := ap.refreshAnimeCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetMediaCollectionEvent)
	event.AnimeCollection = ap.animeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetMediaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (ap *MediaPlatform) GetRawMediaCollection(ctx context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error) {
	if !bypassCache && ap.rawAnimeCollection.IsPresent() {
		event := new(platform.GetCachedRawAnimeCollectionEvent)
		event.AnimeCollection = ap.rawAnimeCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedRawAnimeCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.AnimeCollection, nil
	}

	if _, ok := ap.getUsername(); !ok {
		return nil, nil
	}

	err := ap.refreshAnimeCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetRawMediaCollectionEvent)
	event.AnimeCollection = ap.rawAnimeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawMediaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (ap *MediaPlatform) RefreshAnimeCollection(ctx context.Context) (*mediaapi.AnimeCollection, error) {
	if _, ok := ap.getUsername(); !ok {
		return nil, nil
	}

	err := ap.refreshAnimeCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetMediaCollectionEvent)
	event.AnimeCollection = ap.animeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetMediaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	event2 := new(platform.GetRawMediaCollectionEvent)
	event2.AnimeCollection = ap.rawAnimeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawMediaCollection().Trigger(event2)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (ap *MediaPlatform) refreshAnimeCollection(ctx context.Context) error {
	userName, ok := ap.getUsername()
	if !ok {
		return errors.New("simkl: Username is not set")
	}

	// Else, get the collection from Simkl
	collection, err := ap.simklClient.AnimeCollection(ctx, userName)
	if err != nil {
		return err
	}

	// Merge the custom entries into the collection
	ap.helper.MergeCustomSourceAnimeEntries(collection)

	// Save the raw collection to App (retains the lists with no status)
	ap.rawAnimeCollection = mo.Some(new(*collection))
	ap.rawAnimeCollection.MustGet().MediaListCollection = new(*collection.MediaListCollection)
	listsCopy := make([]*mediaapi.AnimeCollection_MediaListCollection_Lists, len(collection.MediaListCollection.Lists))
	copy(listsCopy, collection.MediaListCollection.Lists)
	ap.rawAnimeCollection.MustGet().MediaListCollection.Lists = listsCopy

	// Remove lists with no status (custom lists)
	collection.MediaListCollection.Lists = ap.helper.FilterOutCustomAnimeLists(collection.MediaListCollection.Lists)

	// Save the collection to App
	ap.animeCollection = mo.Some(collection)

	return nil
}

func (ap *MediaPlatform) GetMediaCollectionWithRelations(ctx context.Context) (*mediaapi.AnimeCollectionWithRelations, error) {
	ap.logger.Trace().Msg("simkl platform: Fetching anime collection with relations")

	userName, ok := ap.getUsername()
	if !ok {
		return nil, nil
	}

	ret, err := ap.simklClient.AnimeCollectionWithRelations(ctx, userName)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (ap *MediaPlatform) GetMangaCollection(ctx context.Context, bypassCache bool) (*mediaapi.MangaCollection, error) {

	if !bypassCache && ap.mangaCollection.IsPresent() {
		event := new(platform.GetCachedMangaCollectionEvent)
		event.MangaCollection = ap.mangaCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedMangaCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.MangaCollection, nil
	}

	if _, ok := ap.getUsername(); !ok {
		return nil, nil
	}

	err := ap.refreshMangaCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetMangaCollectionEvent)
	event.MangaCollection = ap.mangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (ap *MediaPlatform) GetRawMangaCollection(ctx context.Context, bypassCache bool) (*mediaapi.MangaCollection, error) {
	ap.logger.Trace().Msg("simkl platform: Fetching raw manga collection")

	if !bypassCache && ap.rawMangaCollection.IsPresent() {
		ap.logger.Trace().Msg("simkl platform: Returning raw manga collection from cache")
		event := new(platform.GetCachedRawMangaCollectionEvent)
		event.MangaCollection = ap.rawMangaCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedRawMangaCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.MangaCollection, nil
	}

	if _, ok := ap.getUsername(); !ok {
		return nil, nil
	}

	err := ap.refreshMangaCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetRawMangaCollectionEvent)
	event.MangaCollection = ap.rawMangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (ap *MediaPlatform) RefreshMangaCollection(ctx context.Context) (*mediaapi.MangaCollection, error) {
	if _, ok := ap.getUsername(); !ok {
		return nil, nil
	}

	err := ap.refreshMangaCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetMangaCollectionEvent)
	event.MangaCollection = ap.mangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	event2 := new(platform.GetRawMangaCollectionEvent)
	event2.MangaCollection = ap.rawMangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawMangaCollection().Trigger(event2)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (ap *MediaPlatform) refreshMangaCollection(ctx context.Context) error {
	userName, ok := ap.getUsername()
	if !ok {
		return errors.New("simkl: Username is not set")
	}

	collection, err := ap.simklClient.MangaCollection(ctx, userName)
	if err != nil {
		return err
	}

	// Merge the custom entries into the collection
	ap.helper.MergeCustomSourceMangaEntries(collection)

	// Save the raw collection to App (retains the lists with no status)
	ap.rawMangaCollection = mo.Some(new(*collection))
	ap.rawMangaCollection.MustGet().MediaListCollection = new(*collection.MediaListCollection)
	listsCopy := make([]*mediaapi.MangaCollection_MediaListCollection_Lists, len(collection.MediaListCollection.Lists))
	copy(listsCopy, collection.MediaListCollection.Lists)
	ap.rawMangaCollection.MustGet().MediaListCollection.Lists = listsCopy

	// Remove lists with no status (custom lists)
	collection.MediaListCollection.Lists = ap.helper.FilterOutCustomMangaLists(collection.MediaListCollection.Lists)

	// Remove Novels from both collections
	ap.helper.RemoveNovelsFromMangaCollection(collection)
	ap.helper.RemoveNovelsFromMangaCollection(ap.rawMangaCollection.MustGet())

	// Save the collection to App
	ap.mangaCollection = mo.Some(collection)

	return nil
}

func (ap *MediaPlatform) AddMediaToCollection(ctx context.Context, mIds []int) error {
	ap.logger.Trace().Msg("simkl platform: Adding media to collection")
	if len(mIds) == 0 {
		ap.logger.Debug().Msg("simkl: No media added to planning list")
		return nil
	}

	rateLimiter := limiter.NewLimiter(1*time.Second, 1) // 1 request per second

	wg := sync.WaitGroup{}
	for _, _id := range mIds {
		wg.Add(1)
		go func(id int) {
			rateLimiter.Wait()
			defer wg.Done()

			if customsource.IsExtensionId(id) {
				_, err := ap.helper.HandleCustomSourceUpdateEntry(ctx,
					id,
					new(mediaapi.MediaListStatusPlanning),
					new(0),
					new(0),
					nil,
					nil,
				)
				if err != nil {
					ap.logger.Error().Msg("simkl: An error occurred while adding media to planning list: " + err.Error())
				}
				return
			}

			_, err := ap.simklClient.UpdateMediaListEntry(
				ctx,
				&id,
				new(mediaapi.MediaListStatusPlanning),
				new(0),
				new(0),
				nil,
				nil,
			)
			if err != nil {
				ap.logger.Error().Msg("simkl: An error occurred while adding media to planning list: " + err.Error())
			}
		}(_id)
	}
	wg.Wait()

	ap.logger.Debug().Any("count", len(mIds)).Msg("simkl: Media added to planning list")
	return nil
}

func (ap *MediaPlatform) GetStudioDetails(ctx context.Context, studioID int) (*mediaapi.StudioDetails, error) {
	ap.logger.Trace().Msg("simkl platform: Fetching studio details")
	ret, err := ap.simklClient.StudioDetails(ctx, &studioID)
	if err != nil {
		return nil, err
	}

	return ap.helper.TriggerGetStudioDetailsEvent(ret)
}

func (ap *MediaPlatform) GetMediaApiClient() mediaapi.MediaApiClient {
	return ap.simklClient
}

func (ap *MediaPlatform) GetViewerStats(ctx context.Context) (*mediaapi.ViewerStats, error) {
	if ap.username.IsAbsent() {
		return nil, errors.New("simkl: Username is not set")
	}

	ap.logger.Trace().Msg("simkl platform: Fetching viewer stats")
	ret, err := ap.simklClient.ViewerStats(ctx)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (ap *MediaPlatform) GetAnimeAiringSchedule(ctx context.Context) (*mediaapi.AnimeAiringSchedule, error) {
	if ap.username.IsAbsent() {
		return nil, errors.New("simkl: Username is not set")
	}

	collection, err := ap.GetMediaCollection(ctx, false)
	if err != nil {
		return nil, err
	}

	return ap.helper.BuildAnimeAiringSchedule(ctx, collection, ap.simklClient)
}
