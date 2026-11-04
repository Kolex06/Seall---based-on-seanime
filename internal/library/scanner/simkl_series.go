package scanner

import (
	"regexp"
	"strconv"
	"strings"

	"seall/internal/api/metadata"
	"seall/internal/library/anime"
)

var seriesEpisodeRange = regexp.MustCompile(`(?i)e\d+\s*(?:[-+]\s*e?\d+|e\d+)`)

// Resolve against SIMKL's episode table, never against the files present on disk.
func hydrateSIMKLSeriesEpisode(lf *anime.LocalFile, table *metadata.AnimeMetadata) bool {
	if lf == nil || lf.ParsedData == nil || table == nil || table.GetMappings().Type != "tv" {
		return false
	}
	if len(lf.ParsedData.EpisodeRange) != 0 || len(lf.ParsedData.SeasonRange) != 0 || seriesEpisodeRange.MatchString(lf.Name) {
		return false
	}
	number, err := strconv.Atoi(lf.ParsedData.Episode)
	if err != nil || number <= 0 {
		return false
	}
	season, explicit := 0, false
	for _, parsed := range append([]*anime.LocalFileParsedData{lf.ParsedData}, lf.ParsedFolderData...) {
		if parsed == nil || parsed.Season == "" {
			continue
		}
		season, err = strconv.Atoi(parsed.Season)
		if err != nil || season < 0 {
			return false
		}
		explicit = true
		break
	}
	var match *metadata.EpisodeMetadata
	var key string
	if explicit {
		for candidateKey, ep := range table.Episodes {
			if ep == nil || ep.SeasonNumber != season || ep.SeasonEpisodeNumber != number {
				continue
			}
			if (season == 0) != strings.HasPrefix(candidateKey, "S") {
				continue
			}
			if match != nil {
				return false
			}
			match, key = ep, candidateKey
		}
	} else {
		// Without a season marker, accept only explicitly absolute main-episode numbering.
		if lf.IsProbablySpecial() {
			return false
		}
		key = strconv.Itoa(number)
		match = table.Episodes[key]
	}
	if match == nil || match.EpisodeNumber <= 0 {
		return false
	}
	kind := anime.LocalFileTypeMain
	if strings.HasPrefix(key, "S") {
		kind = anime.LocalFileTypeSpecial
	}
	lf.Metadata = &anime.LocalFileMetadata{Episode: match.EpisodeNumber, AniDBEpisode: key, Type: kind}
	return true
}
