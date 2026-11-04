package metadata_provider

import (
	"encoding/json"
	"seall/internal/api/simkl"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildSimklShowMetadataSeasonMappingOffline(t *testing.T) {
	var episodes []simkl.Episode
	require.NoError(t, json.Unmarshal([]byte(`[
		{"season":2,"episode":1,"title":"S2E1"},
		{"season":1,"episode":2,"title":"S1E2"},
		{"season":1,"episode":1,"title":"S1E1"},
		{"episode":99,"type":"special","tvdb":{"season":0,"episode":1},"title":"Special"},
		{"episode":4,"tvdb":{"season":2,"episode":2},"title":"S2E2"}
	]`), &episodes))
	ret := buildSimklShowMetadata(987654324, episodes)
	require.Equal(t, "tv", ret.Mappings.Type)
	require.Equal(t, 4, ret.EpisodeCount)
	require.Equal(t, "S1E1", ret.Episodes["1"].Title)
	require.Equal(t, "S2E1", ret.Episodes["3"].Title)
	require.Equal(t, 2, ret.Episodes["4"].SeasonNumber)
	require.Equal(t, 2, ret.Episodes["4"].SeasonEpisodeNumber)
	require.Equal(t, 0, ret.Episodes["S1"].SeasonNumber)
	require.Equal(t, "S2E1", episodes[0].Title)
}

func TestBuildSimklMovieMetadataUsesMovieSummary(t *testing.T) {
	ret := buildSimklMovieMetadata(987654323, &simkl.StandardMedia{Title: "Test Movie", Overview: "Movie overview", IDs: simkl.IDs{Simkl: 987654323}, Fanart: "12/123", Released: "2026-12-01"})
	require.Equal(t, 1, ret.EpisodeCount)
	require.Equal(t, "movie", ret.Mappings.Type)
	require.Equal(t, "Test Movie", ret.Episodes["1"].Title)
	require.Equal(t, "Movie overview", ret.Episodes["1"].Overview)
	require.NotEmpty(t, ret.Episodes["1"].Image)
}

func TestBuildSimklAnimeMetadataUsesEpisodeImages(t *testing.T) {
	metadata := buildSimklAnimeMetadata(2529412, []simkl.Episode{
		{
			Title:       "Those Who Shed Petals",
			Description: "A SIMKL episode description",
			Episode:     1,
			Type:        "episode",
			Image:       "19/1975492388e552455a",
			Date:        "2026-04-03T00:30:00+09:00",
		},
	})

	episode := metadata.Episodes["1"]
	require.NotNil(t, episode)
	require.Equal(t, "Those Who Shed Petals", episode.Title)
	require.Equal(t, simkl.ImageURL(simkl.ImageKindEpisode, "19/1975492388e552455a", simkl.ImageSizeEpisodeWide), episode.Image)
	require.Equal(t, "2026-04-03", episode.AirDate)
	require.True(t, episode.HasImage)
	require.Equal(t, 1, metadata.EpisodeCount)
}

func TestBuildSimklAnimeMetadataUsesAbsoluteEpisodeKeys(t *testing.T) {
	metadata := buildSimklAnimeMetadata(4390, []simkl.Episode{
		{Title: "S1E1", Episode: 1, Season: 1, Type: "episode"},
		{Title: "S1E2", Episode: 2, Season: 1, Type: "episode"},
		{Title: "S2E1", Episode: 1, Season: 2, Type: "episode"},
	})

	require.Equal(t, "S1E1", metadata.Episodes["1"].Title)
	require.Equal(t, "S1E2", metadata.Episodes["2"].Title)
	require.Equal(t, "S2E1", metadata.Episodes["3"].Title)
	require.Equal(t, 3, metadata.EpisodeCount)
}
