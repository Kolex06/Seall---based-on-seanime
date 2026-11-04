package scanner

import (
	"testing"

	"github.com/stretchr/testify/require"
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata"
	"seall/internal/library/anime"
)

func TestSIMKLSeriesLocalMappingOffline(t *testing.T) {
	table := &metadata.AnimeMetadata{Mappings: &metadata.AnimeMappings{Type: "tv"}, Episodes: map[string]*metadata.EpisodeMetadata{
		"1":  {EpisodeNumber: 1, SeasonNumber: 1, SeasonEpisodeNumber: 1},
		"2":  {EpisodeNumber: 2, SeasonNumber: 1, SeasonEpisodeNumber: 2},
		"3":  {EpisodeNumber: 3, SeasonNumber: 2, SeasonEpisodeNumber: 1},
		"4":  {EpisodeNumber: 4, SeasonNumber: 2, SeasonEpisodeNumber: 2},
		"5":  {EpisodeNumber: 5, SeasonNumber: 3, SeasonEpisodeNumber: 1},
		"S1": {EpisodeNumber: 1, SeasonNumber: 0, SeasonEpisodeNumber: 1},
	}}
	for _, tc := range []struct {
		name, path, key string
		ok              bool
	}{
		{"season one", "Show/Show.S01E02.mkv", "2", true},
		{"season two", "Show/Show.S02E01.mkv", "3", true},
		{"season three", "Show/Show.S03E01.mkv", "5", true},
		{"folder season", "Show/Season 2/Show - 02.mkv", "4", true},
		{"special", "Show/Show.S00E01.mkv", "S1", true},
		{"absolute", "Show/Show - 04.mkv", "4", true},
		{"unknown season", "Show/Show.S04E01.mkv", "", false},
		{"unknown episode", "Show/Show.S02E09.mkv", "", false},
		{"missing number", "Show/Show.mkv", "", false},
		{"multi episode", "Show/Show.S02E01-E02.mkv", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lf := anime.NewLocalFile("C:/Library/"+tc.path, "C:/Library")
			require.Equal(t, tc.ok, hydrateSIMKLSeriesEpisode(lf, table))
			if tc.ok {
				require.Equal(t, tc.key, lf.Metadata.AniDBEpisode)
			}
		})
	}
	t.Run("ambiguous mapping", func(t *testing.T) {
		table.Episodes["6"] = &metadata.EpisodeMetadata{EpisodeNumber: 6, SeasonNumber: 2, SeasonEpisodeNumber: 1}
		lf := anime.NewLocalFile("C:/Library/Show.S02E01.mkv", "C:/Library")
		require.False(t, hydrateSIMKLSeriesEpisode(lf, table))
	})
}

func TestSIMKLSeriesClassificationOffline(t *testing.T) {
	for _, tc := range []struct {
		url  string
		want bool
	}{
		{"https://simkl.com/tv/123/show", true},
		{"https://simkl.com/anime/123/show", false},
		{"https://simkl.com/movie/123/movie", false},
		{"https://example.com/tv/123", false},
		{"", false},
	} {
		m := anime.NewNormalizedMedia(&mediaapi.BaseAnime{SiteURL: &tc.url})
		require.Equal(t, tc.want, m.IsSeries, tc.url)
	}
}
