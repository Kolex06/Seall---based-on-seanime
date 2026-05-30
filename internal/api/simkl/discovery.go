package simkl

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DiscoveryMedia struct {
	Kind  MediaType
	Media StandardMedia
}

type restDiscoveryMedia struct {
	Title         string          `json:"title,omitempty"`
	TitleRomaji   string          `json:"title_romaji,omitempty"`
	Year          int             `json:"year,omitempty"`
	EndpointType  string          `json:"endpoint_type,omitempty"`
	Type          string          `json:"type,omitempty"`
	Poster        string          `json:"poster,omitempty"`
	Fanart        string          `json:"fanart,omitempty"`
	URL           string          `json:"url,omitempty"`
	ReleaseDate   string          `json:"release_date,omitempty"`
	Released      string          `json:"released,omitempty"`
	Status        string          `json:"status,omitempty"`
	Overview      string          `json:"overview,omitempty"`
	Genres        []string        `json:"genres,omitempty"`
	AnimeType     string          `json:"anime_type,omitempty"`
	Country       string          `json:"country,omitempty"`
	IDs           restDiscoveryID `json:"ids,omitempty"`
	Ratings       Ratings         `json:"ratings,omitempty"`
	Runtime       interface{}     `json:"runtime,omitempty"`
	TotalEpisodes int             `json:"total_episodes,omitempty"`
}

type restDiscoveryID struct {
	Simkl       interface{} `json:"simkl,omitempty"`
	SimklID     interface{} `json:"simkl_id,omitempty"`
	Slug        string      `json:"slug,omitempty"`
	IMDB        interface{} `json:"imdb,omitempty"`
	TMDB        interface{} `json:"tmdb,omitempty"`
	TVDB        interface{} `json:"tvdb,omitempty"`
	MAL         interface{} `json:"mal,omitempty"`
	AniDB       interface{} `json:"anidb,omitempty"`
	Kitsu       interface{} `json:"kitsu,omitempty"`
	Crunchyroll interface{} `json:"crunchyroll,omitempty"`
}

func (c *Client) SearchMedia(ctx context.Context, mediaType MediaType, query string) ([]DiscoveryMedia, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	q := url.Values{}
	q.Set("q", query)

	var ret []restDiscoveryMedia
	if err := c.doJSON(ctx, http.MethodGet, searchMediaPath(mediaType), q, nil, &ret, false, false); err != nil {
		return nil, err
	}

	return toDiscoveryMedia(mediaType, ret), nil
}

func (c *Client) TrendingMedia(ctx context.Context, mediaType MediaType, page int, limit int) ([]DiscoveryMedia, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	q := url.Values{}
	q.Set("extended", "overview")
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(limit))

	var ret []restDiscoveryMedia
	if err := c.doJSON(ctx, http.MethodGet, trendingMediaPath(mediaType), q, nil, &ret, false, false); err != nil {
		return nil, err
	}

	return toDiscoveryMedia(mediaType, ret), nil
}

func searchMediaPath(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeMovies:
		return "/search/movie"
	case MediaTypeShows:
		return "/search/tv"
	case MediaTypeAnime:
		return "/search/anime"
	default:
		return "/search/tv"
	}
}

func trendingMediaPath(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeMovies:
		return "/movies/trending/"
	case MediaTypeShows:
		return "/tv/trending/"
	case MediaTypeAnime:
		return "/anime/trending/"
	default:
		return "/tv/trending/"
	}
}

func toDiscoveryMedia(fallback MediaType, items []restDiscoveryMedia) []DiscoveryMedia {
	ret := make([]DiscoveryMedia, 0, len(items))
	for _, item := range items {
		title := firstNonEmpty(item.Title, item.TitleRomaji)
		if title == "" {
			continue
		}

		kind := kindFromEndpoint(firstNonEmpty(item.EndpointType, item.Type), fallback)
		released := firstNonEmpty(item.Released, item.ReleaseDate)
		year := parseReleaseYear(item.Year, released)
		runtime := parseRuntimeMinutes(item.Runtime)
		ids := item.IDs.toIDs()
		totalEpisodes := intPtrIfPositive(item.TotalEpisodes)
		mediaURL := normaliseSimklURL(item.URL)
		if ids.PrimarySimklID() == 0 {
			ids.Simkl = simklIDFromURL(mediaURL)
		}

		if kind == MediaTypeMovies && totalEpisodes == nil {
			totalEpisodes = intPtr(1)
		}

		media := StandardMedia{
			Title:         title,
			Year:          year,
			Type:          string(kind),
			Status:        item.Status,
			Poster:        item.Poster,
			Fanart:        item.Fanart,
			URL:           mediaURL,
			Runtime:       runtime,
			Overview:      item.Overview,
			Genres:        item.Genres,
			Country:       item.Country,
			Released:      released,
			AnimeType:     item.AnimeType,
			EnglishName:   title,
			IDs:           ids,
			Ratings:       item.Ratings,
			TotalEpisodes: totalEpisodes,
		}
		cacheDiscoveryMedia(kind, media)

		ret = append(ret, DiscoveryMedia{
			Kind:  kind,
			Media: media,
		})
	}
	return ret
}

func kindFromEndpoint(endpoint string, fallback MediaType) MediaType {
	switch strings.ToLower(strings.TrimSpace(endpoint)) {
	case "movie", "movies":
		return MediaTypeMovies
	case "show", "shows", "tv":
		return MediaTypeShows
	case "anime":
		return MediaTypeAnime
	default:
		if fallback != MediaTypeAll && fallback != "" {
			return fallback
		}
		return MediaTypeShows
	}
}

func parseReleaseYear(year int, released string) int {
	if year != 0 {
		return year
	}
	released = strings.TrimSpace(released)
	if len(released) < 4 {
		return 0
	}
	if parsed, err := strconv.Atoi(released[:4]); err == nil {
		return parsed
	}
	if parts := strings.Split(released, "/"); len(parts) == 3 {
		if parsed, err := strconv.Atoi(parts[2]); err == nil {
			return parsed
		}
	}
	if parts := strings.Split(released, "-"); len(parts) == 3 {
		if parsed, err := strconv.Atoi(parts[0]); err == nil {
			return parsed
		}
	}
	return 0
}

func parseRuntimeMinutes(raw interface{}) *int {
	switch v := raw.(type) {
	case nil:
		return nil
	case int:
		return intPtrIfPositive(v)
	case int64:
		return intPtrIfPositive(int(v))
	case float64:
		return intPtrIfPositive(int(v))
	case string:
		value := strings.TrimSpace(strings.ToLower(v))
		if value == "" {
			return nil
		}
		if parsed, err := strconv.Atoi(value); err == nil {
			return intPtrIfPositive(parsed)
		}

		total := 0
		for _, part := range strings.Fields(value) {
			switch {
			case strings.HasSuffix(part, "h"):
				if parsed, err := strconv.Atoi(strings.TrimSuffix(part, "h")); err == nil {
					total += parsed * 60
				}
			case strings.HasSuffix(part, "m"):
				if parsed, err := strconv.Atoi(strings.TrimSuffix(part, "m")); err == nil {
					total += parsed
				}
			}
		}
		return intPtrIfPositive(total)
	default:
		return nil
	}
}

func (ids restDiscoveryID) toIDs() IDs {
	return IDs{
		Simkl:       flexibleInt(ids.Simkl),
		SimklID:     flexibleInt(ids.SimklID),
		Slug:        ids.Slug,
		IMDB:        flexibleString(ids.IMDB),
		TMDB:        flexibleString(ids.TMDB),
		TVDB:        flexibleString(ids.TVDB),
		MAL:         flexibleString(ids.MAL),
		AniDB:       flexibleString(ids.AniDB),
		Kitsu:       flexibleString(ids.Kitsu),
		Crunchyroll: flexibleString(ids.Crunchyroll),
	}
}

func flexibleInt(raw interface{}) int {
	switch v := raw.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		parsed, _ := strconv.Atoi(v)
		return parsed
	default:
		return 0
	}
}

func flexibleString(raw interface{}) string {
	switch v := raw.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func normaliseSimklURL(raw string) string {
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	if strings.HasPrefix(raw, "/") {
		return "https://simkl.com" + raw
	}
	return "https://simkl.com/" + strings.TrimLeft(raw, "/")
}

func simklIDFromURL(raw string) int {
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

func intPtrIfPositive(v int) *int {
	if v <= 0 {
		return nil
	}
	return &v
}
