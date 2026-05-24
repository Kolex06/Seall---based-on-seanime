package anime

import (
	"fmt"
	"seall/internal/api/mediaapi"
	simklapi "seall/internal/api/simkl"
	"seall/internal/customsource"
	"seall/internal/hook"
	"seall/internal/util/result"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

type ScheduleItem struct {
	MediaId int    `json:"mediaId"`
	Title   string `json:"title"`
	// MediaType is the SIMKL media bucket: movies, shows, or anime.
	MediaType string `json:"mediaType,omitempty"`
	// Time is in 15:04 format, UTC.
	// The frontend should derive local time from DateTime instead.
	Time string `json:"time"`
	// DateTime is in UTC
	DateTime       time.Time                 `json:"dateTime"`
	Image          string                    `json:"image"`
	EpisodeNumber  int                       `json:"episodeNumber"`
	IsMovie        bool                      `json:"isMovie"`
	IsSeasonFinale bool                      `json:"isSeasonFinale"`
	ListStatus     *mediaapi.MediaListStatus `json:"listStatus,omitempty"`
	ListProgress   *int                      `json:"listProgress,omitempty"`
}

var scheduleCache = result.NewCache[string, []*ScheduleItem]()

func GetScheduleCache(source string) ([]*ScheduleItem, bool) {
	return scheduleCache.Get(source)
}

func SetScheduleCache(source string, val []*ScheduleItem) {
	scheduleCache.Set(source, val)
}

func ClearScheduleCache() {
	scheduleCache.Clear()
}

func GetScheduleItems(animeSchedule *mediaapi.AnimeAiringSchedule, animeCollection *mediaapi.AnimeCollection) []*ScheduleItem {
	if animeSchedule == nil || animeCollection == nil || animeCollection.MediaListCollection == nil {
		return []*ScheduleItem{}
	}

	animeEntryMap := make(map[int]*mediaapi.AnimeListEntry)
	for _, list := range animeCollection.MediaListCollection.GetLists() {
		for _, entry := range list.GetEntries() {
			if customsource.IsExtensionId(entry.Media.GetID()) {
				continue
			}
			animeEntryMap[entry.GetMedia().GetID()] = entry
		}
	}

	type animeScheduleNode interface {
		GetAiringAt() int
		GetTimeUntilAiring() int
		GetEpisode() int
	}

	type animeScheduleMedia interface {
		GetMedia() []*mediaapi.AnimeSchedule
	}

	formatNodeItem := func(node animeScheduleNode, entry *mediaapi.AnimeListEntry) *ScheduleItem {
		t := time.Unix(int64(node.GetAiringAt()), 0)
		item := &ScheduleItem{
			MediaId:        entry.GetMedia().GetID(),
			Title:          entry.GetMedia().GetPreferredTitle(),
			MediaType:      string(simklapi.MediaTypeAnime),
			Time:           t.UTC().Format("15:04"),
			DateTime:       t.UTC(),
			Image:          entry.GetMedia().GetCoverImageSafe(),
			EpisodeNumber:  node.GetEpisode(),
			IsMovie:        entry.GetMedia().IsMovie(),
			IsSeasonFinale: false,
			ListStatus:     entry.Status,
			ListProgress:   entry.Progress,
		}
		if entry.GetMedia().GetTotalEpisodeCount() > 0 && node.GetEpisode() == entry.GetMedia().GetTotalEpisodeCount() {
			item.IsSeasonFinale = true
		}
		return item
	}

	formatPart := func(m animeScheduleMedia) ([]*ScheduleItem, bool) {
		if m == nil {
			return nil, false
		}
		ret := make([]*ScheduleItem, 0)
		for _, m := range m.GetMedia() {
			entry, ok := animeEntryMap[m.GetID()]
			if !ok || entry.Status == nil || *entry.Status == mediaapi.MediaListStatusDropped {
				continue
			}
			for _, n := range m.GetPrevious().GetNodes() {
				ret = append(ret, formatNodeItem(n, entry))
			}
			for _, n := range m.GetUpcoming().GetNodes() {
				ret = append(ret, formatNodeItem(n, entry))
			}
		}
		return ret, true
	}

	ongoingItems, _ := formatPart(animeSchedule.GetOngoing())
	ongoingNextItems, _ := formatPart(animeSchedule.GetOngoingNext())
	precedingItems, _ := formatPart(animeSchedule.GetPreceding())
	upcomingItems, _ := formatPart(animeSchedule.GetUpcoming())
	upcomingNextItems, _ := formatPart(animeSchedule.GetUpcomingNext())

	allItems := make([]*ScheduleItem, 0)
	allItems = append(allItems, ongoingItems...)
	allItems = append(allItems, ongoingNextItems...)
	allItems = append(allItems, precedingItems...)
	allItems = append(allItems, upcomingItems...)
	allItems = append(allItems, upcomingNextItems...)

	return FinalizeScheduleItems(animeCollection, allItems)
}

func GetScheduleItemsFromSIMKLCalendar(calendarItems []simklapi.CalendarItem, animeCollection *mediaapi.AnimeCollection, includeUnlisted bool) []*ScheduleItem {
	typedEntryMap := make(map[string]*mediaapi.AnimeListEntry)
	entryMap := make(map[int]*mediaapi.AnimeListEntry)
	entryIDCounts := make(map[int]int)
	if animeCollection != nil && animeCollection.MediaListCollection != nil {
		for _, list := range animeCollection.MediaListCollection.GetLists() {
			for _, entry := range list.GetEntries() {
				if entry == nil || entry.GetMedia() == nil {
					continue
				}
				mediaID := entry.GetMedia().GetID()
				if customsource.IsExtensionId(mediaID) {
					continue
				}
				typedEntryMap[scheduleEntryKey(inferSIMKLMediaType(entry.GetMedia()), mediaID)] = entry
				entryMap[mediaID] = entry
				entryIDCounts[mediaID]++
			}
		}
	}

	ret := make([]*ScheduleItem, 0, len(calendarItems))
	for _, calendarItem := range calendarItems {
		mediaID := calendarItem.IDs.PrimarySimklID()
		if mediaID == 0 {
			mediaID = simklIDFromCalendarURL(calendarItem.URL)
		}

		dateTime, ok := parseSIMKLCalendarDate(calendarItem.Date, calendarItem.ReleaseDate)
		if !ok {
			continue
		}

		episodeNumber := 1
		if calendarItem.Episode != nil && calendarItem.Episode.Episode > 0 {
			episodeNumber = calendarItem.Episode.Episode
		}

		entry := typedEntryMap[scheduleEntryKey(calendarItem.Kind, mediaID)]
		if entry == nil && entryIDCounts[mediaID] == 1 {
			entry = entryMap[mediaID]
		}
		if entry != nil && entry.Status != nil && *entry.Status == mediaapi.MediaListStatusDropped {
			continue
		}
		if entry == nil && !includeUnlisted {
			continue
		}

		item := &ScheduleItem{
			MediaId:       mediaID,
			Title:         firstNonEmptyString(calendarItem.Title, "Untitled"),
			MediaType:     string(calendarItem.Kind),
			Time:          dateTime.UTC().Format("15:04"),
			DateTime:      dateTime.UTC(),
			Image:         simklapi.ImageURL(simklapi.ImageKindPoster, calendarItem.Poster, simklapi.ImageSizePosterMedium),
			EpisodeNumber: episodeNumber,
			IsMovie:       calendarItem.Kind == simklapi.MediaTypeMovies,
		}
		if entry != nil {
			media := entry.GetMedia()
			item.Title = firstNonEmptyString(media.GetPreferredTitle(), calendarItem.Title, "Untitled")
			item.Image = firstNonEmptyString(media.GetCoverImageSafe(), item.Image)
			item.IsMovie = media.IsMovie() || item.IsMovie
			item.ListStatus = entry.Status
			item.ListProgress = entry.Progress
			if media.GetTotalEpisodeCount() > 0 && episodeNumber == media.GetTotalEpisodeCount() {
				item.IsSeasonFinale = true
			}
		} else if item.IsMovie {
			item.IsSeasonFinale = true
		}
		ret = append(ret, item)
	}

	return FinalizeScheduleItems(animeCollection, ret)
}

func FinalizeScheduleItems(animeCollection *mediaapi.AnimeCollection, items []*ScheduleItem) []*ScheduleItem {
	ret := lo.UniqBy(items, func(item *ScheduleItem) string {
		if item == nil {
			return ""
		}
		return fmt.Sprintf("%s-%d-%d-%d", item.MediaType, item.MediaId, item.EpisodeNumber, item.DateTime.Unix())
	})

	ret = lo.Filter(ret, func(item *ScheduleItem, _ int) bool {
		return item != nil
	})

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].DateTime.Before(ret[j].DateTime)
	})

	event := &AnimeScheduleItemsEvent{
		AnimeCollection: animeCollection,
		Items:           ret,
	}
	err := hook.GlobalHookManager.OnAnimeScheduleItems().Trigger(event)
	if err != nil {
		return ret
	}

	return event.Items
}

func scheduleEntryKey(mediaType simklapi.MediaType, mediaID int) string {
	return fmt.Sprintf("%s:%d", mediaType, mediaID)
}

func inferSIMKLMediaType(media *mediaapi.BaseAnime) simklapi.MediaType {
	if media == nil {
		return simklapi.MediaTypeShows
	}
	if media.IsMovie() {
		return simklapi.MediaTypeMovies
	}
	siteURL := ""
	if media.SiteURL != nil {
		siteURL = strings.ToLower(*media.SiteURL)
	}
	switch {
	case strings.Contains(siteURL, "/anime/"):
		return simklapi.MediaTypeAnime
	case strings.Contains(siteURL, "/movie/"), strings.Contains(siteURL, "/movies/"):
		return simklapi.MediaTypeMovies
	case strings.Contains(siteURL, "/tv/"), strings.Contains(siteURL, "/show/"), strings.Contains(siteURL, "/shows/"):
		return simklapi.MediaTypeShows
	default:
		return simklapi.MediaTypeShows
	}
}

func parseSIMKLCalendarDate(values ...string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
	}
	for _, value := range values {
		if value == "" {
			continue
		}
		for _, layout := range layouts {
			if parsed, err := time.Parse(layout, value); err == nil {
				return parsed, true
			}
		}
	}
	return time.Time{}, false
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func simklIDFromCalendarURL(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "simkl.com/")
	raw = strings.TrimPrefix(raw, "www.simkl.com/")
	parts := strings.Split(strings.Trim(raw, "/"), "/")
	for i, part := range parts {
		switch strings.ToLower(part) {
		case "movie", "movies", "tv", "show", "shows", "anime":
			if i+1 < len(parts) {
				if id, err := strconv.Atoi(parts[i+1]); err == nil && id > 0 {
					return id
				}
			}
		}
	}
	return 0
}
