package onlinestream

import (
	"seall/internal/api/mediaapi"
	"testing"
)

func TestMediaTypeForBaseMediaUsesSimklURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{name: "movie", url: "https://simkl.com/movie/1718303712/example", want: "movies"},
		{name: "show", url: "https://simkl.com/tv/123/example", want: "shows"},
		{name: "anime", url: "https://simkl.com/anime/456/example", want: "anime"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := &mediaapi.BaseAnime{SiteURL: &tt.url}
			if got := MediaTypeForBaseMedia(media); got != tt.want {
				t.Fatalf("MediaTypeForBaseMedia() = %q, want %q", got, tt.want)
			}
		})
	}
}
