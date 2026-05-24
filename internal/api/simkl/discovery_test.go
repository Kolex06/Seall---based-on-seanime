package simkl

import "testing"

func TestDiscoveryUsesSimklIDFromURLWhenIDsAreMissing(t *testing.T) {
	items := toDiscoveryMedia(MediaTypeMovies, []restDiscoveryMedia{{
		Title: "A Movie",
		Year:  2026,
		URL:   "/movie/1306562/a-movie",
	}})

	if len(items) != 1 {
		t.Fatalf("expected one discovery item, got %d", len(items))
	}

	if got := items[0].Media.IDs.PrimarySimklID(); got != 1306562 {
		t.Fatalf("expected SIMKL id from URL, got %d", got)
	}
}

func TestToBaseAnimeUsesSimklIDFromURLBeforeFallback(t *testing.T) {
	media := ToBaseAnime(MediaTypeMovies, &StandardMedia{
		Title: "A Movie",
		Year:  2026,
		URL:   "https://simkl.com/movie/1306562/a-movie",
	})

	if media == nil {
		t.Fatal("expected media")
	}
	if media.ID != 1306562 {
		t.Fatalf("expected SIMKL id from URL, got %d", media.ID)
	}
}
