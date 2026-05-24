package simkl

import (
	"seall/internal/api/mediaapi"
	"testing"
)

func TestDiscoveryMediaFromCalendarItemsPreservesUpcomingMovieDate(t *testing.T) {
	items := []CalendarItem{
		{
			Kind:        MediaTypeMovies,
			Title:       "Future Movie",
			ReleaseDate: "2099-06-15",
			URL:         "/movie/123/future-movie",
			IDs:         IDs{Simkl: 123},
		},
	}

	discoveryItems := DiscoveryMediaFromCalendarItems(items)
	if len(discoveryItems) != 1 {
		t.Fatalf("expected one discovery item, got %d", len(discoveryItems))
	}

	media := ToBaseAnime(discoveryItems[0].Kind, &discoveryItems[0].Media)
	if media.GetFormat() == nil || *media.GetFormat() != mediaapi.MediaFormatMovie {
		t.Fatalf("expected movie format, got %v", media.GetFormat())
	}
	if media.GetStatus() == nil || *media.GetStatus() != mediaapi.MediaStatusNotYetReleased {
		t.Fatalf("expected upcoming status, got %v", media.GetStatus())
	}
	if media.GetStartDate().GetYear() == nil || *media.GetStartDate().GetYear() != 2099 {
		t.Fatalf("expected year 2099, got %v", media.GetStartDate().GetYear())
	}
	if media.GetStartDate().GetMonth() == nil || *media.GetStartDate().GetMonth() != 6 {
		t.Fatalf("expected month 6, got %v", media.GetStartDate().GetMonth())
	}
	if media.GetStartDate().GetDay() == nil || *media.GetStartDate().GetDay() != 15 {
		t.Fatalf("expected day 15, got %v", media.GetStartDate().GetDay())
	}
}
