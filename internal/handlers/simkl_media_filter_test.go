package handlers

import (
	"seall/internal/api/mediaapi"
	simklapi "seall/internal/api/simkl"
	"testing"
	"time"
)

func TestFilterDiscoveryMediaRespectsSIMKLStatusAndFormat(t *testing.T) {
	format := mediaapi.MediaFormatMovie
	status := mediaapi.MediaStatusNotYetReleased
	future := time.Now().UTC().AddDate(0, 1, 0).Format("2006-01-02")
	past := time.Now().UTC().AddDate(-1, 0, 0).Format("2006-01-02")

	items := []simklapi.DiscoveryMedia{
		{
			Kind: simklapi.MediaTypeMovies,
			Media: simklapi.StandardMedia{
				Title:    "Future Movie",
				Released: future,
				IDs:      simklapi.IDs{Simkl: 1},
			},
		},
		{
			Kind: simklapi.MediaTypeMovies,
			Media: simklapi.StandardMedia{
				Title:    "Past Movie",
				Released: past,
				IDs:      simklapi.IDs{Simkl: 2},
			},
		},
		{
			Kind: simklapi.MediaTypeShows,
			Media: simklapi.StandardMedia{
				Title:    "Future Show",
				Released: future,
				IDs:      simklapi.IDs{Simkl: 3},
			},
		},
	}

	filtered := filterDiscoveryMedia(items, listMediaRequest{
		Format: &format,
		Status: []*mediaapi.MediaStatus{&status},
	})

	if len(filtered) != 1 {
		t.Fatalf("expected one future movie, got %d", len(filtered))
	}
	if filtered[0].Media.Title != "Future Movie" {
		t.Fatalf("expected Future Movie, got %q", filtered[0].Media.Title)
	}
}

func TestSortDiscoveryMediaBySIMKLScore(t *testing.T) {
	scoreSort := mediaapi.MediaSortScoreDesc
	items := []simklapi.DiscoveryMedia{
		{
			Kind: simklapi.MediaTypeMovies,
			Media: simklapi.StandardMedia{
				Title:   "Lower",
				IDs:     simklapi.IDs{Simkl: 1},
				Ratings: simklapi.Ratings{Simkl: &simklapi.Rating{Rating: 6.5}},
			},
		},
		{
			Kind: simklapi.MediaTypeMovies,
			Media: simklapi.StandardMedia{
				Title:   "Higher",
				IDs:     simklapi.IDs{Simkl: 2},
				Ratings: simklapi.Ratings{Simkl: &simklapi.Rating{Rating: 9.1}},
			},
		},
	}

	sortDiscoveryMedia(items, []*mediaapi.MediaSort{&scoreSort})

	if items[0].Media.Title != "Higher" {
		t.Fatalf("expected highest SIMKL score first, got %q", items[0].Media.Title)
	}
}
