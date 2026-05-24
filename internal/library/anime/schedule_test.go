package anime_test

import (
	"seall/internal/api/mediaapi"
	simklapi "seall/internal/api/simkl"
	"seall/internal/customsource"
	"seall/internal/library/anime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetScheduleItemsFormatsDeduplicates(t *testing.T) {
	// schedule items are merged from all schedule buckets,
	// deduped by media/episode/time
	h := newAnimeTestWrapper(t)

	patchAnimeCollectionEntry(t, h.animeCollection, 154587, mediaapi.AnimeCollectionEntryPatch{
		Status:        new(mediaapi.MediaListStatusCurrent),
		AiredEpisodes: new(12),
	})
	patchCollectionEntryEpisodeCount(t, h.animeCollection, 154587, 12)

	patchAnimeCollectionEntry(t, h.animeCollection, 146065, mediaapi.AnimeCollectionEntryPatch{
		Status:        new(mediaapi.MediaListStatusCurrent),
		AiredEpisodes: new(1),
	})
	patchCollectionEntryEpisodeCount(t, h.animeCollection, 146065, 1)
	movieFormat := mediaapi.MediaFormatMovie
	patchCollectionEntryFormat(t, h.animeCollection, 146065, movieFormat)
	movieEntry := findCollectionEntryByMediaID(t, h.animeCollection, 146065)
	fallbackTitle := "movie fallback"
	movieEntry.Media.Title.UserPreferred = nil
	movieEntry.Media.Title.English = &fallbackTitle

	// extension-backed ids should not leak into the schedule list.
	extensionEntry := findCollectionEntryByMediaID(t, h.animeCollection, 21)
	extensionID := customsource.GenerateMediaId(1, 99)
	extensionEntry.Media.ID = extensionID
	extensionEntry.Status = new(mediaapi.MediaListStatusCurrent)

	animeSchedule := &mediaapi.AnimeAiringSchedule{
		Ongoing: &mediaapi.AnimeAiringSchedule_Ongoing{Media: []*mediaapi.AnimeSchedule{
			newAnimeSchedule(154587,
				[]*mediaapi.AnimeSchedule_Previous_Nodes{newPreviousScheduleNode(1_700_000_100, 11, -100)},
				[]*mediaapi.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_200, 12, 200)},
			),
			newAnimeSchedule(extensionID, nil, []*mediaapi.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_050, 1, 50)}),
		}},
		OngoingNext: &mediaapi.AnimeAiringSchedule_OngoingNext{Media: []*mediaapi.AnimeSchedule{
			newAnimeSchedule(154587, nil, []*mediaapi.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_200, 12, 200)}),
		}},
		Upcoming: &mediaapi.AnimeAiringSchedule_Upcoming{Media: []*mediaapi.AnimeSchedule{
			newAnimeSchedule(146065, nil, []*mediaapi.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_300, 1, 300)}),
		}},
	}

	items := anime.GetScheduleItems(animeSchedule, h.animeCollection)

	require.Len(t, items, 3)
	require.Len(t, findScheduleItems(items, 154587, 12), 1)
	require.Empty(t, findScheduleItems(items, extensionID, 1))

	previousItem := findScheduleItem(t, items, 154587, 11)
	require.Equal(t, time.Unix(1_700_000_100, 0).UTC(), previousItem.DateTime)
	require.Equal(t, previousItem.DateTime.Format("15:04"), previousItem.Time)
	require.False(t, previousItem.IsSeasonFinale)
	require.False(t, previousItem.IsMovie)

	finaleItem := findScheduleItem(t, items, 154587, 12)
	require.True(t, finaleItem.IsSeasonFinale)

	movieItem := findScheduleItem(t, items, 146065, 1)
	require.Equal(t, fallbackTitle, movieItem.Title)
	require.True(t, movieItem.IsMovie)
	require.True(t, movieItem.IsSeasonFinale)
}

func TestGetScheduleItemsHandlesNilInputs(t *testing.T) {
	// nil inputs should just give the caller an empty slice instead of exploding.
	require.Empty(t, anime.GetScheduleItems(nil, nil))
}

func TestGetScheduleItemsFromSIMKLCalendarMatchesCollection(t *testing.T) {
	currentStatus := mediaapi.MediaListStatusCurrent
	planningStatus := mediaapi.MediaListStatusPlanning
	showFormat := mediaapi.MediaFormatTv
	movieFormat := mediaapi.MediaFormatMovie
	showEpisodes := 12
	movieEpisodes := 1
	showTitle := "Collection show"
	movieTitle := "Collection movie"
	animeCollection := &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{
				{
					Status: &currentStatus,
					Entries: []*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{
						{
							ID:     154587,
							Status: &currentStatus,
							Media: &mediaapi.BaseAnime{
								ID:       154587,
								Format:   &showFormat,
								Episodes: &showEpisodes,
								Title: &mediaapi.BaseAnime_Title{
									UserPreferred: &showTitle,
								},
							},
						},
					},
				},
				{
					Status: &planningStatus,
					Entries: []*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{
						{
							ID:     146065,
							Status: &planningStatus,
							Media: &mediaapi.BaseAnime{
								ID:       146065,
								Format:   &movieFormat,
								Episodes: &movieEpisodes,
								Title: &mediaapi.BaseAnime_Title{
									UserPreferred: &movieTitle,
								},
							},
						},
					},
				},
			},
		},
	}

	calendarItems := []simklapi.CalendarItem{
		{
			Kind:  simklapi.MediaTypeShows,
			Title: "Collection show",
			Date:  "2026-05-22T23:00:00-05:00",
			IDs:   simklapi.IDs{SimklID: 154587},
			Episode: &simklapi.CalendarEpisode{
				Episode: 12,
			},
		},
		{
			Kind:        simklapi.MediaTypeMovies,
			Title:       "Collection movie",
			Date:        "2026-05-23T00:00:00-05:00",
			ReleaseDate: "2026-05-23",
			IDs:         simklapi.IDs{SimklID: 146065},
		},
		{
			Kind:  simklapi.MediaTypeShows,
			Title: "Not in collection",
			Date:  "2026-05-24T00:00:00-05:00",
			IDs:   simklapi.IDs{SimklID: 999999},
			Episode: &simklapi.CalendarEpisode{
				Episode: 1,
			},
		},
	}

	items := anime.GetScheduleItemsFromSIMKLCalendar(calendarItems, animeCollection, true)

	require.Len(t, items, 3)

	showItem := findScheduleItem(t, items, 154587, 12)
	require.Equal(t, time.Date(2026, 5, 23, 4, 0, 0, 0, time.UTC), showItem.DateTime)
	require.True(t, showItem.IsSeasonFinale)
	require.False(t, showItem.IsMovie)
	require.Equal(t, string(simklapi.MediaTypeShows), showItem.MediaType)
	require.NotNil(t, showItem.ListStatus)
	require.Equal(t, currentStatus, *showItem.ListStatus)

	movieItem := findScheduleItem(t, items, 146065, 1)
	require.True(t, movieItem.IsMovie)
	require.True(t, movieItem.IsSeasonFinale)
	require.Equal(t, string(simklapi.MediaTypeMovies), movieItem.MediaType)

	publicItem := findScheduleItem(t, items, 999999, 1)
	require.Equal(t, "Not in collection", publicItem.Title)
	require.Equal(t, string(simklapi.MediaTypeShows), publicItem.MediaType)
	require.Nil(t, publicItem.ListStatus)

	listItems := anime.GetScheduleItemsFromSIMKLCalendar(calendarItems, animeCollection, false)
	require.Len(t, listItems, 2)
	require.NotNil(t, findScheduleItem(t, listItems, 154587, 12).ListStatus)
	require.NotNil(t, findScheduleItem(t, listItems, 146065, 1).ListStatus)
	require.Empty(t, findScheduleItems(listItems, 999999, 1))
}

func newAnimeSchedule(mediaID int, previous []*mediaapi.AnimeSchedule_Previous_Nodes, upcoming []*mediaapi.AnimeSchedule_Upcoming_Nodes) *mediaapi.AnimeSchedule {
	ret := &mediaapi.AnimeSchedule{ID: mediaID}
	if previous != nil {
		ret.Previous = &mediaapi.AnimeSchedule_Previous{Nodes: previous}
	}
	if upcoming != nil {
		ret.Upcoming = &mediaapi.AnimeSchedule_Upcoming{Nodes: upcoming}
	}
	return ret
}

func newPreviousScheduleNode(airingAt int, episode int, timeUntilAiring int) *mediaapi.AnimeSchedule_Previous_Nodes {
	return &mediaapi.AnimeSchedule_Previous_Nodes{
		AiringAt:        airingAt,
		Episode:         episode,
		TimeUntilAiring: timeUntilAiring,
	}
}

func newUpcomingScheduleNode(airingAt int, episode int, timeUntilAiring int) *mediaapi.AnimeSchedule_Upcoming_Nodes {
	return &mediaapi.AnimeSchedule_Upcoming_Nodes{
		AiringAt:        airingAt,
		Episode:         episode,
		TimeUntilAiring: timeUntilAiring,
	}
}

func findScheduleItem(t *testing.T, items []*anime.ScheduleItem, mediaID int, episodeNumber int) *anime.ScheduleItem {
	t.Helper()
	matching := findScheduleItems(items, mediaID, episodeNumber)
	require.Len(t, matching, 1)
	return matching[0]
}

func findScheduleItems(items []*anime.ScheduleItem, mediaID int, episodeNumber int) []*anime.ScheduleItem {
	ret := make([]*anime.ScheduleItem, 0)
	for _, item := range items {
		if item.MediaId == mediaID && item.EpisodeNumber == episodeNumber {
			ret = append(ret, item)
		}
	}
	return ret
}
