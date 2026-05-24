package playlist

import "seall/internal/library/anime"

func isLocalFile(e *anime.PlaylistEpisode) bool {
	return e.Episode.LocalFile != nil && !e.IsNakama
}
