package animeofflinedb

import (
	"bufio"
	"errors"
	"net/http"
	"seall/internal/api/mediaapi"
	"seall/internal/library/anime"
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-json"
)

const (
	DatabaseURL = "https://github.com/manami-project/anime-offline-database/releases/download/latest/anime-offline-database.jsonl"
)

type animeEntry struct {
	Sources     []string    `json:"sources"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
	Episodes    int         `json:"episodes"`
	Status      string      `json:"status"`
	AnimeSeason animeSeason `json:"animeSeason"`
	Picture     string      `json:"picture"`
	Thumbnail   string      `json:"thumbnail"`
	Synonyms    []string    `json:"synonyms"`
}

type animeSeason struct {
	Season string `json:"season"`
	Year   int    `json:"year"`
}

const (
	simklPrefix = "https://mediaapi.co/anime/"
	malPrefix   = "https://myanimelist.net/anime/"
)

var (
	normalizedMediaCache   []*anime.NormalizedMedia
	normalizedMediaCacheMu sync.RWMutex
)

// FetchAndConvertDatabase fetches the database and converts entries to NormalizedMedia.
// Only entries with valid SIMKL IDs are included.
// Entries that already exist in existingMediaIDs are excluded.
func FetchAndConvertDatabase(existingMediaIDs map[int]bool) ([]*anime.NormalizedMedia, error) {
	// check cache first
	normalizedMediaCacheMu.RLock()
	if normalizedMediaCache != nil {
		// filter cached results by existingMediaIDs
		result := filterByExistingIDs(normalizedMediaCache, existingMediaIDs)
		normalizedMediaCacheMu.RUnlock()
		return result, nil
	}
	normalizedMediaCacheMu.RUnlock()

	resp, err := http.Get(DatabaseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch database: " + resp.Status)
	}

	// stream and convert directly to NormalizedMedia
	// estimate ~20300 entries with simkl ids
	allMedia := make([]*anime.NormalizedMedia, 0, 20300)
	result := make([]*anime.NormalizedMedia, 0, 20300)

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum == 1 {
			continue // skip metadata line
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// parse entry
		var entry animeEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}

		// convert immediately and discard raw entry
		media := convertEntryToNormalizedMedia(&entry)
		if media == nil {
			continue // no simkl id
		}

		// add to cache (all media with simkl ids)
		allMedia = append(allMedia, media)

		// check if should be included in result
		if existingMediaIDs == nil || !existingMediaIDs[media.ID] {
			result = append(result, media)
		}
		// entry goes out of scope here and can be GC'd
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// cache all media for future calls?
	normalizedMediaCacheMu.Lock()
	normalizedMediaCache = allMedia
	normalizedMediaCacheMu.Unlock()

	return result, nil
}

// filterByExistingIDs filters cached media by existing IDs
func filterByExistingIDs(media []*anime.NormalizedMedia, existingMediaIDs map[int]bool) []*anime.NormalizedMedia {
	if existingMediaIDs == nil || len(existingMediaIDs) == 0 {
		return media
	}

	result := make([]*anime.NormalizedMedia, 0, len(media))
	for _, m := range media {
		if !existingMediaIDs[m.ID] {
			result = append(result, m)
		}
	}
	return result
}

// ClearCache clears the normalized media cache
func ClearCache() {
	normalizedMediaCacheMu.Lock()
	normalizedMediaCache = nil
	normalizedMediaCacheMu.Unlock()
}

// convertEntryToNormalizedMedia converts an animeEntry to NormalizedMedia.
// Returns nil if the entry has no simkl id.
func convertEntryToNormalizedMedia(e *animeEntry) *anime.NormalizedMedia {
	// extract simkl id
	simklID := extractMediaID(e.Sources)
	if simklID == 0 {
		return nil
	}

	malID := extractMALID(e.Sources)
	var malIDPtr *int
	if malID > 0 {
		malIDPtr = &malID
	}

	// convert type to mediaapi.MediaFormat
	var format *mediaapi.MediaFormat
	switch e.Type {
	case "TV":
		f := mediaapi.MediaFormatTv
		format = &f
	case "MOVIE":
		f := mediaapi.MediaFormatMovie
		format = &f
	case "OVA":
		f := mediaapi.MediaFormatOva
		format = &f
	case "ONA":
		f := mediaapi.MediaFormatOna
		format = &f
	case "SPECIAL":
		f := mediaapi.MediaFormatSpecial
		format = &f
	}

	// convert status to mediaapi.MediaStatus
	var status *mediaapi.MediaStatus
	switch e.Status {
	case "FINISHED":
		s := mediaapi.MediaStatusFinished
		status = &s
	case "ONGOING":
		s := mediaapi.MediaStatusReleasing
		status = &s
	case "UPCOMING":
		s := mediaapi.MediaStatusNotYetReleased
		status = &s
	}

	// convert season to mediaapi.MediaSeason
	var season *mediaapi.MediaSeason
	switch e.AnimeSeason.Season {
	case "SPRING":
		s := mediaapi.MediaSeasonSpring
		season = &s
	case "SUMMER":
		s := mediaapi.MediaSeasonSummer
		season = &s
	case "FALL":
		s := mediaapi.MediaSeasonFall
		season = &s
	case "WINTER":
		s := mediaapi.MediaSeasonWinter
		season = &s
	}

	// reuse the same string pointer for all title fields
	title := e.Title
	titleObj := &anime.NormalizedMediaTitle{
		Romaji:        &title,
		English:       &title,
		UserPreferred: &title,
	}

	// build synonyms
	var synonyms []*string
	if len(e.Synonyms) > 0 {
		synonyms = make([]*string, len(e.Synonyms))
		for i := range e.Synonyms {
			synonyms[i] = &e.Synonyms[i]
		}
	}

	// build start date
	var startDate *anime.NormalizedMediaDate
	if e.AnimeSeason.Year > 0 {
		year := e.AnimeSeason.Year
		startDate = &anime.NormalizedMediaDate{
			Year: &year,
		}
	}

	var episodes *int
	if e.Episodes > 0 {
		ep := e.Episodes
		episodes = &ep
	}

	var year *int
	if e.AnimeSeason.Year > 0 {
		y := e.AnimeSeason.Year
		year = &y
	}

	var coverImage *anime.NormalizedMediaCoverImage
	if e.Thumbnail != "" || e.Picture != "" {
		coverImage = &anime.NormalizedMediaCoverImage{
			Large:  &e.Picture,
			Medium: &e.Thumbnail,
		}
	}

	return anime.NewNormalizedMediaFromOfflineDB(
		simklID,
		malIDPtr,
		titleObj,
		synonyms,
		format,
		status,
		season,
		year,
		startDate,
		episodes,
		coverImage,
	)
}

func extractMediaID(sources []string) int {
	for _, source := range sources {
		if strings.HasPrefix(source, simklPrefix) {
			idStr := source[len(simklPrefix):]
			// handle potential trailing slashes or query params
			if idx := strings.IndexAny(idStr, "/?"); idx != -1 {
				idStr = idStr[:idx]
			}
			if id, err := strconv.Atoi(idStr); err == nil {
				return id
			}
		}
	}
	return 0
}

func extractMALID(sources []string) int {
	for _, source := range sources {
		if strings.HasPrefix(source, malPrefix) {
			idStr := source[len(malPrefix):]
			if idx := strings.IndexAny(idStr, "/?"); idx != -1 {
				idStr = idStr[:idx]
			}
			if id, err := strconv.Atoi(idStr); err == nil {
				return id
			}
		}
	}
	return 0
}
