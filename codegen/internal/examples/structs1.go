package codegen

import (
	"seall/internal/api/mediaapi"
	hibiketorrent "seall/internal/extension/hibike/torrent"
)

//type Struct1 struct {
//	Struct2
//}
//
//type Struct2 struct {
//	Text string `json:"text"`
//}

//type Struct3 []string

type Struct4 struct {
	Torrents    []hibiketorrent.MediaTorrent `json:"torrents"`
	Destination string                       `json:"destination"`
	SmartSelect struct {
		Enabled               bool  `json:"enabled"`
		MissingEpisodeNumbers []int `json:"missingEpisodeNumbers"`
	} `json:"smartSelect"`
	Media *mediaapi.BaseAnime `json:"media"`
}
