package onlinestream

import (
	"fmt"
	"seall/internal/api/mediaapi"
	hibikeonlinestream "seall/internal/extension/hibike/onlinestream"
	"strings"
)

func MediaTypeForBaseMedia(media *mediaapi.BaseAnime) string {
	if media == nil {
		return hibikeonlinestream.MediaTypeAnime
	}

	if siteURL := media.GetSiteURL(); siteURL != nil {
		normalizedURL := strings.ToLower(*siteURL)
		switch {
		case strings.Contains(normalizedURL, "simkl.com/movie/"):
			return hibikeonlinestream.MediaTypeMovies
		case strings.Contains(normalizedURL, "simkl.com/tv/"):
			return hibikeonlinestream.MediaTypeShows
		case strings.Contains(normalizedURL, "simkl.com/anime/"):
			return hibikeonlinestream.MediaTypeAnime
		}
	}

	if media.GetFormat() != nil && *media.GetFormat() == mediaapi.MediaFormatMovie && media.GetIDMal() == nil {
		return hibikeonlinestream.MediaTypeMovies
	}

	return hibikeonlinestream.MediaTypeAnime
}

func mediaTypeLabel(mediaType string) string {
	switch hibikeonlinestream.NormalizeMediaType(mediaType) {
	case hibikeonlinestream.MediaTypeMovies:
		return "movies"
	case hibikeonlinestream.MediaTypeShows:
		return "shows"
	case hibikeonlinestream.MediaTypeAnime:
		return "anime"
	default:
		return "media"
	}
}

func queryMediaFromBaseMedia(media *mediaapi.BaseAnime, mediaType string) hibikeonlinestream.Media {
	if media == nil {
		return hibikeonlinestream.Media{MediaType: mediaType}
	}

	var status string
	if media.GetStatus() != nil {
		status = string(*media.GetStatus())
	}

	var format string
	if media.GetFormat() != nil {
		format = string(*media.GetFormat())
	}

	var startDate *hibikeonlinestream.FuzzyDate
	if media.GetStartDate() != nil {
		startDate = &hibikeonlinestream.FuzzyDate{
			Year:  0,
			Month: media.GetStartDate().GetMonth(),
			Day:   media.GetStartDate().GetDay(),
		}
		if media.GetStartDate().GetYear() != nil {
			startDate.Year = *media.GetStartDate().GetYear()
		}
	}

	return hibikeonlinestream.Media{
		ID:           media.ID,
		IDMal:        media.GetIDMal(),
		SiteURL:      media.GetSiteURL(),
		MediaType:    mediaType,
		Status:       status,
		Format:       format,
		EnglishTitle: media.GetTitle().GetEnglish(),
		RomajiTitle:  media.GetRomajiTitleSafe(),
		EpisodeCount: media.GetTotalEpisodeCount(),
		Synonyms:     media.GetSynonymsContainingSeason(),
		IsAdult:      media.GetIsAdult() != nil && *media.GetIsAdult(),
		StartDate:    startDate,
	}
}

func unsupportedProviderMediaTypeError(provider string, supported []string, mediaType string) error {
	return fmt.Errorf("provider '%s' supports %s, not %s", provider, strings.Join(supported, ", "), mediaTypeLabel(mediaType))
}
