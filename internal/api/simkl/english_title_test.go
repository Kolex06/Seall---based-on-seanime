package simkl

import "testing"

func TestEnglishTitlesAcrossSIMKLSources(t *testing.T) {
	for _, tc := range []struct{ name, english, alias, want string }{
		{"summary", "English title", "", "English title"},
		{"search alias", "", "Search English title", "Search English title"},
		{"fallback", "", "", "Original title"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			media := StandardMedia{Title: "Original title", EnglishName: tc.english, TitleEnglish: tc.alias, IDs: IDs{Simkl: 987654321}}
			discoveryMediaCache.Delete(987654321)
			t.Cleanup(func() { discoveryMediaCache.Delete(987654321) })
			if got := ToBaseAnime(MediaTypeAnime, &media).Title; *got.UserPreferred != tc.want || *got.English != tc.want {
				t.Fatalf("expected %q, got %+v", tc.want, got)
			}
		})
	}
	items := toDiscoveryMedia(MediaTypeAnime, []restDiscoveryMedia{{Title: "Original", TitleEnglish: "English", IDs: restDiscoveryID{Simkl: 987654322}}})
	t.Cleanup(func() { discoveryMediaCache.Delete(987654322) })
	if got := ToBaseAnime(items[0].Kind, &items[0].Media); *got.Title.UserPreferred != "English" {
		t.Fatal("discovery lost the English title")
	}
	calendar := DiscoveryMediaFromCalendarItems([]CalendarItem{{Kind: MediaTypeAnime, Title: "Original", IDs: IDs{Simkl: 987654322}}})
	if got := ToBaseAnime(calendar[0].Kind, &calendar[0].Media); *got.Title.UserPreferred != "English" {
		t.Fatal("calendar replaced the cached English title")
	}
}
