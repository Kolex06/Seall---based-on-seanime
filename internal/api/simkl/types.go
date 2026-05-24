package simkl

import "github.com/goccy/go-json"

type MediaType string

const (
	MediaTypeAll    MediaType = ""
	MediaTypeShows  MediaType = "shows"
	MediaTypeMovies MediaType = "movies"
	MediaTypeAnime  MediaType = "anime"
)

type WatchStatus string

const (
	WatchStatusWatching    WatchStatus = "watching"
	WatchStatusPlanToWatch WatchStatus = "plantowatch"
	WatchStatusCompleted   WatchStatus = "completed"
	WatchStatusDropped     WatchStatus = "dropped"
	WatchStatusHold        WatchStatus = "hold"
)

type ImageKind string

const (
	ImageKindPoster  ImageKind = "posters"
	ImageKindFanart  ImageKind = "fanart"
	ImageKindAvatar  ImageKind = "avatars"
	ImageKindEpisode ImageKind = "episodes"
)

type ImageSize string

const (
	ImageSizePosterCard   ImageSize = "_c"
	ImageSizePosterMedium ImageSize = "_m"
	ImageSizePosterWide   ImageSize = "_w"
	ImageSizeFanartWide   ImageSize = "_medium"
	ImageSizeAvatarMedium ImageSize = "_100"
	ImageSizeEpisodeWide  ImageSize = "_w"
)

type IDs struct {
	Simkl       int    `json:"simkl,omitempty"`
	SimklID     int    `json:"simkl_id,omitempty"`
	Slug        string `json:"slug,omitempty"`
	IMDB        string `json:"imdb,omitempty"`
	TMDB        string `json:"tmdb,omitempty"`
	TVDB        string `json:"tvdb,omitempty"`
	MAL         string `json:"mal,omitempty"`
	AniDB       string `json:"anidb,omitempty"`
	Kitsu       string `json:"kitsu,omitempty"`
	Crunchyroll string `json:"crunchyroll,omitempty"`
}

func (ids *IDs) UnmarshalJSON(data []byte) error {
	type rawIDs struct {
		Simkl       interface{} `json:"simkl,omitempty"`
		SimklID     interface{} `json:"simkl_id,omitempty"`
		Slug        string      `json:"slug,omitempty"`
		IMDB        interface{} `json:"imdb,omitempty"`
		TMDB        interface{} `json:"tmdb,omitempty"`
		TVDB        interface{} `json:"tvdb,omitempty"`
		MAL         interface{} `json:"mal,omitempty"`
		AniDB       interface{} `json:"anidb,omitempty"`
		Kitsu       interface{} `json:"kitsu,omitempty"`
		Crunchyroll interface{} `json:"crunchyroll,omitempty"`
	}

	var raw rawIDs
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	ids.Simkl = flexibleInt(raw.Simkl)
	ids.SimklID = flexibleInt(raw.SimklID)
	ids.Slug = raw.Slug
	ids.IMDB = flexibleString(raw.IMDB)
	ids.TMDB = flexibleString(raw.TMDB)
	ids.TVDB = flexibleString(raw.TVDB)
	ids.MAL = flexibleString(raw.MAL)
	ids.AniDB = flexibleString(raw.AniDB)
	ids.Kitsu = flexibleString(raw.Kitsu)
	ids.Crunchyroll = flexibleString(raw.Crunchyroll)
	return nil
}

func (ids IDs) PrimarySimklID() int {
	if ids.Simkl != 0 {
		return ids.Simkl
	}
	return ids.SimklID
}

type Rating struct {
	Rating float64 `json:"rating,omitempty"`
	Votes  int     `json:"votes,omitempty"`
	Rank   int     `json:"rank,omitempty"`
}

type Ratings struct {
	Simkl *Rating `json:"simkl,omitempty"`
	IMDB  *Rating `json:"imdb,omitempty"`
	MAL   *Rating `json:"mal,omitempty"`
}

type StandardMedia struct {
	Title         string      `json:"title,omitempty"`
	Year          int         `json:"year,omitempty"`
	Type          string      `json:"type,omitempty"`
	Status        string      `json:"status,omitempty"`
	To            WatchStatus `json:"to,omitempty"`
	Rating        *int        `json:"rating,omitempty"`
	AddedAt       string      `json:"added_at,omitempty"`
	WatchedAt     string      `json:"watched_at,omitempty"`
	Poster        string      `json:"poster,omitempty"`
	Fanart        string      `json:"fanart,omitempty"`
	URL           string      `json:"url,omitempty"`
	Runtime       *int        `json:"runtime,omitempty"`
	Overview      string      `json:"overview,omitempty"`
	Genres        []string    `json:"genres,omitempty"`
	Country       string      `json:"country,omitempty"`
	Released      string      `json:"released,omitempty"`
	TotalEpisodes *int        `json:"total_episodes,omitempty"`
	AnimeType     string      `json:"anime_type,omitempty"`
	EnglishName   string      `json:"en_title,omitempty"`
	IDs           IDs         `json:"ids,omitempty"`
	Ratings       Ratings     `json:"ratings,omitempty"`
}

type Episode struct {
	Number      int    `json:"number,omitempty"`
	Episode     int    `json:"episode,omitempty"`
	Season      int    `json:"season,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Aired       bool   `json:"aired,omitempty"`
	Image       string `json:"img,omitempty"`
	Date        string `json:"date,omitempty"`
	WatchedAt   string `json:"watched_at,omitempty"`
	TVDB        *struct {
		Season  int `json:"season,omitempty"`
		Episode int `json:"episode,omitempty"`
	} `json:"tvdb,omitempty"`
	IDs IDs `json:"ids,omitempty"`
}

type Season struct {
	Number   int       `json:"number,omitempty"`
	Episodes []Episode `json:"episodes,omitempty"`
}

type WatchlistItem struct {
	AddedToWatchlistAt    string         `json:"added_to_watchlist_at,omitempty"`
	LastWatchedAt         string         `json:"last_watched_at,omitempty"`
	UserRatedAt           string         `json:"user_rated_at,omitempty"`
	UserRating            *int           `json:"user_rating,omitempty"`
	Status                WatchStatus    `json:"status,omitempty"`
	LastWatched           string         `json:"last_watched,omitempty"`
	NextToWatch           string         `json:"next_to_watch,omitempty"`
	WatchedEpisodesCount  *int           `json:"watched_episodes_count,omitempty"`
	TotalEpisodesCount    *int           `json:"total_episodes_count,omitempty"`
	NotAiredEpisodesCount *int           `json:"not_aired_episodes_count,omitempty"`
	AnimeType             string         `json:"anime_type,omitempty"`
	Show                  *StandardMedia `json:"show,omitempty"`
	Movie                 *StandardMedia `json:"movie,omitempty"`
	Seasons               []Season       `json:"seasons,omitempty"`
	TVDBSeasons           []int          `json:"tvdb_seasons,omitempty"`
	Memo                  *Memo          `json:"memo,omitempty"`
}

type Memo struct {
	Text      string `json:"text,omitempty"`
	IsPrivate bool   `json:"is_private,omitempty"`
}

type AllItems struct {
	Shows  []WatchlistItem `json:"shows,omitempty"`
	Anime  []WatchlistItem `json:"anime,omitempty"`
	Movies []WatchlistItem `json:"movies,omitempty"`
}

type UserSettings struct {
	User struct {
		Name     string `json:"name,omitempty"`
		JoinedAt string `json:"joined_at,omitempty"`
		Gender   string `json:"gender,omitempty"`
		Avatar   string `json:"avatar,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Location string `json:"loc,omitempty"`
		Age      string `json:"age,omitempty"`
	} `json:"user,omitempty"`
	Account struct {
		ID       int    `json:"id,omitempty"`
		Timezone string `json:"timezone,omitempty"`
		Type     string `json:"type,omitempty"`
	} `json:"account,omitempty"`
	Connections map[string]bool `json:"connections,omitempty"`
}

type ActivityBucket struct {
	All             *string `json:"all,omitempty"`
	RatedAt         *string `json:"rated_at,omitempty"`
	Playback        *string `json:"playback,omitempty"`
	PlanToWatch     *string `json:"plantowatch,omitempty"`
	Watching        *string `json:"watching,omitempty"`
	Completed       *string `json:"completed,omitempty"`
	Hold            *string `json:"hold,omitempty"`
	Dropped         *string `json:"dropped,omitempty"`
	RemovedFromList *string `json:"removed_from_list,omitempty"`
}

type Activities struct {
	All      *string        `json:"all,omitempty"`
	Settings ActivityBucket `json:"settings,omitempty"`
	TVShows  ActivityBucket `json:"tv_shows,omitempty"`
	Anime    ActivityBucket `json:"anime,omitempty"`
	Movies   ActivityBucket `json:"movies,omitempty"`
}

type PinCode struct {
	Result          string `json:"result,omitempty"`
	DeviceCode      string `json:"device_code,omitempty"`
	UserCode        string `json:"user_code,omitempty"`
	VerificationURL string `json:"verification_url,omitempty"`
	ExpiresIn       int    `json:"expires_in,omitempty"`
	Interval        int    `json:"interval,omitempty"`
}

type PinStatus struct {
	Result      string `json:"result,omitempty"`
	Message     string `json:"message,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Error       string `json:"error,omitempty"`
}

type AddItemsRequest struct {
	Movies   []StandardMedia `json:"movies,omitempty"`
	Shows    []StandardMedia `json:"shows,omitempty"`
	Episodes []Episode       `json:"episodes,omitempty"`
}

type AddItemsResponse struct {
	Added struct {
		Movies   int `json:"movies,omitempty"`
		Shows    int `json:"shows,omitempty"`
		Episodes int `json:"episodes,omitempty"`
		Statuses []struct {
			Request  StandardMedia `json:"request,omitempty"`
			Response struct {
				Status    WatchStatus `json:"status,omitempty"`
				SimklType string      `json:"simkl_type,omitempty"`
				AnimeType string      `json:"anime_type,omitempty"`
			} `json:"response,omitempty"`
		} `json:"statuses,omitempty"`
	} `json:"added,omitempty"`
	NotFound struct {
		Movies   []StandardMedia `json:"movies,omitempty"`
		Shows    []StandardMedia `json:"shows,omitempty"`
		Episodes []Episode       `json:"episodes,omitempty"`
	} `json:"not_found,omitempty"`
}
