package simkl_platform

import (
	"seall/internal/api/mediaapi"
	"testing"
)

func TestSimklViewerStatsFromCollection(t *testing.T) {
	movieFormat := mediaapi.MediaFormatMovie
	tvFormat := mediaapi.MediaFormatTv
	completed := mediaapi.MediaListStatusCompleted
	current := mediaapi.MediaListStatusCurrent
	scoreOne := 100.0
	scoreTwo := 80.0
	progressOne := 1
	progressTwo := 10
	durationOne := 120
	durationTwo := 24
	yearOne := 2001
	yearTwo := 2024
	genre := "Action"

	stats := simklViewerStatsFromCollection(&mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{
				{
					Entries: []*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{
						{
							ID:       1,
							Progress: &progressOne,
							Score:    &scoreOne,
							Status:   &completed,
							Media: &mediaapi.BaseAnime{
								ID:         1,
								Duration:   &durationOne,
								Format:     &movieFormat,
								SeasonYear: &yearOne,
								Genres:     []*string{&genre},
							},
						},
						{
							ID:       2,
							Progress: &progressTwo,
							Score:    &scoreTwo,
							Status:   &current,
							Media: &mediaapi.BaseAnime{
								ID:         2,
								Duration:   &durationTwo,
								Format:     &tvFormat,
								SeasonYear: &yearTwo,
							},
						},
					},
				},
			},
		},
	})

	anime := stats.GetViewer().GetStatistics().GetAnime()
	if anime.Count != 2 {
		t.Fatalf("expected count 2, got %d", anime.Count)
	}
	if anime.MinutesWatched != 360 {
		t.Fatalf("expected 360 watched minutes, got %d", anime.MinutesWatched)
	}
	if anime.EpisodesWatched != 11 {
		t.Fatalf("expected 11 watched episodes, got %d", anime.EpisodesWatched)
	}
	if anime.MeanScore != 90 {
		t.Fatalf("expected mean score 90, got %f", anime.MeanScore)
	}
	if len(anime.Formats) != 2 || len(anime.Statuses) != 2 || len(anime.ReleaseYears) != 2 || len(anime.Genres) != 1 {
		t.Fatalf("expected populated format/status/year/genre stats, got formats=%d statuses=%d years=%d genres=%d", len(anime.Formats), len(anime.Statuses), len(anime.ReleaseYears), len(anime.Genres))
	}
}
