package anime

import (
	"context"
	"errors"
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata"
	"seall/internal/api/metadata_provider"
	"seall/internal/platforms/platform"
	"seall/internal/util"
	"sort"
	"strconv"

	"github.com/sourcegraph/conc/pool"
)

type (
	SimpleEntry struct {
		MediaId             int                 `json:"mediaId"`
		Media               *mediaapi.BaseAnime `json:"media"`
		EntryListData       *EntryListData      `json:"listData"`
		EntryLibraryData    *EntryLibraryData   `json:"libraryData"`
		Episodes            []*Episode          `json:"episodes"`
		NextEpisode         *Episode            `json:"nextEpisode"`
		LocalFiles          []*LocalFile        `json:"localFiles"`
		CurrentEpisodeCount int                 `json:"currentEpisodeCount"`
	}

	SimpleEntryListData struct {
		Progress    int                       `json:"progress,omitempty"`
		Score       float64                   `json:"score,omitempty"`
		Status      *mediaapi.MediaListStatus `json:"status,omitempty"`
		StartedAt   string                    `json:"startedAt,omitempty"`
		CompletedAt string                    `json:"completedAt,omitempty"`
	}

	NewSimpleAnimeEntryOptions struct {
		MediaId             int
		LocalFiles          []*LocalFile // All local files
		AnimeCollection     *mediaapi.AnimeCollection
		PlatformRef         *util.Ref[platform.Platform]
		MetadataProviderRef *util.Ref[metadata_provider.Provider]
	}
)

func NewSimpleEntry(ctx context.Context, opts *NewSimpleAnimeEntryOptions) (*SimpleEntry, error) {

	if opts.AnimeCollection == nil ||
		opts.PlatformRef.IsAbsent() {
		return nil, errors.New("missing arguments when creating simple media entry")
	}
	// Create new Entry
	entry := new(SimpleEntry)
	entry.MediaId = opts.MediaId

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
		fetchedMedia, err := opts.PlatformRef.Get().GetAnime(ctx, opts.MediaId)
		if err != nil {
			return nil, err
		}
		entry.Media = fetchedMedia
	} else {
		entry.Media = simklEntry.Media
	}

	entry.CurrentEpisodeCount = entry.Media.GetCurrentEpisodeCount()

	// +---------------------+
	// |   Local files       |
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

	// Instantiate EntryListData
	// If the media exist in the user's anime list, add the details
	if found {
		entry.EntryListData = &EntryListData{
			Progress:    simklEntry.GetProgressSafe(),
			Score:       simklEntry.GetScoreSafe(),
			Status:      simklEntry.Status,
			Repeat:      simklEntry.GetRepeatSafe(),
			StartedAt:   mediaapi.ToEntryStartDate(simklEntry.StartedAt),
			CompletedAt: mediaapi.ToEntryCompletionDate(simklEntry.CompletedAt),
		}
	}

	// +---------------------+
	// |       Episodes      |
	// +---------------------+

	amw := opts.MetadataProviderRef.Get().GetAnimeMetadataWrapper(simklEntry.Media, nil)

	// Create episode entities
	entry.hydrateEntryEpisodeData(amw)

	return entry, nil

}

//----------------------------------------------------------------------------------------------------------------------

// hydrateEntryEpisodeData
// Metadata, Media and LocalFiles should be defined
func (e *SimpleEntry) hydrateEntryEpisodeData(amw metadata_provider.AnimeMetadataWrapper) {

	// +---------------------+
	// |       Episodes      |
	// +---------------------+

	p := pool.NewWithResults[*Episode]()
	for _, lf := range e.LocalFiles {
		lf := lf
		p.Go(func() *Episode {
			return NewSimpleEpisode(&NewSimpleEpisodeOptions{
				LocalFile:       lf,
				Media:           e.Media,
				IsDownloaded:    true,
				MetadataWrapper: amw,
			})
		})
	}
	episodes := p.Wait()
	// Sort by progress number
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
	})
	e.Episodes = episodes

	nextEp, found := e.FindNextEpisode()
	if found {
		e.NextEpisode = nextEp
	}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewAnimeMetadataFromEntry(media *mediaapi.BaseAnime, episodes []*Episode) *metadata.AnimeMetadata {
	animeMetadata := &metadata.AnimeMetadata{
		Titles:       make(map[string]string),
		Episodes:     make(map[string]*metadata.EpisodeMetadata),
		EpisodeCount: 0,
		SpecialCount: 0,
		Mappings: &metadata.AnimeMappings{
			MediaId: media.GetID(),
		},
	}
	animeMetadata.Titles["en"] = media.GetTitleSafe()
	animeMetadata.Titles["x-jat"] = media.GetRomajiTitleSafe()

	// Hydrate episodes
	for _, episode := range episodes {
		animeMetadata.Episodes[episode.FileMetadata.AniDBEpisode] = &metadata.EpisodeMetadata{
			AnidbId:               0,
			TvdbId:                0,
			Title:                 episode.DisplayTitle,
			Image:                 episode.EpisodeMetadata.Image,
			AirDate:               "",
			Length:                0,
			Summary:               "",
			Overview:              "",
			EpisodeNumber:         episode.EpisodeNumber,
			Episode:               strconv.Itoa(episode.EpisodeNumber),
			SeasonNumber:          0,
			AbsoluteEpisodeNumber: episode.EpisodeNumber,
			AnidbEid:              0,
			HasImage:              true,
		}
		animeMetadata.EpisodeCount++
	}

	return animeMetadata
}

func NewAnimeMetadataFromEpisodeCount(media *mediaapi.BaseAnime, episodes []int) *metadata.AnimeMetadata {
	animeMetadata := &metadata.AnimeMetadata{
		Titles:       make(map[string]string),
		Episodes:     make(map[string]*metadata.EpisodeMetadata),
		EpisodeCount: 0,
		SpecialCount: 0,
		Mappings: &metadata.AnimeMappings{
			MediaId: media.GetID(),
		},
	}
	animeMetadata.Titles["en"] = media.GetTitleSafe()
	animeMetadata.Titles["x-jat"] = media.GetRomajiTitleSafe()

	// Hydrate episodes
	for _, episode := range episodes {
		animeMetadata.Episodes[strconv.Itoa(episode)] = &metadata.EpisodeMetadata{
			AnidbId:               0,
			TvdbId:                0,
			Title:                 media.GetTitleSafe(),
			Image:                 media.GetBannerImageSafe(),
			AirDate:               "",
			Length:                0,
			Summary:               "",
			Overview:              "",
			EpisodeNumber:         episode,
			Episode:               strconv.Itoa(episode),
			SeasonNumber:          0,
			AbsoluteEpisodeNumber: episode,
			AnidbEid:              0,
			HasImage:              true,
		}
		animeMetadata.EpisodeCount++
	}

	return animeMetadata
}
