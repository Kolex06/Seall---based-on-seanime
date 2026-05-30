package onlinestream

import (
	"seall/internal/api/metadata"
	hibikeonlinestream "seall/internal/extension/hibike/onlinestream"
	"seall/internal/util/comparison"
	"sort"
)

func providerEpisodeSeasonNumbers(animeMetadata *metadata.AnimeMetadata, episodeDetails *hibikeonlinestream.EpisodeDetails) []int {
	if episodeDetails == nil {
		return nil
	}

	if season := comparison.ExtractSeasonNumber(episodeDetails.Title); season > 0 {
		return []int{season}
	}

	if animeMetadata == nil || animeMetadata.Episodes == nil {
		return nil
	}

	seen := make(map[int]struct{})
	for _, episode := range animeMetadata.Episodes {
		if episode == nil || episode.SeasonNumber <= 0 {
			continue
		}
		if episode.EpisodeNumber == episodeDetails.Number || episode.AbsoluteEpisodeNumber == episodeDetails.Number {
			seen[episode.SeasonNumber] = struct{}{}
		}
	}

	ret := make([]int, 0, len(seen))
	for season := range seen {
		ret = append(ret, season)
	}
	sort.Ints(ret)
	return ret
}

func providerEpisodeSeasonNumber(animeMetadata *metadata.AnimeMetadata, episodeDetails *hibikeonlinestream.EpisodeDetails) int {
	seasons := providerEpisodeSeasonNumbers(animeMetadata, episodeDetails)
	if len(seasons) == 1 {
		return seasons[0]
	}
	return 0
}

func providerEpisodeMatchesSeason(animeMetadata *metadata.AnimeMetadata, episodeDetails *hibikeonlinestream.EpisodeDetails, seasonNumber int) bool {
	if seasonNumber <= 0 {
		return true
	}
	for _, season := range providerEpisodeSeasonNumbers(animeMetadata, episodeDetails) {
		if season == seasonNumber {
			return true
		}
	}
	return false
}
