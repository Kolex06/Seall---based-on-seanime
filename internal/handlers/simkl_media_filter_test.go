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

func TestFilterDiscoveryMediaRequiresKnownSeasonMonth(t *testing.T) {
	season := mediaapi.MediaSeasonSpring
	year := 2026
	items := []simklapi.DiscoveryMedia{
		{
			Kind: simklapi.MediaTypeMovies,
			Media: simklapi.StandardMedia{
				Title:    "Spring Movie",
				Year:     year,
				Released: "2026-04-10",
				IDs:      simklapi.IDs{Simkl: 1},
			},
		},
		{
			Kind: simklapi.MediaTypeMovies,
			Media: simklapi.StandardMedia{
				Title: "Year Only Movie",
				Year:  year,
				IDs:   simklapi.IDs{Simkl: 2},
			},
		},
	}

	filtered := filterDiscoveryMedia(items, listMediaRequest{
		Season:     &season,
		SeasonYear: &year,
	})

	if len(filtered) != 1 {
		t.Fatalf("expected only the dated spring item, got %d", len(filtered))
	}
	if filtered[0].Media.Title != "Spring Movie" {
		t.Fatalf("expected Spring Movie, got %q", filtered[0].Media.Title)
	}
}

func TestMatchesGenresNormalizesSIMKLAliases(t *testing.T) {
	wantedSciFi := "Science Fiction"
	wantedSports := "Sports"
	wantedKids := "Kids"
	wantedSuperPower := "Super Power"

	if !matchesGenres([]string{"Sci-Fi", "Sport", "Children", "Superpower"}, []*string{&wantedSciFi, &wantedSports, &wantedKids, &wantedSuperPower}) {
		t.Fatalf("expected SIMKL genre aliases to match discover genre filters")
	}
}

func TestInterleaveDiscoveryMediaGroups(t *testing.T) {
	groups := [][]simklapi.DiscoveryMedia{
		{
			{Kind: simklapi.MediaTypeMovies, Media: simklapi.StandardMedia{Title: "Movie 1"}},
			{Kind: simklapi.MediaTypeMovies, Media: simklapi.StandardMedia{Title: "Movie 2"}},
		},
		{
			{Kind: simklapi.MediaTypeShows, Media: simklapi.StandardMedia{Title: "Show 1"}},
		},
		{
			{Kind: simklapi.MediaTypeAnime, Media: simklapi.StandardMedia{Title: "Anime 1"}},
		},
	}

	ret := interleaveDiscoveryMediaGroups(groups)
	titles := []string{ret[0].Media.Title, ret[1].Media.Title, ret[2].Media.Title, ret[3].Media.Title}
	expected := []string{"Movie 1", "Show 1", "Anime 1", "Movie 2"}
	for i := range expected {
		if titles[i] != expected[i] {
			t.Fatalf("expected %v, got %v", expected, titles)
		}
	}
}

func TestBalanceDiscoveryMediaKindsKeepsAnimeVisible(t *testing.T) {
	items := []simklapi.DiscoveryMedia{
		{Kind: simklapi.MediaTypeMovies, Media: simklapi.StandardMedia{Title: "Movie 1"}},
		{Kind: simklapi.MediaTypeMovies, Media: simklapi.StandardMedia{Title: "Movie 2"}},
		{Kind: simklapi.MediaTypeMovies, Media: simklapi.StandardMedia{Title: "Movie 3"}},
		{Kind: simklapi.MediaTypeShows, Media: simklapi.StandardMedia{Title: "Show 1"}},
		{Kind: simklapi.MediaTypeShows, Media: simklapi.StandardMedia{Title: "Show 2"}},
		{Kind: simklapi.MediaTypeAnime, Media: simklapi.StandardMedia{Title: "Anime 1", Type: string(simklapi.MediaTypeAnime)}},
		{Kind: simklapi.MediaTypeAnime, Media: simklapi.StandardMedia{Title: "Anime 2", Type: string(simklapi.MediaTypeAnime)}},
	}

	ret := balanceDiscoveryMediaKinds(items)
	if ret[2].Media.Title != "Anime 1" {
		t.Fatalf("expected anime to be pulled into the first mixed group, got %q", ret[2].Media.Title)
	}
	if ret[5].Media.Title != "Anime 2" {
		t.Fatalf("expected anime to stay visible in the second mixed group, got %q", ret[5].Media.Title)
	}
}

func TestListRecentAnimeFromCalendarItemsUsesSIMKLCalendarDates(t *testing.T) {
	sortTime := mediaapi.AiringSortTime
	greater := int(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second).Unix())
	lesser := int(time.Date(2026, 6, 3, 0, 0, 0, 0, time.UTC).Unix())

	ret := listRecentAnimeFromCalendarItems([]simklapi.CalendarItem{
		{
			Kind:        simklapi.MediaTypeMovies,
			Title:       "Calendar Movie",
			ReleaseDate: "2026-06-01",
			IDs:         simklapi.IDs{Simkl: 101},
		},
		{
			Kind:  simklapi.MediaTypeShows,
			Title: "Calendar Show",
			Date:  "2026-06-02T12:30:00Z",
			IDs:   simklapi.IDs{Simkl: 202},
			Episode: &simklapi.CalendarEpisode{
				Episode: 4,
			},
		},
	}, 1, 10, "", &greater, &lesser, []*mediaapi.AiringSort{&sortTime})

	schedules := ret.GetPage().GetAiringSchedules()
	if len(schedules) != 2 {
		t.Fatalf("expected two calendar schedules, got %d", len(schedules))
	}
	if schedules[0].GetMedia().GetID() != 101 || schedules[0].GetAiringAt() != greater+1 {
		t.Fatalf("expected movie release first at the calendar date, got id=%d airingAt=%d", schedules[0].GetMedia().GetID(), schedules[0].GetAiringAt())
	}
	if schedules[1].GetMedia().GetID() != 202 || schedules[1].GetEpisode() != 4 {
		t.Fatalf("expected show episode 4 second, got id=%d episode=%d", schedules[1].GetMedia().GetID(), schedules[1].GetEpisode())
	}
	if schedules[1].GetMedia().GetNextAiringEpisode().GetAiringAt() != schedules[1].GetAiringAt() {
		t.Fatalf("expected media next airing data to match schedule airing time")
	}
}
