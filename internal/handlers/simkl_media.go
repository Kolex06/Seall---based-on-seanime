package handlers

import (
	"context"
	"errors"
	"fmt"
	"seall/internal/api/mediaapi"
	simklapi "seall/internal/api/simkl"
	"seall/internal/platforms/shared_platform"
	"seall/internal/util/result"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// HandleGetMediaCollection
//
//	@summary returns the user's SIMKL anime collection.
//	@desc Calling GET will return the cached anime collection.
//	@desc The manga collection is also refreshed in the background, and upon completion, a WebSocket event is sent.
//	@desc Calling POST will refetch both the anime and manga collections.
//	@returns mediaapi.AnimeCollection
//	@route /api/v1/simkl/collection [GET,POST]
func (h *Handler) HandleGetMediaCollection(c echo.Context) error {

	bypassCache := c.Request().Method == "POST"

	if !bypassCache {
		// Get the user's simkl collection
		animeCollection, err := h.App.GetMediaCollection(false)
		if err != nil {
			return h.RespondWithError(c, err)
		}
		return h.RespondWithData(c, animeCollection)
	}

	animeCollection, err := h.App.RefreshAnimeCollection()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	go func() {
		_, _ = h.App.RefreshMangaCollection()
	}()

	return h.RespondWithData(c, animeCollection)
}

// HandleGetRawMediaCollection
//
//	@summary returns the user's SIMKL anime collection without filtering out custom lists.
//	@desc Calling GET will return the cached anime collection.
//	@returns mediaapi.AnimeCollection
//	@route /api/v1/simkl/collection/raw [GET,POST]
func (h *Handler) HandleGetRawMediaCollection(c echo.Context) error {

	bypassCache := c.Request().Method == "POST"

	// Get the user's simkl collection
	animeCollection, err := h.App.GetRawMediaCollection(bypassCache)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, animeCollection)
}

var tagsCache *mediaapi.MediaTagMap
var simklGenreTagsCache = result.NewCache[int, []string]()

// HandleGetRawMediaCollectionTags
//
//	@summary returns the SIMKL genres for the user's raw media collection.
//	@desc The response keeps the tag-map shape used by the lists page filters, but the values come from SIMKL genres.
//	@returns mediaapi.MediaTagMap
//	@route /api/v1/simkl/collection/raw/tags [GET]
func (h *Handler) HandleGetRawMediaCollectionTags(c echo.Context) error {
	h.App.OnRefreshMediaCollectionFuncs.Set("HandleGetRawMediaCollectionTags", func() {
		tagsCache = nil
	})

	if tagsCache != nil {
		return h.RespondWithData(c, *tagsCache)
	}

	collection, err := h.App.GetRawMediaCollection(false)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	tags := h.mediaTagMapFromSimklGenres(c.Request().Context(), collection)
	tagsCache = &tags

	return h.RespondWithData(c, tags)
}

func (h *Handler) mediaTagMapFromSimklGenres(ctx context.Context, collection *mediaapi.AnimeCollection) mediaapi.MediaTagMap {
	tags := mediaapi.MediaTagMapFromAnimeCollectionGenres(collection)
	if collection == nil || collection.GetMediaListCollection() == nil {
		return tags
	}

	client := h.App.SimklClientRef.Get()
	if client == nil {
		return tags
	}

	type genreJob struct {
		id   int
		kind simklapi.MediaType
	}

	queued := make(map[int]struct{})
	jobs := make(chan genreJob)
	var wg sync.WaitGroup
	var mu sync.Mutex

	addGenres := func(mediaID int, genres []string) {
		if len(genres) == 0 {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for _, genre := range genres {
			addMediaTag(tags, mediaID, genre)
		}
	}

	workerCount := 8
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				media, err := client.MediaDetails(ctx, job.kind, strconv.Itoa(job.id), "full")
				if err != nil || media == nil {
					continue
				}
				genres := cleanGenreList(media.Genres)
				simklGenreTagsCache.SetT(job.id, genres, 24*time.Hour)
				addGenres(job.id, genres)
			}
		}()
	}

	for _, list := range collection.GetMediaListCollection().GetLists() {
		if list == nil {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry == nil || entry.GetMedia() == nil {
				continue
			}
			media := entry.GetMedia()
			mediaID := media.GetID()
			if mediaID == 0 {
				continue
			}
			if genres, ok := simklGenreTagsCache.Get(mediaID); ok {
				addGenres(mediaID, genres)
				continue
			}
			if _, ok := queued[mediaID]; ok {
				continue
			}
			queued[mediaID] = struct{}{}
			jobs <- genreJob{
				id:   mediaID,
				kind: simklKindForBaseAnime(media),
			}
		}
	}

	close(jobs)
	wg.Wait()

	return tags
}

func simklKindForBaseAnime(media *mediaapi.BaseAnime) simklapi.MediaType {
	if media == nil {
		return simklapi.MediaTypeShows
	}
	if cached, ok := simklapi.CachedDiscoveryMedia(media.GetID()); ok && cached.Kind != "" {
		return cached.Kind
	}
	if siteURL := media.GetSiteURL(); siteURL != nil {
		site := strings.ToLower(*siteURL)
		switch {
		case strings.Contains(site, "/movie/") || strings.Contains(site, "/movies/"):
			return simklapi.MediaTypeMovies
		case strings.Contains(site, "/anime/"):
			return simklapi.MediaTypeAnime
		case strings.Contains(site, "/tv/") || strings.Contains(site, "/show/") || strings.Contains(site, "/shows/"):
			return simklapi.MediaTypeShows
		}
	}
	if format := media.GetFormat(); format != nil && *format == mediaapi.MediaFormatMovie {
		return simklapi.MediaTypeMovies
	}
	return simklapi.MediaTypeShows
}

func cleanGenreList(genres []string) []string {
	ret := make([]string, 0, len(genres))
	seen := make(map[string]struct{}, len(genres))
	for _, genre := range genres {
		genre = strings.TrimSpace(genre)
		if genre == "" {
			continue
		}
		key := strings.ToLower(genre)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		ret = append(ret, genre)
	}
	return ret
}

func addMediaTag(tags mediaapi.MediaTagMap, mediaID int, tagName string) {
	tagName = strings.TrimSpace(tagName)
	if mediaID == 0 || tagName == "" {
		return
	}
	for _, current := range tags[mediaID] {
		if strings.EqualFold(current, tagName) {
			return
		}
	}
	tags[mediaID] = append(tags[mediaID], tagName)
}

// HandleEditMediaListEntry
//
//	@summary updates the user's list entry on Simkl.
//	@desc This is used to edit an entry on SIMKL.
//	@desc The "type" field is used to determine if the entry is an anime or manga and refreshes the collection accordingly.
//	@desc The client should refetch collection-dependent queries after this mutation.
//	@returns true
//	@route /api/v1/simkl/list-entry [POST]
func (h *Handler) HandleEditMediaListEntry(c echo.Context) error {

	type body struct {
		MediaId   *int                      `json:"mediaId"`
		Status    *mediaapi.MediaListStatus `json:"status"`
		Score     *int                      `json:"score"`
		Progress  *int                      `json:"progress"`
		StartDate *mediaapi.FuzzyDateInput  `json:"startedAt"`
		EndDate   *mediaapi.FuzzyDateInput  `json:"completedAt"`
		Type      string                    `json:"type"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	err := h.App.MediaPlatformRef.Get().UpdateEntry(
		c.Request().Context(),
		*p.MediaId,
		p.Status,
		p.Score,
		p.Progress,
		p.StartDate,
		p.EndDate,
	)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	switch p.Type {
	case "anime":
		_, _ = h.App.RefreshAnimeCollection()
	case "manga":
		_, _ = h.App.RefreshMangaCollection()
	default:
		_, _ = h.App.RefreshAnimeCollection()
		_, _ = h.App.RefreshMangaCollection()
	}

	return h.RespondWithData(c, true)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

var (
	detailsCache = result.NewCache[int, *mediaapi.AnimeDetailsById_Media]()
)

// HandleGetMediaDetails
//
//	@summary returns more details about an SIMKL anime entry.
//	@desc This fetches more fields omitted from the base queries.
//	@param id - int - true - "The SIMKL anime ID"
//	@returns mediaapi.AnimeDetailsById_Media
//	@route /api/v1/simkl/media-details/{id} [GET]
func (h *Handler) HandleGetMediaDetails(c echo.Context) error {

	mId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.RespondWithError(c, err)
	}

	if details, ok := detailsCache.Get(mId); ok {
		return h.RespondWithData(c, details)
	}
	details, err := h.App.MediaPlatformRef.Get().GetAnimeDetails(c.Request().Context(), mId)
	if err != nil {
		return h.RespondWithError(c, err)
	}
	detailsCache.Set(mId, details)

	return h.RespondWithData(c, details)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

var studioDetailsMap = result.NewMap[int, *mediaapi.StudioDetails]()

// HandleGetStudioDetails
//
//	@summary returns details about a studio.
//	@desc This fetches media produced by the studio.
//	@param id - int - true - "The SIMKL studio ID"
//	@returns mediaapi.StudioDetails
//	@route /api/v1/simkl/studio-details/{id} [GET]
func (h *Handler) HandleGetStudioDetails(c echo.Context) error {

	mId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.RespondWithError(c, err)
	}

	if details, ok := studioDetailsMap.Get(mId); ok {
		return h.RespondWithData(c, details)
	}
	details, err := h.App.MediaPlatformRef.Get().GetStudioDetails(c.Request().Context(), mId)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	go func() {
		if details != nil {
			studioDetailsMap.Set(mId, details)
		}
	}()

	return h.RespondWithData(c, details)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

// HandleDeleteMediaListEntry
//
//	@summary deletes an entry from the user's SIMKL list.
//	@desc This is used to delete an entry on SIMKL.
//	@desc The "type" field is used to determine if the entry is an anime or manga and refreshes the collection accordingly.
//	@desc The client should refetch collection-dependent queries after this mutation.
//	@route /api/v1/simkl/list-entry [DELETE]
//	@returns bool
func (h *Handler) HandleDeleteMediaListEntry(c echo.Context) error {

	type body struct {
		MediaId *int    `json:"mediaId"`
		Type    *string `json:"type"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	if p.Type == nil || p.MediaId == nil {
		return h.RespondWithError(c, errors.New("missing parameters"))
	}

	var listEntryID int

	switch *p.Type {
	case "anime":
		// Get the list entry ID
		animeCollection, err := h.App.GetMediaCollection(false)
		if err != nil {
			return h.RespondWithError(c, err)
		}

		listEntry, found := animeCollection.GetListEntryFromAnimeId(*p.MediaId)
		if !found {
			return h.RespondWithError(c, errors.New("list entry not found"))
		}
		listEntryID = listEntry.ID
	case "manga":
		// Get the list entry ID
		mangaCollection, err := h.App.GetMangaCollection(false)
		if err != nil {
			return h.RespondWithError(c, err)
		}

		listEntry, found := mangaCollection.GetListEntryFromMangaId(*p.MediaId)
		if !found {
			return h.RespondWithError(c, errors.New("list entry not found"))
		}
		listEntryID = listEntry.ID
	}

	// Delete the list entry
	err := h.App.MediaPlatformRef.Get().DeleteEntry(c.Request().Context(), *p.MediaId, listEntryID)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	switch *p.Type {
	case "anime":
		_, _ = h.App.RefreshAnimeCollection()
	case "manga":
		_, _ = h.App.RefreshMangaCollection()
	}

	return h.RespondWithData(c, true)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	simklListAnimeCache       = result.NewCache[string, *mediaapi.ListAnime]()
	simklListRecentAnimeCache = result.NewCache[string, *mediaapi.ListRecentAnime]() // holds 1 value
)

type listMediaRequest struct {
	Page                *int                    `json:"page,omitempty"`
	Search              *string                 `json:"search,omitempty"`
	PerPage             *int                    `json:"perPage,omitempty"`
	Sort                []*mediaapi.MediaSort   `json:"sort,omitempty"`
	Status              []*mediaapi.MediaStatus `json:"status,omitempty"`
	Genres              []*string               `json:"genres,omitempty"`
	Tags                []*string               `json:"tags,omitempty"`
	AverageScoreGreater *int                    `json:"averageScore_greater,omitempty"`
	Season              *mediaapi.MediaSeason   `json:"season,omitempty"`
	SeasonYear          *int                    `json:"seasonYear,omitempty"`
	Format              *mediaapi.MediaFormat   `json:"format,omitempty"`
	IsAdult             *bool                   `json:"isAdult,omitempty"`
	CountryOfOrigin     *string                 `json:"countryOfOrigin,omitempty"`
}

// HandleListMedia
//
//	@summary returns a list of anime based on the search parameters.
//	@desc This is used by the "Discover" and "Advanced Search".
//	@route /api/v1/simkl/list-media [POST]
//	@returns mediaapi.ListAnime
func (h *Handler) HandleListMedia(c echo.Context) error {

	p := new(listMediaRequest)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	page := intValue(p.Page, 1)
	perPage := intValue(p.PerPage, 20)

	var isAdult *bool = nil
	if p.IsAdult != nil {
		allowedAdult := *p.IsAdult && h.App.Settings.GetMediaApi().EnableAdultContent
		isAdult = &allowedAdult
	}
	p.IsAdult = isAdult

	cacheKey := mediaapi.ListAnimeCacheKey(
		&page,
		p.Search,
		&perPage,
		p.Sort,
		p.Status,
		p.Genres,
		p.Tags,
		p.AverageScoreGreater,
		p.Season,
		p.SeasonYear,
		p.Format,
		isAdult,
		p.CountryOfOrigin,
	)

	cached, ok := simklListAnimeCache.Get(cacheKey)
	if ok {
		return h.RespondWithData(c, cached)
	}

	ret, err := h.listMediaViaSimklREST(c.Request().Context(), *p)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	if ret != nil {
		simklListAnimeCache.SetT(cacheKey, ret, time.Minute*10)
	}

	return h.RespondWithData(c, ret)
}

// HandleListRecentAiringMedia
//
//	@summary returns a list of recently aired anime.
//	@desc This is used by the "Schedule" page to display recently aired anime.
//	@route /api/v1/simkl/list-recent-airing-media [POST]
//	@returns mediaapi.ListRecentAnime
func (h *Handler) HandleListRecentAiringMedia(c echo.Context) error {

	type body struct {
		Page            *int                   `json:"page,omitempty"`
		Search          *string                `json:"search,omitempty"`
		PerPage         *int                   `json:"perPage,omitempty"`
		AiringAtGreater *int                   `json:"airingAt_greater,omitempty"`
		AiringAtLesser  *int                   `json:"airingAt_lesser,omitempty"`
		NotYetAired     *bool                  `json:"notYetAired,omitempty"`
		Sort            []*mediaapi.AiringSort `json:"sort,omitempty"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	page := intValue(p.Page, 1)
	perPage := intValue(p.PerPage, 50)

	cacheKey := fmt.Sprintf("%v-%v-%v-%v-%v-%v-%v", page, p.Search, perPage, p.AiringAtGreater, p.AiringAtLesser, p.NotYetAired, p.Sort)

	cached, ok := simklListRecentAnimeCache.Get(cacheKey)
	if ok {
		return h.RespondWithData(c, cached)
	}

	ret, err := h.listRecentAiringMediaViaSimklREST(c.Request().Context(), page, perPage, p.Search)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	simklListRecentAnimeCache.SetT(cacheKey, ret, time.Hour*1)

	return h.RespondWithData(c, ret)
}

func (h *Handler) listMediaViaSimklREST(ctx context.Context, request listMediaRequest) (*mediaapi.ListAnime, error) {
	client := h.App.SimklClientRef.Get()
	page := intValue(request.Page, 1)
	perPage := intValue(request.PerPage, 20)
	kinds := mediaKindsForFormat(request.Format)
	term := ""
	if request.Search != nil {
		term = strings.TrimSpace(*request.Search)
	}

	items := make([]simklapi.DiscoveryMedia, 0, perPage*len(kinds))
	fetchLimit := page * perPage
	if term == "" {
		fetchLimit = maxInt(fetchLimit*4, 80)
	}
	for _, kind := range kinds {
		var ret []simklapi.DiscoveryMedia
		var err error
		if term != "" {
			ret, err = client.SearchMedia(ctx, kind, term)
		} else {
			ret, err = client.TrendingMedia(ctx, kind, 1, fetchLimit)
		}
		if err != nil {
			return nil, err
		}
		items = append(items, ret...)
	}
	if term == "" && wantsUnreleasedStatus(request.Status) {
		calendarItems, err := h.getSIMKLCalendarItemsPartial(ctx, time.Now(), kinds, calendarDiscoveryMonths(page, perPage))
		if err != nil && len(items) == 0 {
			return nil, err
		}
		if err == nil {
			items = append(items, simklapi.DiscoveryMediaFromCalendarItems(calendarItems)...)
		}
	}

	items = dedupeDiscoveryMedia(items)
	items = filterDiscoveryMedia(items, request)
	if !sortDiscoveryMedia(items, request.Sort) && wantsUnreleasedStatus(request.Status) {
		sortDiscoveryMediaByRelease(items, false)
	}

	return listAnimeFromDiscovery(items, page, perPage), nil
}

func (h *Handler) getSIMKLCalendarItemsPartial(ctx context.Context, now time.Time, mediaTypes []simklapi.MediaType, monthCount int) ([]simklapi.CalendarItem, error) {
	client := h.App.SimklClientRef.Get()
	if client == nil {
		return nil, errors.New("simkl calendar client is not available")
	}
	if monthCount < 1 {
		monthCount = 1
	}
	if len(mediaTypes) == 0 {
		mediaTypes = []simklapi.MediaType{
			simklapi.MediaTypeShows,
			simklapi.MediaTypeAnime,
			simklapi.MediaTypeMovies,
		}
	}

	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	ret := make([]simklapi.CalendarItem, 0)
	var firstErr error
	for i := 0; i < monthCount; i++ {
		month := start.AddDate(0, i, 0)
		for _, mediaType := range mediaTypes {
			items, err := client.MonthlyCalendar(ctx, mediaType, month.Year(), int(month.Month()))
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			ret = append(ret, items...)
		}
	}
	if len(ret) == 0 && firstErr != nil {
		return nil, firstErr
	}
	return ret, nil
}

func (h *Handler) listRecentAiringMediaViaSimklREST(ctx context.Context, page int, perPage int, search *string) (*mediaapi.ListRecentAnime, error) {
	client := h.App.SimklClientRef.Get()
	term := ""
	if search != nil {
		term = strings.TrimSpace(*search)
	}

	items := make([]simklapi.DiscoveryMedia, 0, perPage)
	if term != "" {
		ret, err := client.SearchMedia(ctx, simklapi.MediaTypeAnime, term)
		if err != nil {
			return nil, err
		}
		items = append(items, ret...)
	} else {
		ret, err := client.TrendingMedia(ctx, simklapi.MediaTypeAnime, 1, page*perPage)
		if err != nil {
			return nil, err
		}
		items = append(items, ret...)
	}

	return listRecentAnimeFromDiscovery(dedupeDiscoveryMedia(items), page, perPage), nil
}

func mediaKindsForFormat(format *mediaapi.MediaFormat) []simklapi.MediaType {
	if format == nil {
		return []simklapi.MediaType{simklapi.MediaTypeMovies, simklapi.MediaTypeShows, simklapi.MediaTypeAnime}
	}
	switch *format {
	case mediaapi.MediaFormatMovie:
		return []simklapi.MediaType{simklapi.MediaTypeMovies}
	case mediaapi.MediaFormatTv, mediaapi.MediaFormatTvShort:
		return []simklapi.MediaType{simklapi.MediaTypeShows, simklapi.MediaTypeAnime}
	case mediaapi.MediaFormatManga, mediaapi.MediaFormatNovel, mediaapi.MediaFormatOneShot:
		return []simklapi.MediaType{simklapi.MediaTypeMovies, simklapi.MediaTypeShows, simklapi.MediaTypeAnime}
	default:
		return []simklapi.MediaType{simklapi.MediaTypeAnime}
	}
}

func listAnimeFromDiscovery(items []simklapi.DiscoveryMedia, page int, perPage int) *mediaapi.ListAnime {
	pageItems, pageInfo := paginateDiscoveryMedia(items, page, perPage)
	media := make([]*mediaapi.BaseAnime, 0, len(pageItems))
	for _, item := range pageItems {
		converted := simklapi.ToBaseAnime(item.Kind, &item.Media)
		if converted != nil {
			media = append(media, converted)
		}
	}

	return &mediaapi.ListAnime{
		Page: &mediaapi.ListAnime_Page{
			Media:    media,
			PageInfo: pageInfo,
		},
	}
}

func listRecentAnimeFromDiscovery(items []simklapi.DiscoveryMedia, page int, perPage int) *mediaapi.ListRecentAnime {
	pageItems, animePageInfo := paginateDiscoveryMedia(items, page, perPage)
	recentPageInfo := &mediaapi.ListRecentAnime_Page_PageInfo{
		CurrentPage: animePageInfo.CurrentPage,
		HasNextPage: animePageInfo.HasNextPage,
		LastPage:    animePageInfo.LastPage,
		PerPage:     animePageInfo.PerPage,
		Total:       animePageInfo.Total,
	}
	schedules := make([]*mediaapi.ListRecentAnime_Page_AiringSchedules, 0, len(pageItems))
	now := int(time.Now().Unix())
	for _, item := range pageItems {
		media := simklapi.ToBaseAnime(item.Kind, &item.Media)
		if media == nil {
			continue
		}
		episode := 1
		if media.Episodes != nil && *media.Episodes > 0 {
			episode = *media.Episodes
		}
		schedules = append(schedules, &mediaapi.ListRecentAnime_Page_AiringSchedules{
			ID:              media.ID,
			AiringAt:        now,
			Episode:         episode,
			Media:           media,
			TimeUntilAiring: 0,
		})
	}

	return &mediaapi.ListRecentAnime{
		Page: &mediaapi.ListRecentAnime_Page{
			AiringSchedules: schedules,
			PageInfo:        recentPageInfo,
		},
	}
}

func paginateDiscoveryMedia(items []simklapi.DiscoveryMedia, page int, perPage int) ([]simklapi.DiscoveryMedia, *mediaapi.ListAnime_Page_PageInfo) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	total := len(items)
	lastPage := 1
	if total > 0 {
		lastPage = (total + perPage - 1) / perPage
	}
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}
	hasNextPage := page < lastPage

	return items[start:end], &mediaapi.ListAnime_Page_PageInfo{
		CurrentPage: &page,
		HasNextPage: &hasNextPage,
		LastPage:    &lastPage,
		PerPage:     &perPage,
		Total:       &total,
	}
}

func dedupeDiscoveryMedia(items []simklapi.DiscoveryMedia) []simklapi.DiscoveryMedia {
	seen := make(map[string]struct{}, len(items))
	ret := make([]simklapi.DiscoveryMedia, 0, len(items))
	for _, item := range items {
		id := item.Media.IDs.PrimarySimklID()
		key := fmt.Sprintf("%s:%d:%s:%d", item.Kind, id, strings.ToLower(item.Media.Title), item.Media.Year)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		ret = append(ret, item)
	}
	return ret
}

func filterDiscoveryMedia(items []simklapi.DiscoveryMedia, request listMediaRequest) []simklapi.DiscoveryMedia {
	ret := make([]simklapi.DiscoveryMedia, 0, len(items))
	wantedGenres := simklGenreFilters(request.Genres, request.Tags)
	for _, item := range items {
		base := simklapi.ToBaseAnime(item.Kind, &item.Media)
		if base == nil {
			continue
		}
		if !matchesMediaStatus(base, request.Status) {
			continue
		}
		if !matchesMediaFormat(base, request.Format) {
			continue
		}
		if !matchesReleaseYear(base, request.SeasonYear) {
			continue
		}
		if !matchesMediaSeason(base, request.Season) {
			continue
		}
		if !matchesGenres(item.Media.Genres, wantedGenres) {
			continue
		}
		if !matchesAverageScore(item.Media.Ratings.Simkl, request.AverageScoreGreater) {
			continue
		}
		if !matchesCountry(item.Media.Country, request.CountryOfOrigin) {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func simklGenreFilters(genres []*string, tags []*string) []*string {
	if len(tags) == 0 {
		return genres
	}
	ret := make([]*string, 0, len(genres)+len(tags))
	ret = append(ret, genres...)
	ret = append(ret, tags...)
	return ret
}

func matchesMediaStatus(media *mediaapi.BaseAnime, statuses []*mediaapi.MediaStatus) bool {
	if len(statuses) == 0 {
		return true
	}
	for _, status := range statuses {
		if status != nil && media.Status != nil && *media.Status == *status {
			return true
		}
	}
	return false
}

func matchesMediaFormat(media *mediaapi.BaseAnime, format *mediaapi.MediaFormat) bool {
	if format == nil {
		return true
	}
	return media.Format != nil && *media.Format == *format
}

func matchesReleaseYear(media *mediaapi.BaseAnime, year *int) bool {
	if year == nil || *year == 0 {
		return true
	}
	if media.SeasonYear != nil && *media.SeasonYear == *year {
		return true
	}
	return media.GetStartDate().GetYear() != nil && *media.GetStartDate().GetYear() == *year
}

func matchesMediaSeason(media *mediaapi.BaseAnime, season *mediaapi.MediaSeason) bool {
	if season == nil {
		return true
	}
	month := media.GetStartDate().GetMonth()
	if month == nil || *month == 0 {
		return true
	}
	switch *season {
	case mediaapi.MediaSeasonWinter:
		return *month >= 1 && *month <= 3
	case mediaapi.MediaSeasonSpring:
		return *month >= 4 && *month <= 6
	case mediaapi.MediaSeasonSummer:
		return *month >= 7 && *month <= 9
	case mediaapi.MediaSeasonFall:
		return *month >= 10 && *month <= 12
	default:
		return true
	}
}

func matchesGenres(mediaGenres []string, wanted []*string) bool {
	if len(wanted) == 0 {
		return true
	}
	genreSet := make(map[string]struct{}, len(mediaGenres))
	for _, genre := range mediaGenres {
		genre = strings.ToLower(strings.TrimSpace(genre))
		if genre != "" {
			genreSet[genre] = struct{}{}
		}
	}
	for _, wantedGenre := range wanted {
		if wantedGenre == nil || strings.TrimSpace(*wantedGenre) == "" {
			continue
		}
		if _, ok := genreSet[strings.ToLower(strings.TrimSpace(*wantedGenre))]; !ok {
			return false
		}
	}
	return true
}

func matchesAverageScore(rating *simklapi.Rating, threshold *int) bool {
	if threshold == nil || *threshold <= 0 {
		return true
	}
	if rating == nil || rating.Rating <= 0 {
		return false
	}
	return int(rating.Rating*10) >= *threshold
}

func matchesCountry(country string, wanted *string) bool {
	if wanted == nil || strings.TrimSpace(*wanted) == "" {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(country), strings.TrimSpace(*wanted))
}

func sortDiscoveryMedia(items []simklapi.DiscoveryMedia, sorts []*mediaapi.MediaSort) bool {
	if len(sorts) == 0 || sorts[0] == nil {
		return false
	}
	sortKey := strings.ToUpper(string(*sorts[0]))
	switch sortKey {
	case "SCORE_DESC", "AVERAGE_SCORE_DESC":
		sort.SliceStable(items, func(i, j int) bool {
			return simklRatingScore(items[i]) > simklRatingScore(items[j])
		})
		return true
	case "SCORE", "AVERAGE_SCORE":
		sort.SliceStable(items, func(i, j int) bool {
			return simklRatingScore(items[i]) < simklRatingScore(items[j])
		})
		return true
	case "START_DATE_DESC":
		sortDiscoveryMediaByRelease(items, true)
		return true
	case "START_DATE":
		sortDiscoveryMediaByRelease(items, false)
		return true
	}
	return false
}

func simklRatingScore(item simklapi.DiscoveryMedia) float64 {
	if item.Media.Ratings.Simkl == nil {
		return 0
	}
	return item.Media.Ratings.Simkl.Rating
}

func discoveryYear(item simklapi.DiscoveryMedia) int {
	media := simklapi.ToBaseAnime(item.Kind, &item.Media)
	if media == nil {
		return 0
	}
	if media.SeasonYear != nil {
		return *media.SeasonYear
	}
	if media.GetStartDate().GetYear() != nil {
		return *media.GetStartDate().GetYear()
	}
	return 0
}

func sortDiscoveryMediaByRelease(items []simklapi.DiscoveryMedia, desc bool) {
	sort.SliceStable(items, func(i, j int) bool {
		left := discoveryDateKey(items[i])
		right := discoveryDateKey(items[j])
		if desc {
			return left > right
		}
		return left < right
	})
}

func discoveryDateKey(item simklapi.DiscoveryMedia) int {
	media := simklapi.ToBaseAnime(item.Kind, &item.Media)
	if media == nil {
		return 0
	}
	year := 0
	month := 1
	day := 1
	if media.GetStartDate().GetYear() != nil {
		year = *media.GetStartDate().GetYear()
	} else if media.SeasonYear != nil {
		year = *media.SeasonYear
	}
	if media.GetStartDate().GetMonth() != nil {
		month = *media.GetStartDate().GetMonth()
	}
	if media.GetStartDate().GetDay() != nil {
		day = *media.GetStartDate().GetDay()
	}
	return year*10000 + month*100 + day
}

func wantsUnreleasedStatus(statuses []*mediaapi.MediaStatus) bool {
	for _, status := range statuses {
		if status != nil && *status == mediaapi.MediaStatusNotYetReleased {
			return true
		}
	}
	return false
}

func calendarDiscoveryMonths(page int, perPage int) int {
	months := page * perPage / 20
	if months < 6 {
		return 6
	}
	if months > 12 {
		return 12
	}
	return months
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func intValue(value *int, fallback int) int {
	if value == nil || *value <= 0 {
		return fallback
	}
	return *value
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var simklMissedSequelsCache = result.NewCache[int, []*mediaapi.BaseAnime]()

// HandleListMissedSequels
//
//	@summary returns a list of sequels not in the user's list.
//	@desc This is used by the "Discover" page to display sequels the user may have missed.
//	@route /api/v1/simkl/list-missed-sequels [GET]
//	@returns []mediaapi.BaseAnime
func (h *Handler) HandleListMissedSequels(c echo.Context) error {

	cached, ok := simklMissedSequelsCache.Get(1)
	if ok {
		return h.RespondWithData(c, cached)
	}

	// Get complete anime collection
	animeCollection, err := h.App.MediaPlatformRef.Get().GetMediaCollectionWithRelations(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}

	ret, err := mediaapi.ListMissedSequels(
		h.App.MediaPlatformRef.Get().GetMediaApiClient(),
		animeCollection,
		h.App.Logger,
		h.App.GetUserSimklToken(),
	)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	simklMissedSequelsCache.SetT(1, ret, time.Hour*4)

	return h.RespondWithData(c, ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var simklStatsCache = result.NewCache[int, *mediaapi.Stats]()

// HandleGetMediaStats
//
//	@summary returns the simkl stats.
//	@desc This returns the SIMKL stats for the user.
//	@route /api/v1/simkl/stats [GET]
//	@returns mediaapi.Stats
func (h *Handler) HandleGetMediaStats(c echo.Context) error {
	cached, ok := simklStatsCache.Get(0)
	if ok {
		return h.RespondWithData(c, cached)
	}

	stats, err := h.App.MediaPlatformRef.Get().GetViewerStats(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}

	ret, err := mediaapi.GetStats(
		c.Request().Context(),
		stats,
	)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	simklStatsCache.SetT(0, ret, time.Hour*1)

	return h.RespondWithData(c, ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// HandleGetMediaCacheLayerStatus
//
//	@summary returns the status of the SIMKL cache layer.
//	@desc This returns the status of the SIMKL cache layer.
//	@route /api/v1/simkl/cache-layer/status [GET]
//	@returns bool
func (h *Handler) HandleGetMediaCacheLayerStatus(c echo.Context) error {
	return h.RespondWithData(c, shared_platform.IsWorking.Load())
}

// HandleToggleMediaCacheLayerStatus
//
//	@summary toggles the status of the SIMKL cache layer.
//	@desc This toggles the status of the SIMKL cache layer.
//	@route /api/v1/simkl/cache-layer/status [POST]
//	@returns bool
func (h *Handler) HandleToggleMediaCacheLayerStatus(c echo.Context) error {
	shared_platform.IsWorking.Store(!shared_platform.IsWorking.Load())
	return h.RespondWithData(c, shared_platform.IsWorking.Load())
}
