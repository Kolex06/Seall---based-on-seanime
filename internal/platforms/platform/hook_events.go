package platform

import (
	"seall/internal/api/mediaapi"
	"seall/internal/hook_resolver"
)

/////////////////////////////
// SIMKL Events
/////////////////////////////

type GetAnimeEvent struct {
	hook_resolver.Event
	Anime *mediaapi.BaseAnime `json:"anime"`
}

type GetAnimeDetailsEvent struct {
	hook_resolver.Event
	Anime *mediaapi.AnimeDetailsById_Media `json:"anime"`
}

type GetMangaEvent struct {
	hook_resolver.Event
	Manga *mediaapi.BaseManga `json:"manga"`
}

type GetMangaDetailsEvent struct {
	hook_resolver.Event
	Manga *mediaapi.MangaDetailsById_Media `json:"manga"`
}

type GetCachedAnimeCollectionEvent struct {
	hook_resolver.Event
	AnimeCollection *mediaapi.AnimeCollection `json:"animeCollection"`
}

type GetCachedMangaCollectionEvent struct {
	hook_resolver.Event
	MangaCollection *mediaapi.MangaCollection `json:"mangaCollection"`
}

type GetMediaCollectionEvent struct {
	hook_resolver.Event
	AnimeCollection *mediaapi.AnimeCollection `json:"animeCollection"`
}

type GetMangaCollectionEvent struct {
	hook_resolver.Event
	MangaCollection *mediaapi.MangaCollection `json:"mangaCollection"`
}

type GetCachedRawAnimeCollectionEvent struct {
	hook_resolver.Event
	AnimeCollection *mediaapi.AnimeCollection `json:"animeCollection"`
}

type GetCachedRawMangaCollectionEvent struct {
	hook_resolver.Event
	MangaCollection *mediaapi.MangaCollection `json:"mangaCollection"`
}

type GetRawMediaCollectionEvent struct {
	hook_resolver.Event
	AnimeCollection *mediaapi.AnimeCollection `json:"animeCollection"`
}

type GetRawMangaCollectionEvent struct {
	hook_resolver.Event
	MangaCollection *mediaapi.MangaCollection `json:"mangaCollection"`
}

type GetStudioDetailsEvent struct {
	hook_resolver.Event
	Studio *mediaapi.StudioDetails `json:"studio"`
}

// PreUpdateEntryEvent is triggered when an entry is about to be updated.
// Prevent default to skip the default update and override the update.
type PreUpdateEntryEvent struct {
	hook_resolver.Event
	MediaID     *int                      `json:"mediaId"`
	Status      *mediaapi.MediaListStatus `json:"status"`
	ScoreRaw    *int                      `json:"scoreRaw"`
	Progress    *int                      `json:"progress"`
	StartedAt   *mediaapi.FuzzyDateInput  `json:"startedAt"`
	CompletedAt *mediaapi.FuzzyDateInput  `json:"completedAt"`
}

type PostUpdateEntryEvent struct {
	hook_resolver.Event
	MediaID *int `json:"mediaId"`
}

// PreUpdateEntryProgressEvent is triggered when an entry's progress is about to be updated.
// Prevent default to skip the default update and override the update.
type PreUpdateEntryProgressEvent struct {
	hook_resolver.Event
	MediaID    *int `json:"mediaId"`
	Progress   *int `json:"progress"`
	TotalCount *int `json:"totalCount"`
	// Defaults to mediaapi.MediaListStatusCurrent
	Status *mediaapi.MediaListStatus `json:"status"`
}

type PostUpdateEntryProgressEvent struct {
	hook_resolver.Event
	MediaID *int `json:"mediaId"`
}

// PreUpdateEntryRepeatEvent is triggered when an entry's repeat is about to be updated.
// Prevent default to skip the default update and override the update.
type PreUpdateEntryRepeatEvent struct {
	hook_resolver.Event
	MediaID *int `json:"mediaId"`
	Repeat  *int `json:"repeat"`
}

type PostUpdateEntryRepeatEvent struct {
	hook_resolver.Event
	MediaID *int `json:"mediaId"`
}

// PreDeleteEntryEvent is triggered when an entry is about to be deleted.
// Prevent default to skip the default deletion and override the deletion.
type PreDeleteEntryEvent struct {
	hook_resolver.Event
	MediaID *int `json:"mediaId"`
	EntryID *int `json:"entryId"`
}

type PostDeleteEntryEvent struct {
	hook_resolver.Event
	MediaID *int `json:"mediaId"`
	EntryID *int `json:"entryId"`
}
