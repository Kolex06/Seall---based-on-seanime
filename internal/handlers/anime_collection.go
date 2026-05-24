package handlers

import (
	"context"
	"errors"
	"seall/internal/api/mediaapi"
	simklapi "seall/internal/api/simkl"
	"seall/internal/customsource"
	"seall/internal/database/db_bridge"
	"seall/internal/library/anime"
	"seall/internal/torrentstream"
	"seall/internal/util"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// HandleGetLibraryCollection
//
//	@summary returns the main local anime collection.
//	@desc This creates a new LibraryCollection struct and returns it.
//	@desc This is used to get the main anime collection of the user.
//	@desc It uses the cached Simkl anime collection for the GET method.
//	@desc It refreshes the SIMKL anime collection if the POST method is used.
//	@route /api/v1/library/collection [GET,POST]
//	@returns anime.LibraryCollection
func (h *Handler) HandleGetLibraryCollection(c echo.Context) error {

	animeCollection, err := h.App.GetMediaCollection(false)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	if animeCollection == nil {
		return h.RespondWithData(c, &anime.LibraryCollection{})
	}

	originalAnimeCollection := animeCollection

	var lfs []*anime.LocalFile
	// If using Nakama's library, fetch it
	nakamaLibrary, fromNakama := h.App.NakamaManager.GetHostAnimeLibrary(c.Request().Context())
	if fromNakama {
		// Save the original anime collection to restore it later
		originalAnimeCollection = animeCollection.Copy()
		lfs = nakamaLibrary.LocalFiles

		// Store all media from the user's collection
		userMediaIds := make(map[int]struct{})
		userCustomSourceMedia := make(map[string]map[int]struct{})
		for _, list := range animeCollection.MediaListCollection.GetLists() {
			for _, entry := range list.GetEntries() {
				mId := entry.GetMedia().GetID()
				userMediaIds[mId] = struct{}{}

				// Add all user custom source media to a map
				// This will be used to avoid duplicates
				if customsource.IsExtensionId(mId) {
					_, localId := customsource.ExtractExtensionData(mId)
					extensionId, ok := customsource.GetCustomSourceExtensionIdFromSiteUrl(entry.GetMedia().GetSiteURL())
					if !ok {
						// couldn't figure out the extension, skip it
						continue
					}
					if _, ok := userCustomSourceMedia[extensionId]; !ok {
						userCustomSourceMedia[extensionId] = make(map[int]struct{})
					}
					userCustomSourceMedia[extensionId][localId] = struct{}{}
				}
			}
		}

		// Store all custom source media from the Nakama host
		nakamaCustomSourceMediaIds := make(map[int]struct{})
		for _, lf := range lfs {
			if lf.MediaId > 0 {
				if customsource.IsExtensionId(lf.MediaId) {
					nakamaCustomSourceMediaIds[lf.MediaId] = struct{}{}
				}
			}
		}

		// Find media entries that are missing from the user's collection
		userMissingSimklMediaIds := make(map[int]struct{})
		for _, lf := range lfs {
			if lf.MediaId > 0 {
				if customsource.IsExtensionId(lf.MediaId) {
					continue
				}
				if _, ok := userMediaIds[lf.MediaId]; !ok {
					userMissingSimklMediaIds[lf.MediaId] = struct{}{}
				}
			}
		}

		nakamaCustomSourceMedia := make(map[int]*mediaapi.AnimeListEntry)

		// Add missing SIMKL entries to the user's collection as "Planning"
		for _, list := range nakamaLibrary.AnimeCollection.MediaListCollection.GetLists() {
			for _, entry := range list.GetEntries() {
				mId := entry.GetMedia().GetID()
				if _, ok := userMissingSimklMediaIds[mId]; ok {
					// create a new entry with blank list data
					newEntry := &mediaapi.AnimeListEntry{
						ID:     entry.GetID(),
						Media:  entry.GetMedia(),
						Status: &[]mediaapi.MediaListStatus{mediaapi.MediaListStatusPlanning}[0],
					}
					animeCollection.MediaListCollection.AddEntryToList(newEntry, mediaapi.MediaListStatusPlanning)
				}
				// Check if the media from a custom source
				if _, ok := nakamaCustomSourceMediaIds[mId]; ok {
					nakamaCustomSourceMedia[mId] = entry
				}
			}
		}

		// Add missing custom source entries to the user's collection as "Planning"
		// We'll find the equivalent
		if len(nakamaCustomSourceMedia) > 0 {
			// Go through all custom source media,
			// For each one, find the extension and replace the generated ID
			for mId, entry := range nakamaCustomSourceMedia {
				//extensionIdentifier, localId := customsource.ExtractExtensionData(mId)
				extensionId, ok := customsource.GetCustomSourceExtensionIdFromSiteUrl(entry.GetMedia().GetSiteURL())
				if !ok {
					// couldn't figure out the extension, skip it
					continue
				}

				_, localId := customsource.ExtractExtensionData(mId)

				// Find the same extension, if it's not installed, skip it
				customSource, ok := h.App.ExtensionRepository.GetCustomSourceExtensionByID(extensionId)
				if !ok {
					continue
				}

				// Generate a new ID for the custom source media
				newId := customsource.GenerateMediaId(customSource.GetExtensionIdentifier(), localId)
				entry.GetMedia().ID = newId

				// Add the entry if the user doesn't already have it
				if _, ok := userCustomSourceMedia[extensionId][localId]; !ok {
					newEntry := &mediaapi.AnimeListEntry{
						ID:     entry.GetID(),
						Media:  entry.GetMedia(),
						Status: &[]mediaapi.MediaListStatus{mediaapi.MediaListStatusPlanning}[0],
					}
					animeCollection.MediaListCollection.AddEntryToList(newEntry, mediaapi.MediaListStatusPlanning)
				}

				// Update the local files
				for _, lf := range lfs {
					if lf.MediaId == mId {
						lf.MediaId = newId
						break
					}
				}
			}
		}

	} else {
		lfs, _, err = db_bridge.GetLocalFiles(h.App.Database)
		if err != nil {
			return h.RespondWithError(c, err)
		}
	}

	libraryCollection, err := anime.NewLibraryCollection(c.Request().Context(), &anime.NewLibraryCollectionOptions{
		AnimeCollection:     animeCollection,
		PlatformRef:         h.App.MediaPlatformRef,
		LocalFiles:          lfs,
		MetadataProviderRef: h.App.MetadataProviderRef,
	})
	if err != nil {
		return h.RespondWithError(c, err)
	}

	// Restore the original anime collection if it was modified
	if fromNakama {
		*animeCollection = *originalAnimeCollection
	}

	if !fromNakama {
		if (h.App.SecondarySettings.Torrentstream != nil && h.App.SecondarySettings.Torrentstream.Enabled && h.App.SecondarySettings.Torrentstream.IncludeInLibrary) ||
			(h.App.Settings.GetLibrary() != nil && h.App.Settings.GetLibrary().EnableOnlinestream && h.App.Settings.GetLibrary().IncludeOnlineStreamingInLibrary) ||
			(h.App.SecondarySettings.Debrid != nil && h.App.SecondarySettings.Debrid.Enabled && h.App.SecondarySettings.Debrid.IncludeDebridStreamInLibrary) {
			h.App.TorrentstreamRepository.HydrateStreamCollection(&torrentstream.HydrateStreamCollectionOptions{
				AnimeCollection:     animeCollection,
				LibraryCollection:   libraryCollection,
				MetadataProviderRef: h.App.MetadataProviderRef,
			})
		}
	}

	// Add and remove necessary metadata when hydrating from Nakama
	if fromNakama {
		for _, ep := range libraryCollection.ContinueWatchingList {
			ep.IsNakamaEpisode = true
		}
		for _, list := range libraryCollection.Lists {
			for _, entry := range list.Entries {
				if entry.EntryLibraryData == nil {
					continue
				}
				entry.NakamaEntryLibraryData = &anime.NakamaEntryLibraryData{
					UnwatchedCount: entry.EntryLibraryData.UnwatchedCount,
					MainFileCount:  entry.EntryLibraryData.MainFileCount,
				}
				entry.EntryLibraryData = nil
			}
		}
	}

	// Hydrate total library size
	if libraryCollection != nil && libraryCollection.Stats != nil {
		libraryCollection.Stats.TotalSize = util.Bytes(h.App.TotalLibrarySize)
	}

	return h.RespondWithData(c, libraryCollection)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

// HandleGetMediaCollectionSchedule
//
//	@summary returns media collection schedule
//	@desc This is used by the "Schedule" page to display the media release schedule.
//	@route /api/v1/library/schedule [GET]
//	@returns []anime.ScheduleItem
func (h *Handler) HandleGetMediaCollectionSchedule(c echo.Context) error {
	source := normalizeScheduleSource(c.QueryParam("source"))

	// Invalidate the cache when the Simkl collection is refreshed
	h.App.AddOnRefreshMediaCollectionFunc("HandleGetMediaCollectionSchedule", func() {
		anime.ClearScheduleCache()
	})

	if ret, ok := anime.GetScheduleCache(source); ok {
		return h.RespondWithData(c, ret)
	}

	animeCollection, err := h.App.GetMediaCollection(false)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	calendarItems, calendarErr := h.getSIMKLCalendarItems(c.Request().Context(), time.Now())
	if calendarErr == nil {
		ret := anime.GetScheduleItemsFromSIMKLCalendar(calendarItems, animeCollection, source == "all")
		anime.SetScheduleCache(source, ret)
		return h.RespondWithData(c, ret)
	}

	animeSchedule, err := h.App.MediaPlatformRef.Get().GetAnimeAiringSchedule(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, calendarErr)
	}
	ret := anime.GetScheduleItems(animeSchedule, animeCollection)

	anime.SetScheduleCache(source, ret)

	return h.RespondWithData(c, ret)
}

func normalizeScheduleSource(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "all", "simkl", "calendar", "full":
		return "all"
	default:
		return "list"
	}
}

func (h *Handler) getSIMKLCalendarItems(ctx context.Context, now time.Time) ([]simklapi.CalendarItem, error) {
	return h.getSIMKLCalendarItemsFor(ctx, now, []simklapi.MediaType{
		simklapi.MediaTypeShows,
		simklapi.MediaTypeAnime,
		simklapi.MediaTypeMovies,
	}, 2)
}

func (h *Handler) getSIMKLCalendarItemsFor(ctx context.Context, now time.Time, mediaTypes []simklapi.MediaType, monthCount int) ([]simklapi.CalendarItem, error) {
	client := h.App.SimklClientRef.Get()
	if client == nil {
		return nil, errors.New("simkl calendar client is not available")
	}
	if monthCount < 1 {
		monthCount = 1
	}
	if len(mediaTypes) == 0 {
		mediaTypes = []simklapi.MediaType{
			simklapi.MediaTypeShows,
			simklapi.MediaTypeAnime,
			simklapi.MediaTypeMovies,
		}
	}

	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	months := make([]time.Time, 0, monthCount)
	for i := 0; i < monthCount; i++ {
		months = append(months, start.AddDate(0, i, 0))
	}

	ret := make([]simklapi.CalendarItem, 0)
	for _, month := range months {
		for _, mediaType := range mediaTypes {
			items, err := client.MonthlyCalendar(ctx, mediaType, month.Year(), int(month.Month()))
			if err != nil {
				return nil, err
			}
			ret = append(ret, items...)
		}
	}
	return ret, nil
}

// HandleAddUnknownMedia
//
//	@summary adds the given media to the user's SIMKL planning collections
//	@desc Since media not found in the user's SIMKL collection are not displayed in the library, this route is used to add them.
//	@desc The response is ignored in the frontend, the client should just refetch the entire library collection.
//	@route /api/v1/library/unknown-media [POST]
//	@returns mediaapi.AnimeCollection
func (h *Handler) HandleAddUnknownMedia(c echo.Context) error {

	type body struct {
		MediaIds []int `json:"mediaIds"`
	}

	b := new(body)
	if err := c.Bind(b); err != nil {
		return h.RespondWithError(c, err)
	}

	// Add non-added media entries to SIMKL collection
	if err := h.App.MediaPlatformRef.Get().AddMediaToCollection(c.Request().Context(), b.MediaIds); err != nil {
		return h.RespondWithError(c, errors.New("error: Simkl responded with an error, this is most likely a rate limit issue"))
	}

	// Bypass the cache
	animeCollection, err := h.App.GetMediaCollection(true)
	if err != nil {
		return h.RespondWithError(c, errors.New("error: Simkl responded with an error, wait one minute before refreshing"))
	}

	return h.RespondWithData(c, animeCollection)

}
