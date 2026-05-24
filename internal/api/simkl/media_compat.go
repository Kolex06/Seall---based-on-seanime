package simkl

import (
	"fmt"
	"seall/internal/api/mediaapi"
	"strconv"
	"strings"
	"time"
)

func ToSIMKLStatus(status WatchStatus) mediaapi.MediaListStatus {
	switch status {
	case WatchStatusWatching:
		return mediaapi.MediaListStatusCurrent
	case WatchStatusPlanToWatch:
		return mediaapi.MediaListStatusPlanning
	case WatchStatusCompleted:
		return mediaapi.MediaListStatusCompleted
	case WatchStatusDropped:
		return mediaapi.MediaListStatusDropped
	case WatchStatusHold:
		return mediaapi.MediaListStatusPaused
	default:
		return mediaapi.MediaListStatusPlanning
	}
}

func FromSIMKLStatus(status *mediaapi.MediaListStatus) WatchStatus {
	if status == nil {
		return WatchStatusPlanToWatch
	}
	switch *status {
	case mediaapi.MediaListStatusCurrent:
		return WatchStatusWatching
	case mediaapi.MediaListStatusCompleted:
		return WatchStatusCompleted
	case mediaapi.MediaListStatusDropped:
		return WatchStatusDropped
	case mediaapi.MediaListStatusPaused:
		return WatchStatusHold
	case mediaapi.MediaListStatusPlanning:
		return WatchStatusPlanToWatch
	default:
		return WatchStatusPlanToWatch
	}
}

func ToSIMKLAnimeCollection(items *AllItems) *mediaapi.AnimeCollection {
	collection := &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: make([]*mediaapi.AnimeCollection_MediaListCollection_Lists, 0),
		},
	}
	if items == nil {
		return collection
	}

	lists := map[mediaapi.MediaListStatus]*mediaapi.AnimeCollection_MediaListCollection_Lists{}
	add := func(kind MediaType, item WatchlistItem) {
		status := ToSIMKLStatus(item.Status)
		list := lists[status]
		if list == nil {
			name := statusListName(status)
			list = &mediaapi.AnimeCollection_MediaListCollection_Lists{
				Name:    &name,
				Status:  &status,
				Entries: make([]*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries, 0),
			}
			lists[status] = list
			collection.MediaListCollection.Lists = append(collection.MediaListCollection.Lists, list)
		}

		entry := ToSIMKLAnimeEntry(kind, item)
		if entry != nil {
			list.Entries = append(list.Entries, entry)
		}
	}

	for _, item := range items.Movies {
		add(MediaTypeMovies, item)
	}
	for _, item := range items.Shows {
		add(MediaTypeShows, item)
	}
	for _, item := range items.Anime {
		add(MediaTypeAnime, item)
	}

	return collection
}

func ToSIMKLAnimeEntry(kind MediaType, item WatchlistItem) *mediaapi.AnimeCollection_MediaListCollection_Lists_Entries {
	simklMedia := item.Media()
	media := ToBaseAnime(kind, simklMedia)
	if media == nil {
		return nil
	}
	status := ToSIMKLStatus(item.Status)
	progress := itemProgress(kind, item)
	score := scoreFromSIMKL(item.UserRating)
	completedAt := firstNonEmpty(item.LastWatchedAt, item.LastWatched)
	if simklMedia != nil {
		completedAt = firstNonEmpty(completedAt, simklMedia.WatchedAt)
	}

	return &mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{
		ID:          media.ID,
		Media:       media,
		Progress:    progress,
		Score:       score,
		StartedAt:   simklStartedAt(item.AddedToWatchlistAt),
		CompletedAt: simklCompletedAt(completedAt),
		Status:      &status,
	}
}

func (item WatchlistItem) Media() *StandardMedia {
	if item.Movie != nil {
		return item.Movie
	}
	return item.Show
}

func ToBaseAnime(kind MediaType, media *StandardMedia) *mediaapi.BaseAnime {
	if media == nil {
		return nil
	}
	cacheDiscoveryMedia(kind, *media)
	kind = KindFromStandardMedia(kind, media)

	id := media.IDs.PrimarySimklID()
	if id == 0 {
		id = simklIDFromURL(media.URL)
	}
	if id == 0 {
		id = stableFallbackID(kind, media.Title, media.Year)
	}
	title := media.Title
	if media.EnglishName != "" {
		title = media.EnglishName
	}
	siteURL := simklSiteURL(kind, media)
	mediaType := mediaapi.MediaTypeAnime
	format := toSIMKLFormat(kind, media)
	status := toSIMKLReleaseStatus(media)
	poster := ImageURL(ImageKindPoster, media.Poster, ImageSizePosterMedium)
	banner := ImageURL(ImageKindFanart, media.Fanart, ImageSizeFanartWide)
	meanScore := simklMeanScore(media)
	idMal := stringInt(media.IDs.MAL)
	genres := stringSlicePointers(media.Genres)
	episodes := media.TotalEpisodes
	if kind == MediaTypeMovies && episodes == nil {
		episodes = intPtr(1)
	}

	ret := &mediaapi.BaseAnime{
		ID:          id,
		IDMal:       idMal,
		SiteURL:     &siteURL,
		Type:        &mediaType,
		Format:      &format,
		Status:      &status,
		SeasonYear:  yearPtr(media.Year),
		BannerImage: emptyToNil(banner),
		Episodes:    episodes,
		Description: emptyToNil(media.Overview),
		Genres:      genres,
		Duration:    media.Runtime,
		MeanScore:   meanScore,
		Title: &mediaapi.BaseAnime_Title{
			English:       emptyToNil(title),
			Romaji:        emptyToNil(title),
			UserPreferred: emptyToNil(title),
		},
		CoverImage: &mediaapi.BaseAnime_CoverImage{
			ExtraLarge: emptyToNil(poster),
			Large:      emptyToNil(poster),
			Medium:     emptyToNil(poster),
		},
		StartDate: simklStartDate(media.Released, media.Year),
	}

	return ret
}

func ToCompleteAnime(kind MediaType, media *StandardMedia) *mediaapi.CompleteAnime {
	base := ToBaseAnime(kind, media)
	if base == nil {
		return nil
	}

	return &mediaapi.CompleteAnime{
		ID:          base.ID,
		IDMal:       base.IDMal,
		SiteURL:     base.SiteURL,
		Status:      base.Status,
		SeasonYear:  base.SeasonYear,
		Type:        base.Type,
		Format:      base.Format,
		BannerImage: base.BannerImage,
		Episodes:    base.Episodes,
		Synonyms:    base.Synonyms,
		IsAdult:     base.IsAdult,
		MeanScore:   base.MeanScore,
		Description: base.Description,
		Genres:      base.Genres,
		Duration:    base.Duration,
		Title: &mediaapi.CompleteAnime_Title{
			English:       base.Title.GetEnglish(),
			Native:        base.Title.GetNative(),
			Romaji:        base.Title.GetRomaji(),
			UserPreferred: base.Title.GetUserPreferred(),
		},
		CoverImage: &mediaapi.CompleteAnime_CoverImage{
			Color:      base.CoverImage.GetColor(),
			ExtraLarge: base.CoverImage.GetExtraLarge(),
			Large:      base.CoverImage.GetLarge(),
			Medium:     base.CoverImage.GetMedium(),
		},
		StartDate: &mediaapi.CompleteAnime_StartDate{
			Day:   base.StartDate.GetDay(),
			Month: base.StartDate.GetMonth(),
			Year:  base.StartDate.GetYear(),
		},
	}
}

func ToAnimeDetails(kind MediaType, media *StandardMedia) *mediaapi.AnimeDetailsById_Media {
	base := ToBaseAnime(kind, media)
	if base == nil {
		return nil
	}
	return &mediaapi.AnimeDetailsById_Media{
		ID:           base.ID,
		SiteURL:      base.SiteURL,
		Description:  base.Description,
		Duration:     base.Duration,
		Genres:       base.Genres,
		MeanScore:    base.MeanScore,
		AverageScore: base.MeanScore,
		StartDate: &mediaapi.AnimeDetailsById_Media_StartDate{
			Day:   base.StartDate.GetDay(),
			Month: base.StartDate.GetMonth(),
			Year:  base.StartDate.GetYear(),
		},
	}
}

func ToSimklMedia(kind MediaType, id int) StandardMedia {
	ids := IDs{Simkl: id}
	return StandardMedia{IDs: ids, Type: string(kind)}
}

func toSIMKLFormat(kind MediaType, media *StandardMedia) mediaapi.MediaFormat {
	if kind == MediaTypeMovies || strings.EqualFold(media.AnimeType, "movie") {
		return mediaapi.MediaFormatMovie
	}
	switch strings.ToLower(media.AnimeType) {
	case "special":
		return mediaapi.MediaFormatSpecial
	case "ova":
		return mediaapi.MediaFormatOva
	case "ona":
		return mediaapi.MediaFormatOna
	case "music video":
		return mediaapi.MediaFormatMusic
	default:
		return mediaapi.MediaFormatTv
	}
}

func KindFromStandardMedia(fallback MediaType, media *StandardMedia) MediaType {
	if media != nil {
		switch strings.ToLower(strings.TrimSpace(media.Type)) {
		case "movie", "movies":
			return MediaTypeMovies
		case "show", "shows", "tv", "series":
			return MediaTypeShows
		case "anime":
			return MediaTypeAnime
		}
		switch strings.ToLower(strings.TrimSpace(media.AnimeType)) {
		case "movie":
			return MediaTypeMovies
		case "tv", "special", "ova", "ona", "music video":
			if fallback == MediaTypeAnime {
				return MediaTypeAnime
			}
		}
	}
	if fallback != MediaTypeAll && fallback != "" {
		return fallback
	}
	return MediaTypeShows
}

func toSIMKLReleaseStatus(media *StandardMedia) mediaapi.MediaStatus {
	if media == nil {
		return mediaapi.MediaStatusFinished
	}
	status := strings.ToLower(strings.TrimSpace(media.Status))
	switch status {
	case "cancelled", "canceled":
		return mediaapi.MediaStatusCancelled
	case "hiatus":
		return mediaapi.MediaStatusHiatus
	case "upcoming", "not yet released", "not released", "tba", "announced", "planned", "pre-production", "post-production", "in production":
		return mediaapi.MediaStatusNotYetReleased
	case "returning series", "releasing", "ongoing":
		return mediaapi.MediaStatusReleasing
	}

	if releaseDate, ok := simklReleaseDate(media.Released); ok && releaseDate.After(time.Now().UTC()) {
		return mediaapi.MediaStatusNotYetReleased
	}
	if media.Year > time.Now().UTC().Year() {
		return mediaapi.MediaStatusNotYetReleased
	}
	return mediaapi.MediaStatusFinished
}

func simklStartDate(raw string, fallbackYear int) *mediaapi.BaseAnime_StartDate {
	year, month, day, ok := simklDateParts(raw)
	if !ok {
		year = fallbackYear
	}
	return &mediaapi.BaseAnime_StartDate{
		Day:   positiveIntPtr(day),
		Month: positiveIntPtr(month),
		Year:  yearPtr(year),
	}
}

func simklReleaseDate(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}

func positiveIntPtr(value int) *int {
	if value <= 0 {
		return nil
	}
	return &value
}

func simklMeanScore(media *StandardMedia) *int {
	if media == nil || media.Ratings.Simkl == nil || media.Ratings.Simkl.Rating == 0 {
		return nil
	}
	score := int(media.Ratings.Simkl.Rating * 10)
	return &score
}

func simklSiteURL(kind MediaType, media *StandardMedia) string {
	if media.URL != "" {
		if strings.HasPrefix(media.URL, "http://") || strings.HasPrefix(media.URL, "https://") {
			return media.URL
		}
		if strings.HasPrefix(media.URL, "/") {
			return "https://simkl.com" + media.URL
		}
		return "https://simkl.com/" + strings.TrimLeft(media.URL, "/")
	}
	id := media.IDs.PrimarySimklID()
	slug := media.IDs.Slug
	switch kind {
	case MediaTypeMovies:
		if slug != "" {
			return fmt.Sprintf("https://simkl.com/movie/%d/%s", id, slug)
		}
		return fmt.Sprintf("https://simkl.com/movie/%d", id)
	case MediaTypeAnime:
		if slug != "" {
			return fmt.Sprintf("https://simkl.com/anime/%d/%s", id, slug)
		}
		return fmt.Sprintf("https://simkl.com/anime/%d", id)
	default:
		if slug != "" {
			return fmt.Sprintf("https://simkl.com/tv/%d/%s", id, slug)
		}
		return fmt.Sprintf("https://simkl.com/tv/%d", id)
	}
}

func statusListName(status mediaapi.MediaListStatus) string {
	switch status {
	case mediaapi.MediaListStatusCurrent:
		return "Watching"
	case mediaapi.MediaListStatusPlanning:
		return "Plan to watch"
	case mediaapi.MediaListStatusCompleted:
		return "Completed"
	case mediaapi.MediaListStatusDropped:
		return "Dropped"
	case mediaapi.MediaListStatusPaused:
		return "On hold"
	default:
		return string(status)
	}
}

func itemProgress(kind MediaType, item WatchlistItem) *int {
	if kind == MediaTypeMovies {
		if item.Status == WatchStatusCompleted {
			return intPtr(1)
		}
		return intPtr(0)
	}
	if item.WatchedEpisodesCount != nil {
		return item.WatchedEpisodesCount
	}
	return intPtr(0)
}

func ScoreToSIMKL(score *int) *int {
	if score == nil {
		return nil
	}
	value := *score
	if value <= 0 {
		return intPtr(0)
	}
	if value > 100 {
		value = 100
	}
	if value <= 10 {
		return intPtr(value)
	}
	value = (value + 5) / 10
	return &value
}

func scoreFromSIMKL(score *int) *float64 {
	if score == nil {
		return nil
	}
	value := *score
	if value <= 10 {
		value *= 10
	}
	if value > 100 {
		value = 100
	}
	f := float64(value)
	return &f
}

func simklStartedAt(raw string) *mediaapi.AnimeCollection_MediaListCollection_Lists_Entries_StartedAt {
	year, month, day, ok := simklDateParts(raw)
	if !ok {
		return nil
	}
	return &mediaapi.AnimeCollection_MediaListCollection_Lists_Entries_StartedAt{
		Year:  &year,
		Month: &month,
		Day:   &day,
	}
}

func simklCompletedAt(raw string) *mediaapi.AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt {
	year, month, day, ok := simklDateParts(raw)
	if !ok {
		return nil
	}
	return &mediaapi.AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt{
		Year:  &year,
		Month: &month,
		Day:   &day,
	}
}

func simklDateParts(raw string) (int, int, int, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, 0, 0, false
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed.Year(), int(parsed.Month()), parsed.Day(), true
		}
	}
	return 0, 0, 0, false
}

func yearPtr(year int) *int {
	if year == 0 {
		return nil
	}
	return &year
}

func intPtr(v int) *int {
	return &v
}

func stringInt(raw string) *int {
	if raw == "" {
		return nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	return &v
}

func stringSlicePointers(values []string) []*string {
	if len(values) == 0 {
		return nil
	}
	ret := make([]*string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		v := value
		ret = append(ret, &v)
	}
	return ret
}

func emptyToNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func stableFallbackID(kind MediaType, title string, year int) int {
	raw := string(kind) + ":" + title + ":" + strconv.Itoa(year)
	var hash uint32 = 2166136261
	for _, b := range []byte(raw) {
		hash ^= uint32(b)
		hash *= 16777619
	}
	return int(hash & 0x7fffffff)
}
