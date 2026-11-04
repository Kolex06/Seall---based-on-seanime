package metadata_provider

import (
	"context"
	"errors"
	"fmt"
	"seall/internal/api/animap"
	"seall/internal/api/metadata"
	"seall/internal/api/simkl"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (p *ProviderImpl) fetchSimklEpisodeMetadata(ctx context.Context, mId int) (*metadata.AnimeMetadata, error) {
	if p.simklClientRef == nil || p.simklClientRef.IsAbsent() {
		return nil, errors.New("simkl client is unavailable")
	}

	client := p.simklClientRef.Get()
	if cached, ok := simkl.CachedDiscoveryMedia(mId); ok && (cached.Kind == simkl.MediaTypeMovies || strings.EqualFold(cached.Media.AnimeType, "movie")) {
		media := cached.Media
		if details, err := client.MediaDetails(ctx, cached.Kind, strconv.Itoa(mId), "full"); err == nil && details != nil {
			media = *details
		}
		return buildSimklMovieMetadata(mId, &media), nil
	}
	var lastErr error
	for _, kind := range simklEpisodeKindCandidates(mId) {
		episodes, err := client.MediaEpisodes(ctx, kind, strconv.Itoa(mId), "full")
		if err != nil {
			lastErr = err
			continue
		}
		if len(episodes) == 0 {
			lastErr = fmt.Errorf("simkl: no episode metadata returned for %d", mId)
			continue
		}
		if kind == simkl.MediaTypeShows {
			return buildSimklShowMetadata(mId, episodes), nil
		}
		return buildSimklAnimeMetadata(mId, episodes), nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("simkl: no episode metadata returned for %d", mId)
	}
	return nil, lastErr
}

func buildSimklMovieMetadata(mId int, media *simkl.StandardMedia) *metadata.AnimeMetadata {
	base := simkl.ToBaseAnime(simkl.MediaTypeMovies, media)
	ret := buildSimklAnimeMetadata(mId, []simkl.Episode{{Episode: 1, Title: base.GetPreferredTitle(), Description: media.Overview, Date: media.Released}})
	ret.Titles["en"] = base.GetPreferredTitle()
	ret.Mappings.Type = "movie"
	ret.Mappings.ImdbId = media.IDs.IMDB
	ret.Mappings.ThemoviedbId = media.IDs.TMDB
	ep := ret.Episodes["1"]
	ep.Image = simkl.ImageURL(simkl.ImageKindFanart, media.Fanart, simkl.ImageSizeFanartWide)
	if ep.Image == "" {
		ep.Image = base.GetCoverImageSafe()
	}
	ep.HasImage = ep.Image != ""
	if media.Runtime != nil {
		ep.Length = *media.Runtime
	}
	return ret
}

func simklEpisodeKindCandidates(mId int) []simkl.MediaType {
	ret := make([]simkl.MediaType, 0, 2)
	add := func(kind simkl.MediaType) {
		if kind != simkl.MediaTypeAnime && kind != simkl.MediaTypeShows {
			return
		}
		for _, existing := range ret {
			if existing == kind {
				return
			}
		}
		ret = append(ret, kind)
	}

	if cached, ok := simkl.CachedDiscoveryMedia(mId); ok {
		add(cached.Kind)
	}
	add(simkl.MediaTypeAnime)
	add(simkl.MediaTypeShows)
	return ret
}

func buildSimklAnimeMetadata(mId int, episodes []simkl.Episode) *metadata.AnimeMetadata {
	ret := &metadata.AnimeMetadata{
		Titles:       make(map[string]string),
		Episodes:     make(map[string]*metadata.EpisodeMetadata),
		EpisodeCount: 0,
		SpecialCount: 0,
		Mappings: &metadata.AnimeMappings{
			MediaId: mId,
		},
	}

	absoluteEpisodeNumber := 0
	for _, ep := range episodes {
		episodeNumber := simklEpisodeNumber(ep)
		if episodeNumber == 0 {
			continue
		}

		isSpecial := strings.EqualFold(ep.Type, "special")
		key := strconv.Itoa(absoluteEpisodeNumber + 1)
		displayEpisodeNumber := absoluteEpisodeNumber + 1
		if isSpecial {
			key = "S" + strconv.Itoa(episodeNumber)
			displayEpisodeNumber = episodeNumber
			ret.SpecialCount++
		} else {
			absoluteEpisodeNumber++
			ret.EpisodeCount++
		}

		image := simkl.ImageURL(simkl.ImageKindEpisode, ep.Image, simkl.ImageSizeEpisodeWide)
		overview := strings.ReplaceAll(ep.Description, "`", "'")
		seasonNumber := ep.Season
		seasonEpisodeNumber := episodeNumber
		if seasonNumber == 0 && ep.TVDB != nil {
			seasonNumber = ep.TVDB.Season
		}
		if ep.TVDB != nil && ep.TVDB.Episode > 0 {
			seasonEpisodeNumber = ep.TVDB.Episode
		}

		ret.Episodes[key] = &metadata.EpisodeMetadata{
			AnidbId:               0,
			TvdbId:                0,
			Title:                 strings.ReplaceAll(ep.Title, "`", "'"),
			Image:                 image,
			AirDate:               simklEpisodeAirDate(ep.Date),
			Length:                0,
			Summary:               overview,
			Overview:              overview,
			EpisodeNumber:         displayEpisodeNumber,
			Episode:               key,
			SeasonNumber:          seasonNumber,
			SeasonEpisodeNumber:   seasonEpisodeNumber,
			AbsoluteEpisodeNumber: displayEpisodeNumber,
			AnidbEid:              0,
			HasImage:              image != "",
		}
	}

	return ret
}

func buildSimklShowMetadata(mId int, episodes []simkl.Episode) *metadata.AnimeMetadata {
	// SIMKL can return seasons out of order. Keep the input untouched for other consumers.
	ordered := append([]simkl.Episode(nil), episodes...)
	season := func(ep simkl.Episode) int {
		if ep.Season == 0 && ep.TVDB != nil {
			return ep.TVDB.Season
		}
		return ep.Season
	}
	number := func(ep simkl.Episode) int {
		if ep.TVDB != nil && ep.TVDB.Episode > 0 {
			return ep.TVDB.Episode
		}
		return simklEpisodeNumber(ep)
	}
	for i := range ordered {
		if strings.EqualFold(ordered[i].Type, "special") || (ordered[i].TVDB != nil && ordered[i].TVDB.Season == 0) {
			ordered[i].Episode = number(ordered[i])
		}
		if season(ordered[i]) == 0 && (strings.EqualFold(ordered[i].Type, "special") || ordered[i].TVDB != nil) {
			ordered[i].Type = "special"
		}
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		if season(ordered[i]) != season(ordered[j]) {
			return season(ordered[i]) < season(ordered[j])
		}
		return number(ordered[i]) < number(ordered[j])
	})
	ret := buildSimklAnimeMetadata(mId, ordered)
	for key, ep := range ret.Episodes {
		if strings.HasPrefix(key, "S") {
			ep.SeasonNumber = 0
		}
	}
	ret.Mappings.Type = "tv"
	return ret
}

func simklEpisodeNumber(ep simkl.Episode) int {
	if ep.Episode != 0 {
		return ep.Episode
	}
	return ep.Number
}

func simklEpisodeAirDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed.Format("2006-01-02")
		}
	}
	return ""
}

func (p *ProviderImpl) fetchAnimapMetadata(platform metadata.Platform, mId int, ret *metadata.AnimeMetadata) (*metadata.AnimeMetadata, error) {
	if p.logger != nil {
		p.logger.Debug().Msgf("animap: Fetching metadata for %d", mId)
	}

	m, err := animap.FetchAnimapMedia(string(platform), mId)
	if err != nil || m == nil {
		return nil, err
	}

	ret.Titles = m.Titles
	ret.EpisodeCount = 0
	ret.SpecialCount = 0
	ret.Mappings.AnimeplanetId = m.Mappings.AnimePlanetID
	ret.Mappings.KitsuId = m.Mappings.KitsuID
	ret.Mappings.MalId = m.Mappings.MalID
	ret.Mappings.Type = m.Mappings.Type
	ret.Mappings.MediaId = m.Mappings.MediaID
	ret.Mappings.AnisearchId = m.Mappings.AnisearchID
	ret.Mappings.AnidbId = m.Mappings.AnidbID
	ret.Mappings.NotifymoeId = m.Mappings.NotifyMoeID
	ret.Mappings.LivechartId = m.Mappings.LivechartID
	ret.Mappings.ThetvdbId = m.Mappings.TheTvdbID
	ret.Mappings.ImdbId = ""
	ret.Mappings.ThemoviedbId = m.Mappings.TheMovieDbID

	for key, ep := range m.Episodes {
		firstChar := key[0]
		if firstChar == 'S' {
			ret.SpecialCount++
		} else {
			if firstChar >= '0' && firstChar <= '9' {
				ret.EpisodeCount++
			}
		}
		em := &metadata.EpisodeMetadata{
			AnidbId:               ep.AnidbId,
			TvdbId:                ep.TvdbId,
			Title:                 ep.AnidbTitle,
			Image:                 ep.Image,
			AirDate:               ep.AirDate,
			Length:                ep.Runtime,
			Summary:               strings.ReplaceAll(ep.Overview, "`", "'"),
			Overview:              strings.ReplaceAll(ep.Overview, "`", "'"),
			EpisodeNumber:         ep.Number,
			Episode:               key,
			SeasonNumber:          ep.SeasonNumber,
			AbsoluteEpisodeNumber: ep.AbsoluteNumber,
			AnidbEid:              ep.AnidbId,
			HasImage:              ep.Image != "",
		}
		if em.Length == 0 && ep.Runtime > 0 {
			em.Length = ep.Runtime
		}
		if em.Summary == "" && ep.Overview != "" {
			em.Summary = ep.Overview
		}
		if em.Overview == "" && ep.Overview != "" {
			em.Overview = ep.Overview
		}
		if ep.TvdbTitle != "" && ep.AnidbTitle == "Episode "+ep.AnidbEpisode {
			em.Title = ep.TvdbTitle
		}
		ret.Episodes[key] = em
	}

	return ret, nil
}
