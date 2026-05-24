package anime

import (
	"context"
	"errors"
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata"
	"seall/internal/api/metadata_provider"
	"seall/internal/hook"
	"seall/internal/platforms/platform"
	"seall/internal/util"
	"sort"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

type (
	// Entry is a container for all data related to a media.
	// It is the primary data structure used by the frontend.
	Entry struct {
		MediaId             int                 `json:"mediaId"`
		Media               *mediaapi.BaseAnime `json:"media"`
		EntryListData       *EntryListData      `json:"listData"`
		EntryLibraryData    *EntryLibraryData   `json:"libraryData"`
		EntryDownloadInfo   *EntryDownloadInfo  `json:"downloadInfo,omitempty"`
		Episodes            []*Episode          `json:"episodes"`
		NextEpisode         *Episode            `json:"nextEpisode"`
		LocalFiles          []*LocalFile        `json:"localFiles"`
		AnidbId             int                 `json:"anidbId"`
		CurrentEpisodeCount int                 `json:"currentEpisodeCount"`

		IsNakamaEntry     bool                    `json:"_isNakamaEntry"`
		NakamaLibraryData *NakamaEntryLibraryData `json:"nakamaLibraryData,omitempty"`
	}

	// EntryListData holds the details of the SIMKL entry.
	EntryListData struct {
		Progress    int                       `json:"progress,omitempty"`
		Score       float64                   `json:"score,omitempty"`
		Status      *mediaapi.MediaListStatus `json:"status,omitempty"`
		Repeat      int                       `json:"repeat,omitempty"`
		StartedAt   string                    `json:"startedAt,omitempty"`
		CompletedAt string                    `json:"completedAt,omitempty"`
	}
)

type (
	// NewEntryOptions is a constructor for Entry.
	NewEntryOptions struct {
		MediaId             int
		LocalFiles          []*LocalFile // All local files
		AnimeCollection     *mediaapi.AnimeCollection
		PlatformRef         *util.Ref[platform.Platform]
		MetadataProviderRef *util.Ref[metadata_provider.Provider]
		IsSimulated         bool // If the account is simulated
	}
)

// NewEntry creates a new Entry based on the media id and a list of local files.
// A Entry is a container for all data related to a media.
// It is the primary data structure used by the frontend.
//
// It has the following properties:
//   - EntryListData: Details of the SIMKL entry (if any)
//   - EntryLibraryData: Details of the local files (if any)
//   - EntryDownloadInfo: Details of the download status
//   - Episodes: List of episodes (if any)
//   - NextEpisode: Next episode to watch (if any)
//   - LocalFiles: List of local files (if any)
//   - AnidbId: AniDB id
//   - CurrentEpisodeCount: Current episode count
func NewEntry(ctx context.Context, opts *NewEntryOptions) (*Entry, error) {
	// Create new Entry
	entry := new(Entry)
	entry.MediaId = opts.MediaId

	reqEvent := new(AnimeEntryRequestedEvent)
	reqEvent.MediaId = opts.MediaId
	reqEvent.LocalFiles = opts.LocalFiles
	reqEvent.AnimeCollection = opts.AnimeCollection
	reqEvent.Entry = entry

	err := hook.GlobalHookManager.OnAnimeEntryRequested().Trigger(reqEvent)
	if err != nil {
		return nil, err
	}
	opts.MediaId = reqEvent.MediaId                 // Override the media ID
	opts.LocalFiles = reqEvent.LocalFiles           // Override the local files
	opts.AnimeCollection = reqEvent.AnimeCollection // Override the anime collection
	entry = reqEvent.Entry                          // Override the entry

	// Default prevented, return the modified entry
	if reqEvent.DefaultPrevented {
		event := new(AnimeEntryEvent)
		event.Entry = reqEvent.Entry
		err = hook.GlobalHookManager.OnAnimeEntry().Trigger(event)
		if err != nil {
			return nil, err
		}

		if event.Entry == nil {
			return nil, errors.New("no entry was returned")
		}
		return event.Entry, nil
	}

	if opts.AnimeCollection == nil ||
		opts.PlatformRef.IsAbsent() {
		return nil, errors.New("missing arguments when creating media entry")
	}

	// +---------------------+
	// |   SIMKL entry     |
	// +---------------------+

	// Get the Simkl List entry
	simklEntry, found := opts.AnimeCollection.GetListEntryFromAnimeId(opts.MediaId)

	// Set the media
	// If the Simkl List entry does not exist, fetch the media from SIMKL
	if !found {
		// If the Simkl entry does not exist, instantiate one with zero values
		simklEntry = &mediaapi.AnimeListEntry{}

		// Fetch the media
		fetchedMedia, err := opts.PlatformRef.Get().GetAnime(ctx, opts.MediaId) // DEVNOTE: Maybe cache it?
		if err != nil {
			return nil, err
		}
		entry.Media = fetchedMedia
	} else {
		animeEvent := new(platform.GetAnimeEvent)
		animeEvent.Anime = simklEntry.Media
		err := hook.GlobalHookManager.OnGetAnime().Trigger(animeEvent)
		if err != nil {
			return nil, err
		}
		entry.Media = animeEvent.Anime
	}

	// If the account is simulated and the media was in the library, we will still fetch
	// the media from SIMKL to ensure we have the latest data
	if opts.IsSimulated && found {
		// Fetch the media
		fetchedMedia, err := opts.PlatformRef.Get().GetAnime(ctx, opts.MediaId) // DEVNOTE: Maybe cache it?
		if err != nil {
			return nil, err
		}
		entry.Media = fetchedMedia
	}

	entry.CurrentEpisodeCount = entry.Media.GetCurrentEpisodeCount()

	// +---------------------+
	// |     Local files     |
	// +---------------------+

	// Get the entry's local files
	lfs := GetLocalFilesFromMediaId(opts.LocalFiles, opts.MediaId)
	entry.LocalFiles = lfs // Returns empty slice if no local files are found

	libraryData, _ := NewEntryLibraryData(&NewEntryLibraryDataOptions{
		EntryLocalFiles: lfs,
		MediaId:         entry.Media.ID,
		CurrentProgress: simklEntry.GetProgressSafe(),
	})
	entry.EntryLibraryData = libraryData

	// +---------------------+
	// |       Animap        |
	// +---------------------+

	// Fetch AniDB data and cache it for 30 minutes
	animeMetadata, err := opts.MetadataProviderRef.Get().GetAnimeMetadata(metadata.MediaPlatform, opts.MediaId)
	if err != nil {

		// +---------------- Start
		// +---------------------+
		// |   Without Animap    |
		// +---------------------+

		// If Animap data is not found, we will still create the Entry without it
		simpleAnimeEntry, err := NewSimpleEntry(ctx, &NewSimpleAnimeEntryOptions{
			MediaId:             opts.MediaId,
			LocalFiles:          opts.LocalFiles,
			AnimeCollection:     opts.AnimeCollection,
			PlatformRef:         opts.PlatformRef,
			MetadataProviderRef: opts.MetadataProviderRef,
		})
		if err != nil {
			return nil, err
		}

		event := &AnimeEntryEvent{
			Entry: &Entry{
				MediaId:             simpleAnimeEntry.MediaId,
				Media:               simpleAnimeEntry.Media,
				EntryListData:       simpleAnimeEntry.EntryListData,
				EntryLibraryData:    simpleAnimeEntry.EntryLibraryData,
				EntryDownloadInfo:   nil,
				Episodes:            simpleAnimeEntry.Episodes,
				NextEpisode:         simpleAnimeEntry.NextEpisode,
				LocalFiles:          simpleAnimeEntry.LocalFiles,
				AnidbId:             0,
				CurrentEpisodeCount: simpleAnimeEntry.CurrentEpisodeCount,
			},
		}
		err = hook.GlobalHookManager.OnAnimeEntry().Trigger(event)
		if err != nil {
			return nil, err
		}

		return event.Entry, nil
		// +--------------- End

	}

	entry.AnidbId = animeMetadata.GetMappings().AnidbId

	// Instantiate EntryListData
	// If the media exist in the user's anime list, add the details
	if found {
		entry.EntryListData = NewEntryListData(simklEntry)
	}

	// +---------------------+
	// |       Episodes      |
	// +---------------------+

	// Create episode entities
	entry.hydrateEntryEpisodeData(simklEntry, animeMetadata, opts.MetadataProviderRef)

	event := &AnimeEntryEvent{
		Entry: entry,
	}
	err = hook.GlobalHookManager.OnAnimeEntry().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.Entry, nil
}

//----------------------------------------------------------------------------------------------------------------------

// hydrateEntryEpisodeData
// Metadata, Media and LocalFiles should be defined
func (e *Entry) hydrateEntryEpisodeData(
	simklEntry *mediaapi.AnimeListEntry,
	animeMetadata *metadata.AnimeMetadata,
	metadataProviderRef *util.Ref[metadata_provider.Provider],
) {

	if animeMetadata.Episodes == nil && len(animeMetadata.Episodes) == 0 {
		return
	}

	// +---------------------+
	// |     Discrepancy     |
	// +---------------------+

	// We offset the progress number by 1 if there is a discrepancy
	progressOffset := 0
	if FindDiscrepancy(e.Media, animeMetadata) == DiscrepancySIMKLCountsEpisodeZero {
		progressOffset = 1

		_, ok := lo.Find(e.LocalFiles, func(lf *LocalFile) bool {
			return lf.Metadata.Episode == 0
		})
		// Remove the offset if episode 0 is not found
		if !ok {
			progressOffset = 0
		}
	}

	amw := metadataProviderRef.Get().GetAnimeMetadataWrapper(simklEntry.Media, animeMetadata)

	// +---------------------+
	// |       Episodes      |
	// +---------------------+

	p := pool.NewWithResults[*Episode]()
	for _, lf := range e.LocalFiles {
		p.Go(func() *Episode {
			return NewEpisode(&NewEpisodeOptions{
				LocalFile:            lf,
				MetadataWrapper:      amw,
				OptionalAniDBEpisode: "",
				AnimeMetadata:        animeMetadata,
				Media:                e.Media,
				ProgressOffset:       progressOffset,
				IsDownloaded:         true,
				MetadataProvider:     metadataProviderRef.Get(),
			})
		})
	}
	episodes := p.Wait()
	// Sort by progress number
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
	})
	e.Episodes = episodes

	// +---------------------+
	// |    Download Info    |
	// +---------------------+

	info, err := NewEntryDownloadInfo(&NewEntryDownloadInfoOptions{
		LocalFiles:          e.LocalFiles,
		AnimeMetadata:       animeMetadata,
		Progress:            simklEntry.Progress,
		Status:              simklEntry.Status,
		Media:               e.Media,
		MetadataProviderRef: metadataProviderRef,
	})
	if err == nil {
		e.EntryDownloadInfo = info
	}

	nextEp, found := e.FindNextEpisode()
	if found {
		e.NextEpisode = nextEp
	}

}

func NewEntryListData(simklEntry *mediaapi.AnimeListEntry) *EntryListData {
	return &EntryListData{
		Progress:    simklEntry.GetProgressSafe(),
		Score:       simklEntry.GetScoreSafe(),
		Status:      simklEntry.Status,
		Repeat:      simklEntry.GetRepeatSafe(),
		StartedAt:   mediaapi.FuzzyDateToString(simklEntry.StartedAt),
		CompletedAt: mediaapi.FuzzyDateToString(simklEntry.CompletedAt),
	}
}

//----------------------------------------------------------------------------------------------------------------------

type Discrepancy int

const (
	DiscrepancySIMKLCountsEpisodeZero Discrepancy = iota
	DiscrepancySIMKLCountsSpecials
	DiscrepancyAniDBHasMore
	DiscrepancyNone
)

// FindDiscrepancy returns the discrepancy between the SIMKL and AniDB episode counts.
// It returns DiscrepancySIMKLCountsEpisodeZero if SIMKL most likely has episode 0 as part of the main count.
// It returns DiscrepancySIMKLCountsSpecials if there is a discrepancy between the SIMKL and AniDB episode counts and specials are included in the SIMKL count.
// It returns DiscrepancyAniDBHasMore if the AniDB episode count is greater than the SIMKL episode count.
// It returns DiscrepancyNone if there is no discrepancy.
func FindDiscrepancy(media *mediaapi.BaseAnime, animeMetadata *metadata.AnimeMetadata) Discrepancy {
	if media == nil || animeMetadata == nil || animeMetadata.Episodes == nil {
		return DiscrepancyNone
	}

	_, aniDBHasS1 := animeMetadata.Episodes["S1"]
	_, aniDBHasS2 := animeMetadata.Episodes["S2"]

	difference := media.GetCurrentEpisodeCount() - animeMetadata.GetMainEpisodeCount()

	if difference == 0 {
		return DiscrepancyNone
	}

	if difference < 0 {
		return DiscrepancyAniDBHasMore
	}

	if difference == 1 && aniDBHasS1 {
		return DiscrepancySIMKLCountsEpisodeZero
	}

	if difference > 1 && aniDBHasS1 && aniDBHasS2 {
		return DiscrepancySIMKLCountsSpecials
	}

	return DiscrepancyNone
}
