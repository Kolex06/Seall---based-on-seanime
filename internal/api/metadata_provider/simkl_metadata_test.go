package metadata_provider

import (
	"seall/internal/api/simkl"
	"testing"

	"github.com/stretchr/testify/require"
)

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
