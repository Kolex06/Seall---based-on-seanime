package simkl

import (
	"seall/internal/api/mediaapi"
	"testing"
	"time"
)

func TestSIMKLUserRatingScalesForSeallUI(t *testing.T) {
	rating := 10
	entry := ToSIMKLAnimeEntry(MediaTypeMovies, WatchlistItem{
		Status:     WatchStatusCompleted,
		UserRating: &rating,
		Movie: &StandardMedia{
			Title: "Perfect Movie",
			IDs:   IDs{Simkl: 1},
		},
	})

	if entry == nil || entry.Score == nil {
		t.Fatal("expected entry score")
	}
	if *entry.Score != 100 {
		t.Fatalf("expected SIMKL 10/10 to become Seall 100/100, got %v", *entry.Score)
	}
}

func TestScoreToSIMKLScalesFromSeallUI(t *testing.T) {
	score := 100
	rating := ScoreToSIMKL(&score)
	if rating == nil || *rating != 10 {
		t.Fatalf("expected Seall 100/100 to become SIMKL 10/10, got %v", rating)
	}
}

func TestFutureSIMKLReleaseBecomesUpcoming(t *testing.T) {
	future := time.Now().UTC().AddDate(0, 1, 0).Format("2006-01-02")
	media := ToBaseAnime(MediaTypeMovies, &StandardMedia{
		Title:    "Future Movie",
		Released: future,
		IDs:      IDs{Simkl: 2},
	})

	if media == nil || media.Status == nil {
		t.Fatal("expected converted media status")
	}
	if *media.Status != mediaapi.MediaStatusNotYetReleased {
		t.Fatalf("expected future SIMKL release to be upcoming, got %v", *media.Status)
	}
	if media.StartDate == nil || media.StartDate.Month == nil || media.StartDate.Day == nil {
		t.Fatalf("expected full SIMKL start date, got %#v", media.StartDate)
	}
}
