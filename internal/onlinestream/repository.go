package onlinestream

import (
	"context"
	"errors"
	"fmt"
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata"
	"seall/internal/api/metadata_provider"
	"seall/internal/database/db"
	"seall/internal/extension"
	hibikeonlinestream "seall/internal/extension/hibike/onlinestream"
	"seall/internal/library/anime"
	"seall/internal/platforms/platform"
	"seall/internal/util"
	"seall/internal/util/filecache"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type (
	Repository struct {
		logger              *zerolog.Logger
		extensionBankRef    *util.Ref[*extension.UnifiedBank]
		fileCacher          *filecache.Cacher
		metadataProviderRef *util.Ref[metadata_provider.Provider]
		platformRef         *util.Ref[platform.Platform]
		simklBaseAnimeCache *mediaapi.BaseAnimeCache
		db                  *db.Database
	}
)

var (
	ErrNoVideoSourceFound = errors.New("no video source found")
)

type (
	Episode struct {
		Number       int            `json:"number"`
		SeasonNumber int            `json:"seasonNumber,omitempty"`
		Title        string         `json:"title,omitempty"`
		Image        string         `json:"image,omitempty"`
		Description  string         `json:"description,omitempty"`
		IsFiller     bool           `json:"isFiller,omitempty"`
		Metadata     *anime.Episode `json:"metadata"`
	}

	EpisodeSource struct {
		Number       int            `json:"number"`
		VideoSources []*VideoSource `json:"videoSources"`
	}

	VideoSource struct {
		Server    string                             `json:"server"`
		Headers   map[string]string                  `json:"headers,omitempty"`
		URL       string                             `json:"url"`
		Label     string                             `json:"label,omitempty"`
		Quality   string                             `json:"quality"`
		Type      hibikeonlinestream.VideoSourceType `json:"type"`
		Subtitles []*Subtitle                        `json:"subtitles,omitempty"`
	}

	EpisodeListResponse struct {
		Episodes []*Episode          `json:"episodes"`
		Media    *mediaapi.BaseAnime `json:"media"`
	}

	Subtitle struct {
		URL      string `json:"url"`
		Language string `json:"language"`
	}
)

type (
	NewRepositoryOptions struct {
		Logger              *zerolog.Logger
		FileCacher          *filecache.Cacher
		MetadataProviderRef *util.Ref[metadata_provider.Provider]
		PlatformRef         *util.Ref[platform.Platform]
		Database            *db.Database
		ExtensionBankRef    *util.Ref[*extension.UnifiedBank]
	}
)

func NewRepository(opts *NewRepositoryOptions) *Repository {
	return &Repository{
		logger:              opts.Logger,
		metadataProviderRef: opts.MetadataProviderRef,
		fileCacher:          opts.FileCacher,
		extensionBankRef:    opts.ExtensionBankRef,
		simklBaseAnimeCache: mediaapi.NewBaseAnimeCache(),
		platformRef:         opts.PlatformRef,
		db:                  opts.Database,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// getFcEpisodeDataBucket returns a episode data bucket for the provider and mediaId.
// "Episode data" refers to the episodeData struct
//
//	e.g., onlinestream_zoro_episode-data_123
func (r *Repository) getFcEpisodeDataBucket(provider string, mediaId int) filecache.Bucket {
	return filecache.NewBucket("onlinestream_"+provider+"_episode-data_"+strconv.Itoa(mediaId), time.Hour*24*2)
}

// getFcEpisodeListBucket returns a episode data bucket for the provider and mediaId.
// "Episode list" refers to a slice of onlinestream_providers.EpisodeDetails
//
//	e.g., onlinestream_zoro_episode-list_123
func (r *Repository) getFcEpisodeListBucket(provider string, mediaId int) filecache.Bucket {
	return filecache.NewBucket("onlinestream_"+provider+"_episode-data_"+strconv.Itoa(mediaId), time.Hour*24*1)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *Repository) getMedia(ctx context.Context, mId int) (*mediaapi.BaseAnime, error) {
	media, err := r.simklBaseAnimeCache.GetOrSet(mId, func() (*mediaapi.BaseAnime, error) {
		media, err := r.platformRef.Get().GetAnime(ctx, mId)
		if err != nil {
			return nil, err
		}
		return media, nil
	})
	if err != nil {
		return nil, err
	}
	return media, nil
}

func (r *Repository) GetMedia(ctx context.Context, mId int) (*mediaapi.BaseAnime, error) {
	return r.getMedia(ctx, mId)
}

func (r *Repository) EmptyCache(mediaId int) error {
	_ = r.fileCacher.RemoveAllBy(func(filename string) bool {
		return strings.HasPrefix(filename, "onlinestream_") && strings.Contains(filename, strconv.Itoa(mediaId))
	})
	// clear all stores
	_ = r.fileCacher.Clear()
	return nil
}

func (r *Repository) GetMediaEpisodes(provider string, media *mediaapi.BaseAnime, dubbed bool) ([]*Episode, error) {
	episodes := make([]*Episode, 0)

	if provider == "" {
		return episodes, nil
	}

	// +---------------------+
	// |       Animap        |
	// +---------------------+

	//animeMetadata, err := r.metadataProvider.GetAnimeMetadata(metadata.MediaPlatform, mId)
	//	//foundAnimeMetadata := err == nil && animeMetadata != nil
	//aw := r.metadataProvider.GetAnimeMetadataWrapper(media, animeMetadata)

	var episodeCollection *anime.EpisodeCollection
	var err error
	mediaType := MediaTypeForBaseMedia(media)
	if mediaType == hibikeonlinestream.MediaTypeAnime {
		episodeCollection, err = anime.NewEpisodeCollection(anime.NewEpisodeCollectionOptions{
			AnimeMetadata:       nil,
			Media:               media,
			MetadataProviderRef: r.metadataProviderRef,
			Logger:              r.logger,
		})
	}
	foundEpisodeCollection := err == nil && episodeCollection != nil
	var animeMetadata *metadata.AnimeMetadata
	if foundEpisodeCollection {
		animeMetadata = episodeCollection.Metadata
	}

	// +---------------------+
	// |    Episode list     |
	// +---------------------+

	// Fetch the episode list from the provider
	// "from" and "to" are set to 0 in order not to fetch episode servers
	ec, err := r.getEpisodeContainer(provider, media, 0, 0, dubbed, media.GetStartYearSafe(), 0, nil)
	if err != nil {
		return nil, err
	}

	for _, episodeDetails := range ec.ProviderEpisodeList {
		seasonNumber := providerEpisodeSeasonNumber(animeMetadata, episodeDetails)

		// If the title contains "[{", it means it's an episode part (e.g. "Episode 6 [{6.5}]", the episode number should be 6)
		if strings.Contains(episodeDetails.Title, "[{") {
			ep := strings.Split(episodeDetails.Title, "[{")[1]
			ep = strings.Split(ep, "}]")[0]
			episodes = append(episodes, &Episode{
				Number:       episodeDetails.Number,
				SeasonNumber: seasonNumber,
				Title:        fmt.Sprintf("Episode %s", ep),
				Image:        media.GetBannerImageSafe(),
				Description:  "",
				IsFiller:     false,
			})

		} else {

			if foundEpisodeCollection {
				episode, found := episodeCollection.FindEpisodeByNumber(episodeDetails.Number)
				if found {
					episodes = append(episodes, &Episode{
						Number:       episodeDetails.Number,
						SeasonNumber: seasonNumber,
						Title:        episode.EpisodeTitle,
						Image:        episode.EpisodeMetadata.Image,
						Description:  episode.EpisodeMetadata.Summary,
						IsFiller:     episode.EpisodeMetadata.IsFiller,
						Metadata:     episode,
					})
				} else {
					episodes = append(episodes, &Episode{
						Number:       episodeDetails.Number,
						SeasonNumber: seasonNumber,
						Title:        episodeDetails.Title,
						Image:        media.GetCoverImageSafe(),
					})
				}
			} else {
				episodes = append(episodes, &Episode{
					Number:       episodeDetails.Number,
					SeasonNumber: seasonNumber,
					Title:        episodeDetails.Title,
					Image:        media.GetCoverImageSafe(),
				})
			}

		}
	}

	episodes = lo.Filter(episodes, func(item *Episode, index int) bool {
		return item != nil
	})

	return episodes, nil
}

func (r *Repository) GetEpisodeSources(ctx context.Context, provider string, mId int, number int, seasonNumber int, dubbed bool, year int) (*EpisodeSource, error) {

	// +---------------------+
	// |        Media        |
	// +---------------------+

	media, err := r.getMedia(ctx, mId)
	if err != nil {
		return nil, err
	}

	// +---------------------+
	// |   Episode servers   |
	// +---------------------+

	var animeMetadata *metadata.AnimeMetadata
	if MediaTypeForBaseMedia(media) == hibikeonlinestream.MediaTypeAnime {
		if episodeCollection, err := anime.NewEpisodeCollection(anime.NewEpisodeCollectionOptions{
			AnimeMetadata:       nil,
			Media:               media,
			MetadataProviderRef: r.metadataProviderRef,
			Logger:              r.logger,
		}); err == nil && episodeCollection != nil {
			animeMetadata = episodeCollection.Metadata
		}
	}

	ec, err := r.getEpisodeContainer(provider, media, number, number, dubbed, year, seasonNumber, animeMetadata)
	if err != nil {
		return nil, err
	}

	var sources *EpisodeSource
	for _, ep := range ec.Episodes {
		if ep.Number == number {
			s := &EpisodeSource{
				Number:       ep.Number,
				VideoSources: make([]*VideoSource, 0),
			}
			for _, es := range ep.Servers {

				for _, vs := range es.VideoSources {
					s.VideoSources = append(s.VideoSources, &VideoSource{
						Server:  es.Server,
						Headers: es.Headers,
						URL:     vs.URL,
						Label:   vs.Label,
						Quality: vs.Quality,
						Type:    vs.Type,
						Subtitles: lo.Map(vs.Subtitles, func(sub *hibikeonlinestream.VideoSubtitle, _ int) *Subtitle {
							return &Subtitle{
								URL:      sub.URL,
								Language: sub.Language,
							}
						}),
					})
				}
			}
			sources = s
			break
		}
	}

	if sources == nil {
		return nil, ErrNoVideoSourceFound
	}

	return sources, nil
}
