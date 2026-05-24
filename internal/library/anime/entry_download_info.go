package anime

import (
	"errors"
	"fmt"
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata"
	"seall/internal/api/metadata_provider"
	"seall/internal/hook"
	"seall/internal/util"
	"strconv"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

type (
	// EntryDownloadInfo is instantiated by the Entry
	EntryDownloadInfo struct {
		EpisodesToDownload    []*EntryDownloadEpisode `json:"episodesToDownload"`
		CanBatch              bool                    `json:"canBatch"`
		BatchAll              bool                    `json:"batchAll"`
		HasInaccurateSchedule bool                    `json:"hasInaccurateSchedule"`
		Rewatch               bool                    `json:"rewatch"`
		AbsoluteOffset        int                     `json:"absoluteOffset"`
	}

	EntryDownloadEpisode struct {
		EpisodeNumber int      `json:"episodeNumber"`
		AniDBEpisode  string   `json:"aniDBEpisode"`
		Episode       *Episode `json:"episode"`
	}
)

type (
	NewEntryDownloadInfoOptions struct {
		// Media's local files
		LocalFiles          []*LocalFile
		AnimeMetadata       *metadata.AnimeMetadata
		Media               *mediaapi.BaseAnime
		Progress            *int
		Status              *mediaapi.MediaListStatus
		MetadataProviderRef *util.Ref[metadata_provider.Provider]
	}
)

// NewEntryDownloadInfo returns a list of episodes to download or episodes for the torrent/debrid streaming views
// based on the options provided.
func NewEntryDownloadInfo(opts *NewEntryDownloadInfoOptions) (*EntryDownloadInfo, error) {

	reqEvent := &AnimeEntryDownloadInfoRequestedEvent{
		LocalFiles:        opts.LocalFiles,
		AnimeMetadata:     opts.AnimeMetadata,
		Media:             opts.Media,
		Progress:          opts.Progress,
		Status:            opts.Status,
		EntryDownloadInfo: &EntryDownloadInfo{},
	}

	err := hook.GlobalHookManager.OnAnimeEntryDownloadInfoRequested().Trigger(reqEvent)
	if err != nil {
		return nil, err
	}

	if reqEvent.DefaultPrevented {
		return reqEvent.EntryDownloadInfo, nil
	}

	opts.LocalFiles = reqEvent.LocalFiles
	opts.AnimeMetadata = reqEvent.AnimeMetadata
	opts.Media = reqEvent.Media
	opts.Progress = reqEvent.Progress
	opts.Status = reqEvent.Status

	if *opts.Media.Status == mediaapi.MediaStatusNotYetReleased {
		return &EntryDownloadInfo{}, nil
	}
	if opts.AnimeMetadata == nil {
		return nil, errors.New("could not get anime metadata")
	}
	currentEpisodeCount := opts.Media.GetCurrentEpisodeCount()
	if currentEpisodeCount == -1 && opts.AnimeMetadata != nil {
		currentEpisodeCount = opts.AnimeMetadata.GetCurrentEpisodeCount()
	}
	if currentEpisodeCount == -1 {
		return nil, errors.New("could not get current media episode count")
	}

	// +---------------------+
	// |     Discrepancy     |
	// +---------------------+

	// Whether SIMKL includes episode 0 as part of main episodes, but AniDB does not, however AniDB has "S1"
	discrepancy := FindDiscrepancy(opts.Media, opts.AnimeMetadata)

	// SIMKL is the source of truth for episode numbers
	epSlice := newEpisodeSlice(currentEpisodeCount)

	// Handle discrepancies
	if discrepancy != DiscrepancyNone {

		// If SIMKL includes episode 0 as part of main episodes, but AniDB does not, however AniDB has "S1"
		if discrepancy == DiscrepancySIMKLCountsEpisodeZero {
			// Add "S1" to the beginning of the episode slice
			epSlice.trimEnd(1)
			epSlice.prepend(0, "S1")
		}

		// If SIMKL includes specials, but AniDB does not
		if discrepancy == DiscrepancySIMKLCountsSpecials {
			diff := currentEpisodeCount - opts.AnimeMetadata.GetMainEpisodeCount()
			epSlice.trimEnd(diff)
			for i := 0; i < diff; i++ {
				epSlice.add(currentEpisodeCount-i, "S"+strconv.Itoa(i+1))
			}
		}

		// If AniDB has more episodes than SIMKL
		if discrepancy == DiscrepancyAniDBHasMore {
			// Do nothing
		}

	}

	// Filter out episodes not aired
	if opts.Media.NextAiringEpisode != nil {
		epSlice.filter(func(item *episodeSliceItem, index int) bool {
			// e.g. if the next airing episode is 13, then filter out episodes 14 and above
			return index+1 < opts.Media.NextAiringEpisode.Episode
		})
	}

	// Get progress, if the media isn't in the user's list, progress is 0
	// If the media is completed, set progress is 0
	progress := 0
	if opts.Progress != nil {
		progress = *opts.Progress
	}
	if opts.Status != nil {
		if *opts.Status == mediaapi.MediaListStatusCompleted {
			progress = 0
		}
	}

	hasInaccurateSchedule := false
	if opts.Media.NextAiringEpisode == nil && *opts.Media.Status == mediaapi.MediaStatusReleasing {
		hasInaccurateSchedule = true
	}

	// Filter out episodes already watched (index+1 is the progress number)
	toDownloadSlice := epSlice.filterNew(func(item *episodeSliceItem, index int) bool {
		return index+1 > progress
	})

	// This slice contains episode numbers that are not downloaded
	// The source of truth is AniDB, but we will handle discrepancies
	lfsEpSlice := newEpisodeSlice(0)
	if opts.LocalFiles != nil {
		// Get all episode numbers of main local files
		for _, lf := range opts.LocalFiles {
			if lf.Metadata.Type == LocalFileTypeMain {
				lfsEpSlice.add(lf.Metadata.Episode, lf.Metadata.AniDBEpisode)
			}
		}
	}

	// Filter out downloaded episodes
	toDownloadSlice.filter(func(item *episodeSliceItem, index int) bool {
		isDownloaded := false
		for _, lf := range opts.LocalFiles {
			if lf.Metadata.Type != LocalFileTypeMain {
				continue
			}
			// If the file episode number matches that of the episode slice item
			if lf.Metadata.Episode == item.episodeNumber {
				isDownloaded = true
			}
			// If the slice episode number is 0 and the file is a main S1
			if discrepancy == DiscrepancySIMKLCountsEpisodeZero && item.episodeNumber == 0 && lf.Metadata.AniDBEpisode == "S1" {
				isDownloaded = true
			}
		}

		return !isDownloaded
	})

	// +---------------------+
	// |   EntryEpisode      |
	// +---------------------+

	// Generate `episodesToDownload` based on `toDownloadSlice`

	// DEVNOTE: The EntryEpisode generated has inaccurate progress numbers since not local files are passed in

	mediaWrapper := opts.MetadataProviderRef.Get().GetAnimeMetadataWrapper(opts.Media, opts.AnimeMetadata)

	progressOffset := 0
	if discrepancy == DiscrepancySIMKLCountsEpisodeZero {
		progressOffset = 1
	}

	p := pool.NewWithResults[*EntryDownloadEpisode]()
	for _, ep := range toDownloadSlice.getSlice() {
		p.Go(func() *EntryDownloadEpisode {
			str := new(EntryDownloadEpisode)
			str.EpisodeNumber = ep.episodeNumber
			str.AniDBEpisode = ep.aniDBEpisode
			// Create a new episode with a placeholder local file
			// We pass that placeholder local file so that all episodes are hydrated as main episodes for consistency
			str.Episode = NewEpisode(&NewEpisodeOptions{
				MetadataWrapper: mediaWrapper,
				LocalFile: &LocalFile{
					ParsedData:       &LocalFileParsedData{},
					ParsedFolderData: []*LocalFileParsedData{},
					Metadata: &LocalFileMetadata{
						Episode:      ep.episodeNumber,
						Type:         LocalFileTypeMain,
						AniDBEpisode: ep.aniDBEpisode,
					},
				},
				OptionalAniDBEpisode: str.AniDBEpisode,
				AnimeMetadata:        opts.AnimeMetadata,
				Media:                opts.Media,
				ProgressOffset:       progressOffset,
				IsDownloaded:         false,
				MetadataProvider:     opts.MetadataProviderRef.Get(),
			})
			str.Episode.AniDBEpisode = ep.aniDBEpisode
			// Reset the local file to nil, since it's a placeholder
			str.Episode.LocalFile = nil
			return str
		})
	}
	episodesToDownload := p.Wait()

	//--------------

	canBatch := false
	if *opts.Media.GetStatus() == mediaapi.MediaStatusFinished && opts.Media.GetTotalEpisodeCount() > 0 {
		canBatch = true
	}
	batchAll := false
	if canBatch && lfsEpSlice.len() == 0 && progress == 0 {
		batchAll = true
	}
	rewatch := false
	if opts.Status != nil && *opts.Status == mediaapi.MediaListStatusCompleted {
		rewatch = true
	}

	downloadInfo := &EntryDownloadInfo{
		EpisodesToDownload:    episodesToDownload,
		CanBatch:              canBatch,
		BatchAll:              batchAll,
		Rewatch:               rewatch,
		HasInaccurateSchedule: hasInaccurateSchedule,
		AbsoluteOffset:        opts.AnimeMetadata.GetOffset(),
	}

	event := &AnimeEntryDownloadInfoEvent{
		EntryDownloadInfo: downloadInfo,
	}
	err = hook.GlobalHookManager.OnAnimeEntryDownloadInfo().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.EntryDownloadInfo, nil
}

type episodeSliceItem struct {
	episodeNumber int
	aniDBEpisode  string
}

type episodeSlice []*episodeSliceItem

func newEpisodeSlice(episodeCount int) *episodeSlice {
	s := make([]*episodeSliceItem, 0)
	for i := 0; i < episodeCount; i++ {
		s = append(s, &episodeSliceItem{episodeNumber: i + 1, aniDBEpisode: strconv.Itoa(i + 1)})
	}
	ret := &episodeSlice{}
	ret.set(s)
	return ret
}

func (s *episodeSlice) set(eps []*episodeSliceItem) {
	*s = eps
}

func (s *episodeSlice) add(episodeNumber int, aniDBEpisode string) {
	*s = append(*s, &episodeSliceItem{episodeNumber: episodeNumber, aniDBEpisode: aniDBEpisode})
}

func (s *episodeSlice) prepend(episodeNumber int, aniDBEpisode string) {
	*s = append([]*episodeSliceItem{{episodeNumber: episodeNumber, aniDBEpisode: aniDBEpisode}}, *s...)
}

func (s *episodeSlice) trimEnd(n int) {
	*s = (*s)[:len(*s)-n]
}

func (s *episodeSlice) trimStart(n int) {
	*s = (*s)[n:]
}

func (s *episodeSlice) len() int {
	return len(*s)
}

func (s *episodeSlice) get(index int) *episodeSliceItem {
	return (*s)[index]
}

func (s *episodeSlice) getEpisodeNumber(episodeNumber int) *episodeSliceItem {
	for _, item := range *s {
		if item.episodeNumber == episodeNumber {
			return item
		}
	}
	return nil
}

func (s *episodeSlice) filter(filter func(*episodeSliceItem, int) bool) {
	*s = lo.Filter(*s, filter)
}

func (s *episodeSlice) filterNew(filter func(*episodeSliceItem, int) bool) *episodeSlice {
	s2 := make(episodeSlice, 0)
	for i, item := range *s {
		if filter(item, i) {
			s2 = append(s2, item)
		}
	}
	return &s2
}

func (s *episodeSlice) copy() *episodeSlice {
	s2 := make(episodeSlice, len(*s), cap(*s))
	for i, item := range *s {
		s2[i] = item
	}
	return &s2
}

func (s *episodeSlice) getSlice() []*episodeSliceItem {
	return *s
}

func (s *episodeSlice) print() {
	for i, item := range *s {
		fmt.Printf("(%d) %d -> %s\n", i, item.episodeNumber, item.aniDBEpisode)
	}
}
