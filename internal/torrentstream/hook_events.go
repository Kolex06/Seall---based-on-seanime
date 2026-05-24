package torrentstream

import (
	"seall/internal/api/mediaapi"
	hibiketorrent "seall/internal/extension/hibike/torrent"
	"seall/internal/hook_resolver"
)

// TorrentStreamAutoSelectTorrentsFetchedEvent is triggered when the torrents are fetched for auto select.
// The torrents are sorted by seeders from highest to lowest.
// This event is triggered before the top 3 torrents are analyzed.
type TorrentStreamAutoSelectTorrentsFetchedEvent struct {
	hook_resolver.Event
	Torrents []*hibiketorrent.MediaTorrent
}

// TorrentStreamSendStreamToMediaPlayerEvent is triggered when the torrent stream is about to send a stream to the media player.
// Prevent default to skip the default playback and override the playback.
type TorrentStreamSendStreamToMediaPlayerEvent struct {
	hook_resolver.Event
	WindowTitle  string              `json:"windowTitle"`
	StreamURL    string              `json:"streamURL"`
	Media        *mediaapi.BaseAnime `json:"media"`
	AniDbEpisode string              `json:"aniDbEpisode"`
	PlaybackType string              `json:"playbackType"`
}
