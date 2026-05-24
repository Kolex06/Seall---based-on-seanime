package metadata_provider

import (
	"context"
	"errors"
	"fmt"
	"seall/internal/api/animap"
	"seall/internal/api/metadata"
	"seall/internal/api/simkl"
	"strconv"
	"strings"
	"time"
)

func (p *ProviderImpl) fetchSimklEpisodeMetadata(ctx context.Context, mId int) (*metadata.AnimeMetadata, error) {
	if p.simklClientRef == nil || p.simklClientRef.IsAbsent() {
		return nil, errors.New("simkl client is unavailable")
	}

	client := p.simklClientRef.Get()
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
		return buildSimklAnimeMetadata(mId, episodes), nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("simkl: no episode metadata returned for %d", mId)
	}
	return nil, lastErr
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
		if seasonNumber == 0 && ep.TVDB != nil {
			seasonNumber = ep.TVDB.Season
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
			AbsoluteEpisodeNumber: displayEpisodeNumber,
			AnidbEid:              0,
			HasImage:              image != "",
		}
	}

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
