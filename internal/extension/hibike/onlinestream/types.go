package hibikeonlinestream

import "strings"

type (
	Provider interface {
		Search(opts SearchOptions) ([]*SearchResult, error)
		// FindMediaItems returns playable items for the given provider media ID.
		FindMediaItems(id string) ([]*EpisodeDetails, error)
		// FindMediaItemServer returns the stream server for the given item.
		// The "server" argument can be "default"
		FindMediaItemServer(item *EpisodeDetails, server string) (*EpisodeServer, error)
		// GetSettings returns the provider settings.
		GetSettings() Settings
	}

	SearchOptions struct {
		// The media object provided by Seall.
		Media Media `json:"media"`
		// The search query.
		Query string `json:"query"`
		// Whether to search for subbed/original or dubbed media.
		Dub bool `json:"dub"`
		// The year the media was released.
		// Will be 0 if the year is not available.
		Year int `json:"year"`
		// SIMKL media type: "movies", "shows", "anime" or "all".
		MediaType string `json:"mediaType,omitempty"`
	}

	Media struct {
		// SIMKL ID of the media.
		ID int `json:"id"`
		// MyAnimeList ID, when SIMKL has one for this media.
		IDMal *int `json:"idMal,omitempty"`
		// SIMKL URL for the media.
		SiteURL *string `json:"siteUrl,omitempty"`
		// SIMKL media type: "movies", "shows", "anime" or "all".
		MediaType string `json:"mediaType,omitempty"`
		// e.g. "FINISHED", "RELEASING", "NOT_YET_RELEASED", "CANCELLED", "HIATUS"
		// This will be set to "NOT_YET_RELEASED" if the status is unknown.
		Status string `json:"status,omitempty"`
		// e.g. "TV", "TV_SHORT", "MOVIE", "SPECIAL", "OVA", "ONA", "MUSIC"
		// This will be set to "TV" if the format is unknown.
		Format string `json:"format,omitempty"`
		// e.g. "Attack on Titan"
		// This will be undefined if the english title is unknown.
		EnglishTitle *string `json:"englishTitle,omitempty"`
		// Main title.
		RomajiTitle string `json:"romajiTitle,omitempty"`
		// TotalEpisodes is total number of episodes of the media.
		// This will be -1 if the total number of episodes is unknown / not applicable.
		EpisodeCount int `json:"episodeCount,omitempty"`
		// All alternative titles of the media.
		Synonyms []string `json:"synonyms"`
		// Whether the media is NSFW.
		IsAdult bool `json:"isAdult"`
		// Start date of the media.
		// This will be undefined if it has no start date.
		StartDate *FuzzyDate `json:"startDate,omitempty"`
	}

	FuzzyDate struct {
		Year  int  `json:"year"`
		Month *int `json:"month"`
		Day   *int `json:"day"`
	}

	Settings struct {
		EpisodeServers      []string `json:"episodeServers"`
		SupportsDub         bool     `json:"supportsDub"`
		SupportedMediaTypes []string `json:"supportedMediaTypes"`
	}

	SearchResult struct {
		// ID is the provider media slug or ID.
		// It is used to fetch playable items.
		ID string `json:"id"`
		// Title is the media title.
		Title string `json:"title"`
		// URL is the provider page URL.
		URL      string   `json:"url"`
		SubOrDub SubOrDub `json:"subOrDub"`
	}

	// EpisodeDetails contains the playable item information from a provider.
	// It is obtained by scraping the provider item list.
	EpisodeDetails struct {
		// "ID" of the extension.
		Provider string `json:"provider"`
		// ID is the provider item slug.
		// e.g. "the-apothecary-diaries-18578".
		ID string `json:"id"`
		// Item number.
		// From 0 to n.
		Number int `json:"number"`
		// Provider item page URL.
		URL string `json:"url"`
		// Item title.
		// Leave it empty if the title is not available.
		Title string `json:"title,omitempty"`
	}

	// EpisodeServer contains the server, headers and video sources for an item.
	EpisodeServer struct {
		// "ID" of the extension.
		Provider string `json:"provider"`
		// Stream server name.
		// e.g. "vidcloud".
		Server string `json:"server"`
		// HTTP headers for the video request.
		Headers map[string]string `json:"headers"`
		// Video sources for the episode.
		VideoSources []*VideoSource `json:"videoSources"`
	}

	SubOrDub string

	VideoSourceType string

	VideoSource struct {
		// URL of the video source.
		URL string `json:"url"`
		// Type of the video source.
		Type VideoSourceType `json:"type"`
		// Label of the video source. (e.g. "English")
		Label string `json:"label,omitempty"`
		// Quality of the video source.
		// e.g. "default", "auto", "1080p".
		Quality string `json:"quality"`
		// Subtitles for the video source.
		Subtitles []*VideoSubtitle `json:"subtitles"`
	}

	VideoSubtitle struct {
		ID  string `json:"id"`
		URL string `json:"url"`
		// e.g. "en", "fr"
		Language  string `json:"language"`
		IsDefault bool   `json:"isDefault"`
	}

	VideoExtractor interface {
		Extract(uri string) ([]*VideoSource, error)
	}
)

const (
	MediaTypeAll    = "all"
	MediaTypeMovies = "movies"
	MediaTypeShows  = "shows"
	MediaTypeAnime  = "anime"
)

const (
	Sub       SubOrDub = "sub"
	Dub       SubOrDub = "dub"
	SubAndDub SubOrDub = "both"
)

const (
	VideoSourceMP4     VideoSourceType = "mp4"
	VideoSourceM3U8    VideoSourceType = "m3u8"
	VideoSourceUnknown VideoSourceType = "unknown"
)

func NormalizeMediaType(mediaType string) string {
	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "", "all":
		return MediaTypeAll
	case "movie", "movies":
		return MediaTypeMovies
	case "show", "shows", "tv", "series":
		return MediaTypeShows
	case "anime":
		return MediaTypeAnime
	default:
		return strings.ToLower(strings.TrimSpace(mediaType))
	}
}

func NormalizeSupportedMediaTypes(settings Settings) []string {
	if len(settings.SupportedMediaTypes) == 0 {
		return []string{MediaTypeAnime}
	}

	seen := make(map[string]struct{}, len(settings.SupportedMediaTypes))
	ret := make([]string, 0, len(settings.SupportedMediaTypes))
	for _, mediaType := range settings.SupportedMediaTypes {
		normalized := NormalizeMediaType(mediaType)
		if normalized == "" {
			continue
		}
		if _, found := seen[normalized]; found {
			continue
		}
		seen[normalized] = struct{}{}
		ret = append(ret, normalized)
	}
	if len(ret) == 0 {
		return []string{MediaTypeAnime}
	}
	return ret
}

func SupportsMediaType(settings Settings, mediaType string) bool {
	mediaType = NormalizeMediaType(mediaType)
	if mediaType == MediaTypeAll || mediaType == "" {
		return true
	}
	for _, supported := range NormalizeSupportedMediaTypes(settings) {
		if supported == MediaTypeAll || supported == mediaType {
			return true
		}
	}
	return false
}
