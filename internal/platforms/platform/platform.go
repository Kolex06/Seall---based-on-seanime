package platform

import (
	"context"
	"seall/internal/api/mediaapi"
)

type Platform interface {
	SetUsername(username string)
	// UpdateEntry updates the entry for the given media ID
	UpdateEntry(context context.Context, mediaID int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput) error
	// UpdateEntryProgress updates the entry progress for the given media ID
	UpdateEntryProgress(context context.Context, mediaID int, progress int, totalEpisodes *int) error
	// UpdateEntryRepeat updates the entry repeat number for the given media ID
	UpdateEntryRepeat(context context.Context, mediaID int, repeat int) error
	// DeleteEntry deletes the entry for the given media ID
	DeleteEntry(context context.Context, mediaID int, entryID int) error
	// GetAnime gets the anime for the given media ID
	GetAnime(context context.Context, mediaID int) (*mediaapi.BaseAnime, error)
	// GetAnimeByMalID gets the anime by MAL ID
	GetAnimeByMalID(context context.Context, malID int) (*mediaapi.BaseAnime, error)
	// GetAnimeWithRelations gets the anime with relations for the given media ID
	// This is used for scanning purposes in order to build the relation tree
	GetAnimeWithRelations(context context.Context, mediaID int) (*mediaapi.CompleteAnime, error)
	// GetAnimeDetails gets the anime details for the given media ID
	// These details are only fetched by the anime page
	GetAnimeDetails(context context.Context, mediaID int) (*mediaapi.AnimeDetailsById_Media, error)
	// GetManga gets the manga for the given media ID
	GetManga(context context.Context, mediaID int) (*mediaapi.BaseManga, error)
	// GetMediaCollection gets the anime collection without custom lists
	// This should not make any API calls and instead should be based on GetRawMediaCollection
	GetMediaCollection(context context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error)
	// GetRawMediaCollection gets the anime collection with custom lists
	GetRawMediaCollection(context context.Context, bypassCache bool) (*mediaapi.AnimeCollection, error)
	// GetMangaDetails gets the manga details for the given media ID
	// These details are only fetched by the manga page
	GetMangaDetails(context context.Context, mediaID int) (*mediaapi.MangaDetailsById_Media, error)
	// GetMediaCollectionWithRelations gets the anime collection with relations
	// This is used for scanning purposes in order to build the relation tree
	GetMediaCollectionWithRelations(context context.Context) (*mediaapi.AnimeCollectionWithRelations, error)
	// GetMangaCollection gets the manga collection without custom lists
	// This should not make any API calls and instead should be based on GetRawMangaCollection
	GetMangaCollection(context context.Context, bypassCache bool) (*mediaapi.MangaCollection, error)
	// GetRawMangaCollection gets the manga collection with custom lists
	GetRawMangaCollection(context context.Context, bypassCache bool) (*mediaapi.MangaCollection, error)
	// AddMediaToCollection adds the media to the collection
	AddMediaToCollection(context context.Context, mIds []int) error
	// GetStudioDetails gets the studio details for the given studio ID
	GetStudioDetails(context context.Context, studioID int) (*mediaapi.StudioDetails, error)
	// GetMediaApiClient gets the SIMKL client
	GetMediaApiClient() mediaapi.MediaApiClient
	// RefreshAnimeCollection refreshes the anime collection
	RefreshAnimeCollection(context context.Context) (*mediaapi.AnimeCollection, error)
	// RefreshMangaCollection refreshes the manga collection
	RefreshMangaCollection(context context.Context) (*mediaapi.MangaCollection, error)
	// GetViewerStats gets the viewer stats
	GetViewerStats(context context.Context) (*mediaapi.ViewerStats, error)
	// GetAnimeAiringSchedule gets the schedule for airing anime in the collection
	GetAnimeAiringSchedule(context context.Context) (*mediaapi.AnimeAiringSchedule, error)
	ClearCache()
	Close()
}
