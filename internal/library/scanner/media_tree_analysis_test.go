package scanner

import (
	"context"
	"seall/internal/api/mediaapi"
	"seall/internal/util"
	"seall/internal/util/limiter"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestMediaTreeAnalysis(t *testing.T) {
	wrapper := newScannerLiveWrapper(t)
	simklClient := wrapper.MediaApiClient
	simklRateLimiter := wrapper.MediaApiRateLimiter

	tests := []struct {
		name                          string
		mediaId                       int
		absoluteEpisodeNumber         int
		expectedRelativeEpisodeNumber int
	}{
		{
			name:                          "Media Tree Analysis for 86 - Eighty Six Part 2",
			mediaId:                       131586, // 86 - Eighty Six Part 2
			absoluteEpisodeNumber:         23,
			expectedRelativeEpisodeNumber: 12,
		},
		{
			name:                          "Oshi no Ko Season 2",
			mediaId:                       150672, // 86 - Eighty Six Part 2
			absoluteEpisodeNumber:         12,
			expectedRelativeEpisodeNumber: 1,
		},
		{
			name:                          "Re:zero",
			mediaId:                       21355, // Re:Zero kara Hajimeru Isekai Seikatsu
			absoluteEpisodeNumber:         51,
			expectedRelativeEpisodeNumber: 1,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mediaF, err := simklClient.BaseAnimeByID(context.Background(), &tt.mediaId)
			if err != nil {
				t.Fatal("expected media, got not found")
			}
			media := mediaF.GetMedia()
			tree := mediaapi.NewCompleteAnimeRelationTree()

			// +---------------------+
			// |     MediaTree       |
			// +---------------------+

			err = media.FetchMediaTree(
				mediaapi.FetchMediaTreeAll,
				simklClient,
				simklRateLimiter,
				tree,
				mediaapi.NewCompleteAnimeCache(),
			)

			if err != nil {
				t.Fatal("expected media tree, got error:", err.Error())
			}

			// +---------------------+
			// |  MediaTreeAnalysis  |
			// +---------------------+

			mta, err := NewMediaTreeAnalysis(&MediaTreeAnalysisOptions{
				tree:                tree,
				metadataProviderRef: util.NewRef(wrapper.MetadataProvider),
				rateLimiter:         limiter.NewLimiter(time.Minute, 25),
			})
			if err != nil {
				t.Fatal("expected media tree analysis, got error:", err.Error())
			}

			// +---------------------+
			// |  Relative Episode   |
			// +---------------------+

			relEp, _, ok := mta.getRelativeEpisodeNumber(tt.absoluteEpisodeNumber)

			if assert.Truef(t, ok, "expected relative episode number %v for absolute episode number %v, nothing found", tt.expectedRelativeEpisodeNumber, tt.absoluteEpisodeNumber) {

				assert.Equal(t, tt.expectedRelativeEpisodeNumber, relEp)

			}

		})

	}

}

func TestMediaTreeAnalysis2(t *testing.T) {
	wrapper := newScannerLiveWrapper(t)
	simklClient := wrapper.MediaApiClient
	simklRateLimiter := wrapper.MediaApiRateLimiter

	tests := []struct {
		name    string
		mediaId int
	}{
		{
			name:    "Media Tree Analysis",
			mediaId: 375, // Soreyuke! Uchuu Senkan Yamamoto Yohko
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			media, err := simklClient.BaseAnimeByID(context.Background(), &tt.mediaId)
			if err != nil {
				t.Fatal("expected media, got error:", err.Error())
			}
			tree := mediaapi.NewCompleteAnimeRelationTree()

			// +---------------------+
			// |     MediaTree       |
			// +---------------------+

			err = media.GetMedia().FetchMediaTree(
				mediaapi.FetchMediaTreeAll,
				simklClient,
				simklRateLimiter,
				tree,
				mediaapi.NewCompleteAnimeCache(),
			)

			if err != nil {
				t.Fatal("expected media tree, got error:", err.Error())
			}

			// +---------------------+
			// |  MediaTreeAnalysis  |
			// +---------------------+

			mta, err := NewMediaTreeAnalysis(&MediaTreeAnalysisOptions{
				tree:                tree,
				metadataProviderRef: util.NewRef(wrapper.MetadataProvider),
				rateLimiter:         limiter.NewLimiter(time.Minute, 25),
			})
			if err != nil {
				t.Fatal("expected media tree analysis, got error:", err.Error())
			}

			t.Log(spew.Sdump(mta))

		})

	}

}
