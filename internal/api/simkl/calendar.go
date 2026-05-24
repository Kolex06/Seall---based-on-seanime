package simkl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

const calendarBaseURL = "https://data.simkl.in/calendar"

type CalendarEpisode struct {
	Season  int    `json:"season,omitempty"`
	Episode int    `json:"episode,omitempty"`
	URL     string `json:"url,omitempty"`
}

type CalendarItem struct {
	Kind        MediaType        `json:"-"`
	Title       string           `json:"title,omitempty"`
	Year        int              `json:"year,omitempty"`
	Date        string           `json:"date,omitempty"`
	ReleaseDate string           `json:"release_date,omitempty"`
	Poster      string           `json:"poster,omitempty"`
	Fanart      string           `json:"fanart,omitempty"`
	URL         string           `json:"url,omitempty"`
	IDs         IDs              `json:"ids,omitempty"`
	Episode     *CalendarEpisode `json:"episode,omitempty"`
	AnimeType   string           `json:"anime_type,omitempty"`
}

func DiscoveryMediaFromCalendarItems(items []CalendarItem) []DiscoveryMedia {
	ret := make([]DiscoveryMedia, 0, len(items))
	for _, item := range items {
		title := item.Title
		if title == "" {
			continue
		}

		kind := item.Kind
		if kind == "" {
			kind = kindFromEndpoint("", MediaTypeShows)
		}
		released := firstNonEmpty(item.ReleaseDate, item.Date)
		ids := item.IDs
		mediaURL := normaliseSimklURL(item.URL)
		if ids.PrimarySimklID() == 0 {
			ids.Simkl = simklIDFromURL(mediaURL)
		}

		totalEpisodes := (*int)(nil)
		if kind == MediaTypeMovies {
			totalEpisodes = intPtr(1)
		}

		media := StandardMedia{
			Title:         title,
			Year:          parseReleaseYear(item.Year, released),
			Type:          string(kind),
			Poster:        item.Poster,
			Fanart:        item.Fanart,
			URL:           mediaURL,
			Released:      released,
			AnimeType:     item.AnimeType,
			EnglishName:   title,
			IDs:           ids,
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

func (c *Client) MonthlyCalendar(ctx context.Context, mediaType MediaType, year int, month int) ([]CalendarItem, error) {
	if c == nil {
		return nil, errors.New("simkl: nil client")
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("simkl: invalid calendar month %d", month)
	}

	filename := calendarFileName(mediaType)
	if filename == "" {
		return nil, fmt.Errorf("simkl: unsupported calendar media type %q", mediaType)
	}

	requestURL := fmt.Sprintf("%s/%d/%d/%s", calendarBaseURL, year, month, filename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("simkl: calendar %s failed with %d: %s", requestURL, resp.StatusCode, string(body))
	}

	var items []CalendarItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("simkl: decode calendar %s: %w", requestURL, err)
	}
	for i := range items {
		items[i].Kind = mediaType
	}
	return items, nil
}

func calendarFileName(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeShows:
		return "tv.json"
	case MediaTypeAnime:
		return "anime.json"
	case MediaTypeMovies:
		return "movie_release.json"
	default:
		return ""
	}
}
