package shared_platform

import (
	"context"
	"errors"
	"fmt"
	"seall/internal/api/mediaapi"
	"seall/internal/events"
	"seall/internal/util"
	"seall/internal/util/filecache"
	"seall/internal/util/result"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

// devnote: I got lazy and used global variables

var ShouldCache = atomic.Bool{}
var IsWorking = atomic.Bool{}
var MediaApiClient = atomic.Value{}

type failureRecord struct {
	timestamp time.Time
	err       error
}

var (
	failureTracking      = make([]failureRecord, 0)
	failureTrackingMutex sync.RWMutex
)

const (
	failureWindow     = 30 * time.Second // time window to consider failures
	failureThreshold  = 4                // number of failures needed to mark as down
	cleanupInterval   = 5 * time.Minute  // how often to clean up old failure records
	maxFailureRecords = 50               // maximum number of failure records to keep
)

func init() {
	ShouldCache.Store(true)
	IsWorking.Store(true)

	go func() {
		// Every 10 seconds, check if the SIMKL client is working
		for {
			time.Sleep(time.Second * 10)
			if !ShouldCache.Load() {
				IsWorking.Store(true)
				continue
			}
			if IsWorking.Load() {
				continue
			}
			if MediaApiClient.Load() == nil {
				IsWorking.Store(true)
				continue
			}
			simklClient, ok := MediaApiClient.Load().(mediaapi.MediaApiClient)
			if !ok {
				IsWorking.Store(true)
				continue
			}
			_, err := simklClient.BaseAnimeByID(context.Background(), new(1))
			if err != nil {
				IsWorking.Store(false)
			} else {
				clearFailureTracking()
				events.GlobalWSEventManager.SendEvent(events.InfoToast, "The SIMKL API is back online")
				IsWorking.Store(true)
			}
		}
	}()

	// periodic cleanup of old failure records
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			cleanupOldFailures()
		}
	}()
}

type (
	// CacheLayer is a "network-first" wrapper around an SIMKL client that caches fetched data in cache files.
	// It detects when the API client is not working and falls back to the cached data instead.
	// When the API client not working, it will still send the requests in the background and transition back to working state when the API client is working again.
	// Entry/progress updates are queued when the API client is not working; other mutations return an error.
	// Caching strategy:
	// - All queries to a specific media that IS in the anime collection or manga collection will be always cached/updated without limit
	// - Media that are NOT in the anime or manga collection will be bounded to a maximum of 10 entries
	CacheLayer struct {
		mediaClientRef         *util.Ref[mediaapi.MediaApiClient]
		fileCacher             *filecache.Cacher
		buckets                map[string]filecache.PermanentBucket
		logger                 *zerolog.Logger
		collectionMediaIDs     *result.Map[int, struct{}] // Track which media IDs are in collections
		lastCollectionUpdate   time.Time                  // When collections were last fetched
		logoutFunc             func()                     // called when an invalid token is detected
		pendingUpdateSyncMutex sync.Mutex
	}
)

const (
	AnimeCollectionBucket          = "anime-collection"
	AnimeCollectionTagsBucket      = "MEDIA-LIBRARY-tags"
	AnimeCollectionRelationsBucket = "MEDIA-LIBRARY-relations"
	MangaCollectionBucket          = "manga-collection"
	MangaCollectionTagsBucket      = "manga-collection-tags"
	BaseAnimeBucket                = "base-anime"
	BaseAnimeMalBucket             = "base-anime-mal"
	CompleteAnimeBucket            = "complete-anime"
	AnimeDetailsBucket             = "anime-details"
	BaseMangaBucket                = "base-manga"
	MangaDetailsBucket             = "manga-details"
	ViewerBucket                   = "viewer"
	ViewerStatsBucket              = "viewer-stats"
	StudioDetailsBucket            = "studio-details"
	AnimeAiringScheduleBucket      = "anime-airing-schedule"
	AnimeAiringScheduleRawBucket   = "anime-airing-schedule-raw"
	ListMediaBucket                = "list-media"
	ListRecentAiringMediaBucket    = "list-recent-airing-media"
	SearchBaseMangaBucket          = "search-base-manga"
	ListMangaBucket                = "list-manga"
	SearchBaseAnimeByIdsBucket     = "search-base-anime-by-ids"
	CustomQueryBucket              = "custom-query"
	PendingMediaListUpdatesBucket  = "pending-media-list-updates"

	maxNonCollectionCacheEntries      = 10
	maxNonCollectionMediaCacheEntries = 50
	// Collection update interval (refresh collection tracking every 30 minutes)
	collectionUpdateInterval = 30 * time.Minute
)

// addFailureRecord adds a new failure record to the tracking
func addFailureRecord(err error) {
	failureTrackingMutex.Lock()
	defer failureTrackingMutex.Unlock()

	now := time.Now()
	failureTracking = append(failureTracking, failureRecord{
		timestamp: now,
		err:       err,
	})

	// keep only the most recent records
	if len(failureTracking) > maxFailureRecords {
		failureTracking = failureTracking[len(failureTracking)-maxFailureRecords:]
	}
}

// getRecentFailureCount returns the number of failures within the failure window
func getRecentFailureCount() int {
	failureTrackingMutex.RLock()
	defer failureTrackingMutex.RUnlock()

	now := time.Now()
	cutoff := now.Add(-failureWindow)
	count := 0

	for _, record := range failureTracking {
		if record.timestamp.After(cutoff) {
			count++
		}
	}

	return count
}

// cleanupOldFailures removes failure records older than the failure window
func cleanupOldFailures() {
	failureTrackingMutex.Lock()
	defer failureTrackingMutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-failureWindow)
	validRecords := make([]failureRecord, 0, len(failureTracking))

	for _, record := range failureTracking {
		if record.timestamp.After(cutoff) {
			validRecords = append(validRecords, record)
		}
	}

	failureTracking = validRecords
}

// clearFailureTracking clears all failure records (called when API comes back online)
func clearFailureTracking() {
	failureTrackingMutex.Lock()
	defer failureTrackingMutex.Unlock()
	failureTracking = failureTracking[:0]
}

// NewCacheLayer returns a new instance of the global cache layer.
// An optional logoutFunc can be passed to perform server-side cleanup when an invalid token is detected.
func NewCacheLayer(mediaClientRef *util.Ref[mediaapi.MediaApiClient], logoutFunc ...func()) mediaapi.MediaApiClient {
	fileCacher, err := filecache.NewCacher(mediaClientRef.Get().GetCacheDir())
	if err != nil {
		return mediaClientRef.Get()
	}

	buckets := make(map[string]filecache.PermanentBucket)
	buckets[AnimeCollectionBucket] = filecache.NewPermanentBucket(AnimeCollectionBucket)
	buckets[AnimeCollectionTagsBucket] = filecache.NewPermanentBucket(AnimeCollectionTagsBucket)
	buckets[AnimeCollectionRelationsBucket] = filecache.NewPermanentBucket(AnimeCollectionRelationsBucket)
	buckets[MangaCollectionBucket] = filecache.NewPermanentBucket(MangaCollectionBucket)
	buckets[MangaCollectionTagsBucket] = filecache.NewPermanentBucket(MangaCollectionTagsBucket)
	buckets[BaseAnimeBucket] = filecache.NewPermanentBucket(BaseAnimeBucket)
	buckets[BaseAnimeMalBucket] = filecache.NewPermanentBucket(BaseAnimeMalBucket)
	buckets[CompleteAnimeBucket] = filecache.NewPermanentBucket(CompleteAnimeBucket)
	buckets[AnimeDetailsBucket] = filecache.NewPermanentBucket(AnimeDetailsBucket)
	buckets[BaseMangaBucket] = filecache.NewPermanentBucket(BaseMangaBucket)
	buckets[MangaDetailsBucket] = filecache.NewPermanentBucket(MangaDetailsBucket)
	buckets[ViewerBucket] = filecache.NewPermanentBucket(ViewerBucket)
	buckets[ViewerStatsBucket] = filecache.NewPermanentBucket(ViewerStatsBucket)
	buckets[StudioDetailsBucket] = filecache.NewPermanentBucket(StudioDetailsBucket)
	buckets[AnimeAiringScheduleBucket] = filecache.NewPermanentBucket(AnimeAiringScheduleBucket)
	buckets[AnimeAiringScheduleRawBucket] = filecache.NewPermanentBucket(AnimeAiringScheduleRawBucket)
	buckets[ListMediaBucket] = filecache.NewPermanentBucket(ListMediaBucket)
	buckets[ListRecentAiringMediaBucket] = filecache.NewPermanentBucket(ListRecentAiringMediaBucket)
	buckets[SearchBaseMangaBucket] = filecache.NewPermanentBucket(SearchBaseMangaBucket)
	buckets[ListMangaBucket] = filecache.NewPermanentBucket(ListMangaBucket)
	buckets[SearchBaseAnimeByIdsBucket] = filecache.NewPermanentBucket(SearchBaseAnimeByIdsBucket)
	buckets[CustomQueryBucket] = filecache.NewPermanentBucket(CustomQueryBucket)
	buckets[PendingMediaListUpdatesBucket] = filecache.NewPermanentBucket(PendingMediaListUpdatesBucket)

	logger := util.NewLogger()

	var logout func()
	if len(logoutFunc) > 0 {
		logout = logoutFunc[0]
	}

	cl := &CacheLayer{
		mediaClientRef:     mediaClientRef,
		fileCacher:         fileCacher,
		buckets:            buckets,
		logger:             logger,
		collectionMediaIDs: result.NewMap[int, struct{}](),
		logoutFunc:         logout,
	}

	MediaApiClient.Store(mediaClientRef.Get())
	cl.startQueuedUpdateSync()

	return cl
}

var _ mediaapi.MediaApiClient = (*CacheLayer)(nil)

func (c *CacheLayer) IsAuthenticated() bool {
	return c.mediaClientRef.Get().IsAuthenticated()
}

func (c *CacheLayer) GetCacheDir() string {
	return c.mediaClientRef.Get().GetCacheDir()
}

func (c *CacheLayer) CustomQuery(body []byte, logger *zerolog.Logger, token ...string) (interface{}, error) {
	// Use the stringified body as cache key
	cacheKey := string(body)
	bucket := c.buckets[CustomQueryBucket]

	// Try network first if API is working
	if IsWorking.Load() {
		res, err := c.mediaClientRef.Get().CustomQuery(body, logger, token...)
		c.checkAndUpdateWorkingState(err)

		if err == nil {
			go func() {
				if !ShouldCache.Load() {
					return
				}
				allData, err := filecache.GetAll[interface{}](c.fileCacher, filecache.NewBucket(bucket.Name(), 0))
				if err == nil && len(allData) >= maxNonCollectionCacheEntries {
					_ = c.fileCacher.DeletePermOldest(bucket)
				}

				if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
					c.logger.Warn().Err(err).Msg("simkl cache: Failed to cache custom query result")
				}
			}()
			return res, nil
		}
	} else {
		// If API is not working, try it in the background to check if it's back
		go func() {
			res, err := c.mediaClientRef.Get().CustomQuery(body, logger, token...)
			c.checkAndUpdateWorkingState(err)
			if err == nil {
				// Cache the result for future use with bounded size
				allData, err := filecache.GetAll[interface{}](c.fileCacher, filecache.NewBucket(bucket.Name(), 0))
				if err == nil && len(allData) >= maxNonCollectionCacheEntries {
					_ = c.fileCacher.DeletePermOldest(bucket)
				}

				if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
					c.logger.Warn().Err(err).Msg("simkl cache: Failed to cache background custom query result")
				}
			}
		}()
	}

	// Fall back to cache
	var cached interface{}
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &cached)
	if err != nil {
		return nil, fmt.Errorf("cache lookup failed: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("no cached data available")
	}

	c.logger.Debug().Str("bucket", CustomQueryBucket).Str("key", cacheKey).Msg("simkl cache: Serving custom query from cache")
	return cached, nil
}

// checkAndUpdateWorkingState checks if the API client is working and updates the state
func (c *CacheLayer) checkAndUpdateWorkingState(err error) {
	if err != nil {
		// Skip context.Canceled errors, not indicative of API issues
		if errors.Is(err, context.Canceled) {
			return
		}

		// skip 404 errors
		if strings.Contains(err.Error(), "404") {
			return
		}
		// skip 429 errors
		if strings.Contains(err.Error(), "429") {
			return
		}

		errStr := strings.ToLower(err.Error())
		// handle invalid token
		if strings.Contains(errStr, "user not found") {
			events.GlobalWSEventManager.SendEvent(events.ServerLoggedOutSimkl, "Your SIMKL session has expired. Please log in again.")
			if c.logoutFunc != nil {
				go c.logoutFunc()
			}
			return
		}

		// Add failure to tracking
		addFailureRecord(err)

		// Only mark as down if we have enough recent failures and are currently marked as working
		if IsWorking.Load() {
			recentFailures := getRecentFailureCount()
			if recentFailures >= failureThreshold {
				c.logger.Warn().
					Err(err).
					Int("recent_failures", recentFailures).
					Dur("within_window", failureWindow).
					Msg("simkl cache: Multiple API failures detected, switching to cache-only mode.")
				events.GlobalWSEventManager.SendEvent(events.WarningToast,
					fmt.Sprintf("The SIMKL API is experiencing issues (%d failures in %v), switching to cache-only mode.",
						recentFailures, failureWindow))
				IsWorking.Store(false)
			} else {
				c.logger.Debug().
					Err(err).
					Int("recent_failures", recentFailures).
					Int("threshold", failureThreshold).
					Msg("simkl cache: API failure recorded, monitoring for more failures")
			}
		}
	} else {
		// clear failure tracking and mark as working if not already
		if !IsWorking.Load() {
			c.logger.Info().Msg("simkl cache: API client is working again, switching back to network-first mode.")
			events.GlobalWSEventManager.SendEvent(events.InfoToast, "The SIMKL API is back online")
			IsWorking.Store(true)
		}
		clearFailureTracking()
	}
}

// generateCacheKey generates a cache key from the given parameters
func (c *CacheLayer) generateCacheKey(params ...interface{}) string {
	var keyParts []string
	for _, param := range params {
		if param == nil {
			keyParts = append(keyParts, "nil")
			continue
		}
		switch v := param.(type) {
		case *int:
			if v != nil {
				keyParts = append(keyParts, strconv.Itoa(*v))
			} else {
				keyParts = append(keyParts, "nil")
			}
		case *string:
			if v != nil {
				keyParts = append(keyParts, *v)
			} else {
				keyParts = append(keyParts, "nil")
			}
		case *bool:
			if v != nil {
				keyParts = append(keyParts, strconv.FormatBool(*v))
			} else {
				keyParts = append(keyParts, "nil")
			}
		case []*int:
			tmp := make([]int, 0, len(v))
			for _, id := range v {
				if id != nil {
					tmp = append(tmp, *id)
				}
			}
			slices.Sort(tmp)
			for _, id := range tmp {
				keyParts = append(keyParts, strconv.Itoa(id))
			}
		case []*string:
			tmp := make([]string, 0, len(v))
			for _, s := range v {
				if s != nil {
					tmp = append(tmp, *s)
				}
			}
			slices.Sort(tmp)
			keyParts = append(keyParts, tmp...)
		default:
			keyParts = append(keyParts, fmt.Sprintf("%v", param))
		}
	}
	return lo.Reduce(keyParts, func(acc, item string, _ int) string {
		if acc == "" {
			return item
		}
		return acc + "-" + item
	}, "")
}

// isInCollection checks if a media ID is in the user's collection
func (c *CacheLayer) isInCollection(mediaID int) bool {
	// Update collection tracking if needed
	c.updateCollectionTracking()
	_, ok := c.collectionMediaIDs.Get(mediaID)
	return ok
}

// updateCollectionTracking updates the collection media IDs tracking
func (c *CacheLayer) updateCollectionTracking() {
	if time.Since(c.lastCollectionUpdate) < collectionUpdateInterval {
		return
	}

	go func() {
		defer func() {
			c.lastCollectionUpdate = time.Now()
		}()

		// Try to fetch anime collection
		if animeCollection, err := c.mediaClientRef.Get().AnimeCollection(context.Background(), nil); err == nil && animeCollection != nil {
			for _, list := range animeCollection.MediaListCollection.Lists {
				if list != nil {
					for _, entry := range list.Entries {
						if entry != nil && entry.Media != nil {
							c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
						}
					}
				}
			}
		}

		// Try to fetch manga collection
		if mangaCollection, err := c.mediaClientRef.Get().MangaCollection(context.Background(), nil); err == nil && mangaCollection != nil {
			for _, list := range mangaCollection.MediaListCollection.Lists {
				if list != nil {
					for _, entry := range list.Entries {
						if entry != nil && entry.Media != nil {
							c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
						}
					}
				}
			}
		}
	}()
}

// networkFirstGet performs a network-first get operation with caching
func networkFirstGet[T any](c *CacheLayer, bucketName string, cacheKey string, networkFn func() (*T, error)) (*T, error) {
	if !ShouldCache.Load() {
		return networkFn()
	}

	bucket := c.buckets[bucketName]

	// Try network first if API is working
	if IsWorking.Load() {
		res, err := networkFn()
		c.checkAndUpdateWorkingState(err)

		if err == nil && res != nil {
			// Cache the successful result
			if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
				c.logger.Warn().Err(err).Msg("simkl cache: Failed to cache result")
			}
			return res, nil
		}
	} else {
		// If API is not working, try it in the background to check if it's back
		go func() {
			res, err := networkFn()
			c.checkAndUpdateWorkingState(err)
			if err == nil && res != nil {
				// Cache the result for future use
				if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
					c.logger.Warn().Err(err).Msg("simkl cache: Failed to cache background result")
				}
			}
		}()
	}

	// Fall back to cache
	var cached T
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &cached)
	if err != nil {
		return nil, fmt.Errorf("cache lookup failed: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("no cached data available")
	}

	c.logger.Debug().Str("bucket", bucketName).Str("key", cacheKey).Msg("simkl cache: Serving from cache")
	return &cached, nil
}

// boundedCacheSet caches data with a limit on non-collection entries
func (c *CacheLayer) boundedCacheSet(bucketName string, cacheKey string, data interface{}, mediaID int) error {
	if !ShouldCache.Load() {
		return nil
	}

	bucket := c.buckets[bucketName]

	// Always cache collection media
	if c.isInCollection(mediaID) {
		return c.fileCacher.SetPerm(bucket, cacheKey, data)
	}

	// For non-collection media, enforce the limit
	allData, err := filecache.GetAll[interface{}](c.fileCacher, filecache.NewBucket(bucket.Name(), 0))
	if err != nil {
		return err
	}

	// If we're at the limit, remove the oldest entry (simple FIFO for now)
	if len(allData) >= maxNonCollectionMediaCacheEntries {
		// Remove the first key we find (this is a simple implementation)
		for key := range allData {
			if err := c.fileCacher.DeletePerm(bucket, key); err == nil {
				break
			}
		}
	}

	return c.fileCacher.SetPerm(bucket, cacheKey, data)
}

// updateCollectionTrackingFromAnimeCollection updates collection tracking from anime collection
func (c *CacheLayer) updateCollectionTrackingFromAnimeCollection(collection *mediaapi.AnimeCollection) {
	if !ShouldCache.Load() {
		return
	}

	if !ShouldCache.Load() || collection == nil || collection.MediaListCollection == nil {
		return
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil {
			for _, entry := range list.Entries {
				if entry != nil && entry.Media != nil {
					c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
				}
			}
		}
	}
	c.lastCollectionUpdate = time.Now()
}

func (c *CacheLayer) updateCollectionTrackingFromAnimeCollectionWithRelations(collection *mediaapi.AnimeCollectionWithRelations) {
	if !ShouldCache.Load() {
		return
	}

	if !ShouldCache.Load() || collection == nil || collection.MediaListCollection == nil {
		return
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil {
			for _, entry := range list.Entries {
				if entry != nil && entry.Media != nil {
					c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
				}
			}
		}
	}
	c.lastCollectionUpdate = time.Now()
}

func (c *CacheLayer) updateCollectionTrackingFromMangaCollection(collection *mediaapi.MangaCollection) {
	if !ShouldCache.Load() {
		return
	}

	if !ShouldCache.Load() || collection == nil || collection.MediaListCollection == nil {
		return
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil {
			for _, entry := range list.Entries {
				if entry != nil && entry.Media != nil {
					c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
				}
			}
		}
	}
	c.lastCollectionUpdate = time.Now()
}

// invalidateMediaCaches invalidates caches for a specific media ID
func (c *CacheLayer) invalidateMediaCaches(mediaID int) {
	if !ShouldCache.Load() {
		return
	}

	mediaIDStr := strconv.Itoa(mediaID)

	// Delete from all media-specific buckets
	buckets := []string{
		BaseAnimeBucket,
		CompleteAnimeBucket,
		AnimeDetailsBucket,
		BaseMangaBucket,
		MangaDetailsBucket,
	}

	for _, bucketName := range buckets {
		bucket := c.buckets[bucketName]
		if err := c.fileCacher.DeletePerm(bucket, mediaIDStr); err != nil {
			c.logger.Debug().Err(err).Str("bucket", bucketName).Int("mediaID", mediaID).Msg("simkl cache: Failed to invalidate cache entry")
		}
	}
}

// invalidateCollectionCaches invalidates all collection caches and custom queries
func (c *CacheLayer) invalidateCollectionCaches() {
	if !ShouldCache.Load() {
		return
	}

	collectionBuckets := []string{
		AnimeCollectionBucket,
		AnimeCollectionTagsBucket,
		AnimeCollectionRelationsBucket,
		MangaCollectionBucket,
		MangaCollectionTagsBucket,
		CustomQueryBucket,
	}

	for _, bucketName := range collectionBuckets {
		bucket := c.buckets[bucketName]
		if err := c.fileCacher.EmptyPerm(bucket); err != nil {
			c.logger.Warn().Err(err).Str("bucket", bucketName).Msg("simkl cache: Failed to invalidate collection cache")
		}
	}

	// Reset collection tracking
	c.collectionMediaIDs.Clear()
	c.lastCollectionUpdate = time.Time{}
}

// extractBaseAnimeFromCollection attempts to extract BaseAnime data from cached anime collection
func (c *CacheLayer) extractBaseAnimeFromCollection(mediaID int) *mediaapi.BaseAnimeByID {
	// Try anime collection
	bucket := c.buckets[AnimeCollectionBucket]
	cacheKey := c.generateCacheKey("collection", nil)
	var animeCollection mediaapi.AnimeCollection
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &animeCollection)
	if err == nil && found && animeCollection.MediaListCollection != nil {
		for _, list := range animeCollection.MediaListCollection.Lists {
			if list != nil {
				for _, entry := range list.Entries {
					if entry != nil && entry.Media != nil && entry.Media.ID == mediaID {
						return &mediaapi.BaseAnimeByID{
							Media: entry.Media,
						}
					}
				}
			}
		}
	}

	// Try anime collection with relations
	relBucket := c.buckets[AnimeCollectionRelationsBucket]
	var animeCollectionRel mediaapi.AnimeCollectionWithRelations
	found, err = c.fileCacher.GetPerm(relBucket, cacheKey, &animeCollectionRel)
	if err == nil && found && animeCollectionRel.MediaListCollection != nil {
		for _, list := range animeCollectionRel.MediaListCollection.Lists {
			if list != nil {
				for _, entry := range list.Entries {
					if entry != nil && entry.Media != nil && entry.Media.ID == mediaID {
						return &mediaapi.BaseAnimeByID{
							Media: entry.Media.ToBaseAnime(),
						}
					}
				}
			}
		}
	}

	return nil
}

// extractBaseMangaFromCollection attempts to extract BaseManga data from cached manga collection
func (c *CacheLayer) extractBaseMangaFromCollection(mediaID int) *mediaapi.BaseMangaByID {
	if !ShouldCache.Load() {
		return nil
	}

	bucket := c.buckets[MangaCollectionBucket]
	cacheKey := c.generateCacheKey("collection", nil)
	var mangaCollection mediaapi.MangaCollection
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &mangaCollection)
	if err == nil && found && mangaCollection.MediaListCollection != nil {
		for _, list := range mangaCollection.MediaListCollection.Lists {
			if list != nil {
				for _, entry := range list.Entries {
					if entry != nil && entry.Media != nil && entry.Media.ID == mediaID {
						return &mediaapi.BaseMangaByID{
							Media: entry.Media,
						}
					}
				}
			}
		}
	}

	return nil
}

// networkFirstGetWithBoundedCache performs a network-first get operation with bounded caching for list/search results
func networkFirstGetWithBoundedCache[T any](c *CacheLayer, bucketName string, cacheKey string, networkFn func() (*T, error)) (*T, error) {
	bucket := c.buckets[bucketName]

	// Try network first if API is working
	if IsWorking.Load() {
		res, err := networkFn()
		c.checkAndUpdateWorkingState(err)

		if err == nil && res != nil {
			// Cache the successful result with bounded size
			go func() {
				// For list/search results, always apply bounded caching
				allData, err := filecache.GetAll[interface{}](c.fileCacher, filecache.NewBucket(bucket.Name(), 0))
				if err == nil && len(allData) >= maxNonCollectionCacheEntries {
					_ = c.fileCacher.DeletePermOldest(bucket)
				}

				if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
					c.logger.Warn().Err(err).Msg("simkl cache: Failed to cache bounded result")
				}
			}()
			return res, nil
		}
	} else {
		// If API is not working, try it in the background to check if it's back
		go func() {
			res, err := networkFn()
			c.checkAndUpdateWorkingState(err)
			if err == nil && res != nil {
				// Cache the result for future use with bounded size
				allData, err := filecache.GetAll[interface{}](c.fileCacher, filecache.NewBucket(bucket.Name(), 0))
				if err == nil && len(allData) >= maxNonCollectionCacheEntries {
					_ = c.fileCacher.DeletePermOldest(bucket)
				}

				if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
					c.logger.Warn().Err(err).Msg("simkl cache: Failed to cache background bounded result")
				}
			}
		}()
	}

	// Fall back to cache
	var cached T
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &cached)
	if err != nil {
		return nil, fmt.Errorf("cache lookup failed: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("no cached data available")
	}

	c.logger.Debug().Str("bucket", bucketName).Str("key", cacheKey).Msg("simkl cache: Serving bounded result from cache")
	return &cached, nil
}

func (c *CacheLayer) AnimeCollection(ctx context.Context, userName *string, interceptors ...clientv2.RequestInterceptor) (*mediaapi.AnimeCollection, error) {
	cacheKey := c.generateCacheKey("collection", nil)
	res, err := networkFirstGet(c, AnimeCollectionBucket, cacheKey, func() (*mediaapi.AnimeCollection, error) {
		return c.mediaClientRef.Get().AnimeCollection(ctx, userName, interceptors...)
	})

	if err == nil && res != nil && c.applyQueuedUpdatesToAnimeCollection(res) {
		if err := c.fileCacher.SetPerm(c.buckets[AnimeCollectionBucket], cacheKey, res); err != nil {
			c.logger.Warn().Err(err).Msg("simkl cache: Failed to apply queued updates to anime collection cache")
		}
	}

	// Update collection tracking with the fetched data
	if err == nil && res != nil {
		go c.updateCollectionTrackingFromAnimeCollection(res)
	}

	return res, err
}

func (c *CacheLayer) AnimeCollectionTags(ctx context.Context, userName *string, interceptors ...clientv2.RequestInterceptor) (*mediaapi.AnimeCollectionTags, error) {
	cacheKey := c.generateCacheKey("collection-tags", userName)
	return networkFirstGet(c, AnimeCollectionTagsBucket, cacheKey, func() (*mediaapi.AnimeCollectionTags, error) {
		return c.mediaClientRef.Get().AnimeCollectionTags(ctx, userName, interceptors...)
	})
}

func (c *CacheLayer) AnimeCollectionWithRelations(ctx context.Context, userName *string, interceptors ...clientv2.RequestInterceptor) (*mediaapi.AnimeCollectionWithRelations, error) {
	cacheKey := c.generateCacheKey("collection-relations", nil)
	res, err := networkFirstGet(c, AnimeCollectionRelationsBucket, cacheKey, func() (*mediaapi.AnimeCollectionWithRelations, error) {
		return c.mediaClientRef.Get().AnimeCollectionWithRelations(ctx, userName, interceptors...)
	})

	// Update collection tracking with the fetched data
	if err == nil && res != nil {
		go c.updateCollectionTrackingFromAnimeCollectionWithRelations(res)
	}

	return res, err
}

func (c *CacheLayer) BaseAnimeByMalID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.BaseAnimeByMalID, error) {
	if id == nil {
		return c.mediaClientRef.Get().BaseAnimeByMalID(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey("mal", id)
	return networkFirstGet(c, BaseAnimeMalBucket, cacheKey, func() (*mediaapi.BaseAnimeByMalID, error) {
		return c.mediaClientRef.Get().BaseAnimeByMalID(ctx, id, interceptors...)
	})
}

func (c *CacheLayer) BaseAnimeByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.BaseAnimeByID, error) {
	if id == nil {
		return c.mediaClientRef.Get().BaseAnimeByID(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, BaseAnimeBucket, cacheKey, func() (*mediaapi.BaseAnimeByID, error) {
		return c.mediaClientRef.Get().BaseAnimeByID(ctx, id, interceptors...)
	})

	// If network and direct cache failed, try to extract from collection cache
	if err != nil {
		if collectionResult := c.extractBaseAnimeFromCollection(*id); collectionResult != nil {
			c.logger.Debug().Int("mediaID", *id).Msg("simkl cache: Extracted BaseAnime from collection cache")
			return collectionResult, nil
		}
	}

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(BaseAnimeBucket, cacheKey, res, *id); err != nil {
				c.logger.Warn().Err(err).Msg("simkl cache: Failed to update bounded cache")
			}
		}()
	}

	return res, err
}

func (c *CacheLayer) SearchBaseAnimeByIds(ctx context.Context, ids []*int, page *int, perPage *int, status []*mediaapi.MediaStatus, inCollection *bool, sort []*mediaapi.MediaSort, season *mediaapi.MediaSeason, year *int, genre *string, format *mediaapi.MediaFormat, interceptors ...clientv2.RequestInterceptor) (*mediaapi.SearchBaseAnimeByIds, error) {
	cacheKey := c.generateCacheKey(ids, page, perPage, status, inCollection, sort, season, year, genre, format)
	return networkFirstGetWithBoundedCache(c, SearchBaseAnimeByIdsBucket, cacheKey, func() (*mediaapi.SearchBaseAnimeByIds, error) {
		return c.mediaClientRef.Get().SearchBaseAnimeByIds(ctx, ids, page, perPage, status, inCollection, sort, season, year, genre, format, interceptors...)
	})
}

func (c *CacheLayer) CompleteAnimeByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.CompleteAnimeByID, error) {
	if id == nil {
		return c.mediaClientRef.Get().CompleteAnimeByID(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, CompleteAnimeBucket, cacheKey, func() (*mediaapi.CompleteAnimeByID, error) {
		return c.mediaClientRef.Get().CompleteAnimeByID(ctx, id, interceptors...)
	})

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(CompleteAnimeBucket, cacheKey, res, *id); err != nil {
				c.logger.Warn().Err(err).Msg("simkl cache: failed to update bounded cache")
			}
		}()
	}

	return res, err
}

func (c *CacheLayer) AnimeDetailsByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.AnimeDetailsByID, error) {
	if id == nil {
		return c.mediaClientRef.Get().AnimeDetailsByID(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, AnimeDetailsBucket, cacheKey, func() (*mediaapi.AnimeDetailsByID, error) {
		return c.mediaClientRef.Get().AnimeDetailsByID(ctx, id, interceptors...)
	})

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(AnimeDetailsBucket, cacheKey, res, *id); err != nil {
				c.logger.Warn().Err(err).Msg("simkl cache: failed to update bounded cache")
			}
		}()
	}

	return res, err
}

func (c *CacheLayer) ListAnime(ctx context.Context, page *int, search *string, perPage *int, sort []*mediaapi.MediaSort, status []*mediaapi.MediaStatus, genres []*string, tags []*string, averageScoreGreater *int, season *mediaapi.MediaSeason, seasonYear *int, format *mediaapi.MediaFormat, isAdult *bool, interceptors ...clientv2.RequestInterceptor) (*mediaapi.ListAnime, error) {
	cacheKey := c.generateCacheKey(page, search, perPage, sort, status, genres, averageScoreGreater, season, seasonYear, format, isAdult)
	return networkFirstGetWithBoundedCache(c, ListMediaBucket, cacheKey, func() (*mediaapi.ListAnime, error) {
		return c.mediaClientRef.Get().ListAnime(ctx, page, search, perPage, sort, status, genres, tags, averageScoreGreater, season, seasonYear, format, isAdult, interceptors...)
	})
}

func (c *CacheLayer) ListRecentAnime(ctx context.Context, page *int, perPage *int, airingAtGreater *int, airingAtLesser *int, notYetAired *bool, interceptors ...clientv2.RequestInterceptor) (*mediaapi.ListRecentAnime, error) {
	// devnote: don't include airingAt params since they're unique for each requests, just return from the other params
	cacheKey := c.generateCacheKey(page, perPage, notYetAired)
	return networkFirstGetWithBoundedCache(c, ListRecentAiringMediaBucket, cacheKey, func() (*mediaapi.ListRecentAnime, error) {
		return c.mediaClientRef.Get().ListRecentAnime(ctx, page, perPage, airingAtGreater, airingAtLesser, notYetAired, interceptors...)
	})
}

func (c *CacheLayer) UpdateMediaListEntry(ctx context.Context, mediaID *int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput, interceptors ...clientv2.RequestInterceptor) (*mediaapi.UpdateMediaListEntry, error) {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		entryID, err := c.queueMediaListEntryUpdate(mediaID, status, scoreRaw, progress, startedAt, completedAt)
		if err != nil {
			return nil, err
		}
		return &mediaapi.UpdateMediaListEntry{SaveMediaListEntry: &mediaapi.UpdateMediaListEntry_SaveMediaListEntry{ID: entryID}}, nil
	}

	res, err := c.sendMediaListEntryUpdate(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt, interceptors...)
	c.checkAndUpdateWorkingState(err)
	if err != nil && shouldQueueMediaListUpdate(err) {
		entryID, queueErr := c.queueMediaListEntryUpdate(mediaID, status, scoreRaw, progress, startedAt, completedAt)
		if queueErr != nil {
			return nil, queueErr
		}
		return &mediaapi.UpdateMediaListEntry{SaveMediaListEntry: &mediaapi.UpdateMediaListEntry_SaveMediaListEntry{ID: entryID}}, nil
	}

	// Invalidate relevant caches on successful mutation
	if err == nil && mediaID != nil {
		c.invalidateMediaCaches(*mediaID)
		c.invalidateCollectionCaches()
	}

	return res, err
}

func (c *CacheLayer) UpdateMediaListEntryProgress(ctx context.Context, mediaID *int, progress *int, status *mediaapi.MediaListStatus, interceptors ...clientv2.RequestInterceptor) (*mediaapi.UpdateMediaListEntryProgress, error) {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		entryID, err := c.queueMediaListEntryProgressUpdate(mediaID, progress, status)
		if err != nil {
			return nil, err
		}
		return &mediaapi.UpdateMediaListEntryProgress{SaveMediaListEntry: &mediaapi.UpdateMediaListEntryProgress_SaveMediaListEntry{ID: entryID}}, nil
	}

	res, err := c.sendMediaListEntryProgressUpdate(ctx, mediaID, progress, status, interceptors...)
	c.checkAndUpdateWorkingState(err)
	if err != nil && shouldQueueMediaListUpdate(err) {
		entryID, queueErr := c.queueMediaListEntryProgressUpdate(mediaID, progress, status)
		if queueErr != nil {
			return nil, queueErr
		}
		return &mediaapi.UpdateMediaListEntryProgress{SaveMediaListEntry: &mediaapi.UpdateMediaListEntryProgress_SaveMediaListEntry{ID: entryID}}, nil
	}

	// Invalidate relevant caches on successful mutation
	if err == nil && mediaID != nil {
		c.invalidateMediaCaches(*mediaID)
		c.invalidateCollectionCaches()
	}

	return res, err
}

func (c *CacheLayer) UpdateMediaListEntryRepeat(ctx context.Context, mediaID *int, repeat *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.UpdateMediaListEntryRepeat, error) {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		return nil, fmt.Errorf("simkl cache: API client is not working, mutation operations are not available")
	}

	res, err := c.mediaClientRef.Get().UpdateMediaListEntryRepeat(ctx, mediaID, repeat, interceptors...)
	c.checkAndUpdateWorkingState(err)

	// Invalidate relevant caches on successful mutation
	if err == nil && mediaID != nil {
		c.invalidateMediaCaches(*mediaID)
		c.invalidateCollectionCaches()
	}

	return res, err
}

func (c *CacheLayer) DeleteEntry(ctx context.Context, mediaListEntryID *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.DeleteEntry, error) {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		return nil, fmt.Errorf("simkl cache: API client is not working, mutation operations are not available")
	}

	res, err := c.mediaClientRef.Get().DeleteEntry(ctx, mediaListEntryID, interceptors...)
	c.checkAndUpdateWorkingState(err)

	// Invalidate collection caches on successful deletion
	if err == nil {
		c.invalidateCollectionCaches()
	}

	return res, err
}

func (c *CacheLayer) MangaCollection(ctx context.Context, userName *string, interceptors ...clientv2.RequestInterceptor) (*mediaapi.MangaCollection, error) {
	cacheKey := c.generateCacheKey("collection", nil)
	res, err := networkFirstGet(c, MangaCollectionBucket, cacheKey, func() (*mediaapi.MangaCollection, error) {
		return c.mediaClientRef.Get().MangaCollection(ctx, userName, interceptors...)
	})

	if err == nil && res != nil && c.applyQueuedUpdatesToMangaCollection(res) {
		if err := c.fileCacher.SetPerm(c.buckets[MangaCollectionBucket], cacheKey, res); err != nil {
			c.logger.Warn().Err(err).Msg("simkl cache: Failed to apply queued updates to manga collection cache")
		}
	}

	// Update collection tracking with the fetched data
	if err == nil && res != nil {
		go c.updateCollectionTrackingFromMangaCollection(res)
	}

	return res, err
}

func (c *CacheLayer) MangaCollectionTags(ctx context.Context, userName *string, interceptors ...clientv2.RequestInterceptor) (*mediaapi.MangaCollectionTags, error) {
	cacheKey := c.generateCacheKey("collection-tags", userName)
	return networkFirstGet(c, MangaCollectionTagsBucket, cacheKey, func() (*mediaapi.MangaCollectionTags, error) {
		return c.mediaClientRef.Get().MangaCollectionTags(ctx, userName, interceptors...)
	})
}

func (c *CacheLayer) SearchBaseManga(ctx context.Context, page *int, perPage *int, sort []*mediaapi.MediaSort, search *string, status []*mediaapi.MediaStatus, interceptors ...clientv2.RequestInterceptor) (*mediaapi.SearchBaseManga, error) {
	cacheKey := c.generateCacheKey(page, perPage, sort, search, status)
	return networkFirstGetWithBoundedCache(c, SearchBaseMangaBucket, cacheKey, func() (*mediaapi.SearchBaseManga, error) {
		return c.mediaClientRef.Get().SearchBaseManga(ctx, page, perPage, sort, search, status, interceptors...)
	})
}

func (c *CacheLayer) BaseMangaByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.BaseMangaByID, error) {
	if id == nil {
		return c.mediaClientRef.Get().BaseMangaByID(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, BaseMangaBucket, cacheKey, func() (*mediaapi.BaseMangaByID, error) {
		return c.mediaClientRef.Get().BaseMangaByID(ctx, id, interceptors...)
	})

	// If network and direct cache failed, try to extract from collection cache
	if err != nil {
		if collectionResult := c.extractBaseMangaFromCollection(*id); collectionResult != nil {
			c.logger.Debug().Int("mediaID", *id).Msg("simkl cache: Extracted BaseManga from collection cache")
			return collectionResult, nil
		}
	}

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(BaseMangaBucket, cacheKey, res, *id); err != nil {
				c.logger.Warn().Err(err).Msg("simkl cache: Failed to update bounded cache")
			}
		}()
	}

	return res, err
}

func (c *CacheLayer) MangaDetailsByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.MangaDetailsByID, error) {
	if id == nil {
		return c.mediaClientRef.Get().MangaDetailsByID(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, MangaDetailsBucket, cacheKey, func() (*mediaapi.MangaDetailsByID, error) {
		return c.mediaClientRef.Get().MangaDetailsByID(ctx, id, interceptors...)
	})

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(MangaDetailsBucket, cacheKey, res, *id); err != nil {
				c.logger.Warn().Err(err).Msg("simkl cache: failed to update bounded cache")
			}
		}()
	}

	return res, err
}

func (c *CacheLayer) ListManga(ctx context.Context, page *int, search *string, perPage *int, sort []*mediaapi.MediaSort, status []*mediaapi.MediaStatus, genres []*string, tags []*string, averageScoreGreater *int, startDateGreater *string, startDateLesser *string, format *mediaapi.MediaFormat, countryOfOrigin *string, isAdult *bool, interceptors ...clientv2.RequestInterceptor) (*mediaapi.ListManga, error) {
	cacheKey := c.generateCacheKey(page, search, perPage, sort, status, genres, tags, averageScoreGreater, startDateGreater, startDateLesser, format, countryOfOrigin, isAdult)
	return networkFirstGetWithBoundedCache(c, ListMangaBucket, cacheKey, func() (*mediaapi.ListManga, error) {
		return c.mediaClientRef.Get().ListManga(ctx, page, search, perPage, sort, status, genres, tags, averageScoreGreater, startDateGreater, startDateLesser, format, countryOfOrigin, isAdult, interceptors...)
	})
}

func (c *CacheLayer) ViewerStats(ctx context.Context, interceptors ...clientv2.RequestInterceptor) (*mediaapi.ViewerStats, error) {
	cacheKey := "stats"
	return networkFirstGet(c, ViewerStatsBucket, cacheKey, func() (*mediaapi.ViewerStats, error) {
		return c.mediaClientRef.Get().ViewerStats(ctx, interceptors...)
	})
}

func (c *CacheLayer) StudioDetails(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.StudioDetails, error) {
	if id == nil {
		return c.mediaClientRef.Get().StudioDetails(ctx, id, interceptors...)
	}

	cacheKey := c.generateCacheKey(id)
	return networkFirstGet(c, StudioDetailsBucket, cacheKey, func() (*mediaapi.StudioDetails, error) {
		return c.mediaClientRef.Get().StudioDetails(ctx, id, interceptors...)
	})
}

func (c *CacheLayer) GetViewer(ctx context.Context, interceptors ...clientv2.RequestInterceptor) (*mediaapi.GetViewer, error) {
	cacheKey := "viewer"
	return networkFirstGet(c, ViewerBucket, cacheKey, func() (*mediaapi.GetViewer, error) {
		return c.mediaClientRef.Get().GetViewer(ctx, interceptors...)
	})
}

func (c *CacheLayer) AnimeAiringSchedule(ctx context.Context, ids []*int, season *mediaapi.MediaSeason, seasonYear *int, previousSeason *mediaapi.MediaSeason, previousSeasonYear *int, nextSeason *mediaapi.MediaSeason, nextSeasonYear *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.AnimeAiringSchedule, error) {
	cacheKey := c.generateCacheKey(ids, season, seasonYear, previousSeason, previousSeasonYear, nextSeason, nextSeasonYear)
	return networkFirstGet(c, AnimeAiringScheduleBucket, cacheKey, func() (*mediaapi.AnimeAiringSchedule, error) {
		return c.mediaClientRef.Get().AnimeAiringSchedule(ctx, ids, season, seasonYear, previousSeason, previousSeasonYear, nextSeason, nextSeasonYear, interceptors...)
	})
}

func (c *CacheLayer) AnimeAiringScheduleRaw(ctx context.Context, ids []*int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.AnimeAiringScheduleRaw, error) {
	cacheKey := c.generateCacheKey(ids)
	return networkFirstGet(c, AnimeAiringScheduleRawBucket, cacheKey, func() (*mediaapi.AnimeAiringScheduleRaw, error) {
		return c.mediaClientRef.Get().AnimeAiringScheduleRaw(ctx, ids, interceptors...)
	})
}
