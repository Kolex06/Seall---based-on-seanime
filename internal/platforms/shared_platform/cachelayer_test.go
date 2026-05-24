package shared_platform

import (
	"context"
	"seall/internal/api/mediaapi"
	"seall/internal/util"
	"strconv"
	"testing"

	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/stretchr/testify/require"
)

type cacheLayerTestClient struct {
	mediaapi.MediaApiClient
	cacheDir            string
	animeCollection     *mediaapi.AnimeCollection
	mangaCollection     *mediaapi.MangaCollection
	updateEntryCalls    []cacheLayerUpdateEntryCall
	updateProgressCalls []cacheLayerUpdateProgressCall
}

type cacheLayerUpdateEntryCall struct {
	MediaID     *int
	Status      *mediaapi.MediaListStatus
	ScoreRaw    *int
	Progress    *int
	StartedAt   *mediaapi.FuzzyDateInput
	CompletedAt *mediaapi.FuzzyDateInput
}

type cacheLayerUpdateProgressCall struct {
	MediaID  *int
	Progress *int
	Status   *mediaapi.MediaListStatus
}

func (c *cacheLayerTestClient) IsAuthenticated() bool {
	return true
}

func (c *cacheLayerTestClient) GetCacheDir() string {
	return c.cacheDir
}

func (c *cacheLayerTestClient) BaseAnimeByID(_ context.Context, id *int, _ ...clientv2.RequestInterceptor) (*mediaapi.BaseAnimeByID, error) {
	mediaID := 0
	if id != nil {
		mediaID = *id
	}
	return &mediaapi.BaseAnimeByID{Media: &mediaapi.BaseAnime{ID: mediaID}}, nil
}

func (c *cacheLayerTestClient) AnimeCollection(_ context.Context, _ *string, _ ...clientv2.RequestInterceptor) (*mediaapi.AnimeCollection, error) {
	return c.animeCollection, nil
}

func (c *cacheLayerTestClient) MangaCollection(_ context.Context, _ *string, _ ...clientv2.RequestInterceptor) (*mediaapi.MangaCollection, error) {
	return c.mangaCollection, nil
}

func (c *cacheLayerTestClient) UpdateMediaListEntry(_ context.Context, mediaID *int, status *mediaapi.MediaListStatus, scoreRaw *int, progress *int, startedAt *mediaapi.FuzzyDateInput, completedAt *mediaapi.FuzzyDateInput, _ ...clientv2.RequestInterceptor) (*mediaapi.UpdateMediaListEntry, error) {
	c.updateEntryCalls = append(c.updateEntryCalls, cacheLayerUpdateEntryCall{
		MediaID:     newCloned(mediaID),
		Status:      newCloned(status),
		ScoreRaw:    newCloned(scoreRaw),
		Progress:    newCloned(progress),
		StartedAt:   cloneFuzzyDateInput(startedAt),
		CompletedAt: cloneFuzzyDateInput(completedAt),
	})
	return &mediaapi.UpdateMediaListEntry{SaveMediaListEntry: &mediaapi.UpdateMediaListEntry_SaveMediaListEntry{ID: 999}}, nil
}

func (c *cacheLayerTestClient) UpdateMediaListEntryProgress(_ context.Context, mediaID *int, progress *int, status *mediaapi.MediaListStatus, _ ...clientv2.RequestInterceptor) (*mediaapi.UpdateMediaListEntryProgress, error) {
	c.updateProgressCalls = append(c.updateProgressCalls, cacheLayerUpdateProgressCall{
		MediaID:  newCloned(mediaID),
		Progress: newCloned(progress),
		Status:   newCloned(status),
	})
	return &mediaapi.UpdateMediaListEntryProgress{SaveMediaListEntry: &mediaapi.UpdateMediaListEntryProgress_SaveMediaListEntry{ID: 999}}, nil
}

func TestCacheLayerQueuesProgressUpdateAndPatchesAnimeCache(t *testing.T) {
	client := &cacheLayerTestClient{
		cacheDir:        t.TempDir(),
		animeCollection: newTestAnimeCollection(101, 321, mediaapi.MediaListStatusCurrent, 2),
	}
	cacheLayer := newTestCacheLayer(t, client)

	// get the collection in cache
	_, err := cacheLayer.AnimeCollection(context.Background(), new("user"))
	require.NoError(t, err)

	IsWorking.Store(false)
	res, err := cacheLayer.UpdateMediaListEntryProgress(context.Background(), new(101), new(6), new(mediaapi.MediaListStatusCompleted))
	require.NoError(t, err)
	require.Equal(t, 321, res.GetSaveMediaListEntry().GetID())
	require.Empty(t, client.updateProgressCalls)

	// the cached entry should move lists immediately so refetches see the local change.
	cached := getCachedAnimeCollection(t, cacheLayer)
	entry, found := cached.GetListEntryFromAnimeId(101)
	require.True(t, found)
	require.Equal(t, 6, *entry.GetProgress())
	require.Equal(t, mediaapi.MediaListStatusCompleted, *entry.GetStatus())
	require.True(t, animeListContains(cached, mediaapi.MediaListStatusCompleted, 101))
	require.False(t, animeListContains(cached, mediaapi.MediaListStatusCurrent, 101))

	queued := getQueuedUpdate(t, cacheLayer, 101)
	require.Equal(t, 101, queued.MediaID)
	require.Equal(t, 6, *queued.Progress)
	require.Equal(t, mediaapi.MediaListStatusCompleted, *queued.Status)
	require.False(t, queued.FullUpdate)
}

func TestCacheLayerQueuesEntryUpdateAndSyncsWhenOnline(t *testing.T) {
	client := &cacheLayerTestClient{
		cacheDir:        t.TempDir(),
		mangaCollection: newTestMangaCollection(202, 654, mediaapi.MediaListStatusCurrent, 4),
	}
	cacheLayer := newTestCacheLayer(t, client)

	// seed manga cache, then queue an edit while the api is marked down
	_, err := cacheLayer.MangaCollection(context.Background(), new("user"))
	require.NoError(t, err)

	IsWorking.Store(false)
	startedAt := &mediaapi.FuzzyDateInput{Year: new(2025), Month: new(1), Day: new(2)}
	completedAt := &mediaapi.FuzzyDateInput{Year: new(2025), Month: new(2), Day: new(3)}
	res, err := cacheLayer.UpdateMediaListEntry(context.Background(), new(202), new(mediaapi.MediaListStatusCompleted), new(85), new(12), startedAt, completedAt)
	require.NoError(t, err)
	require.Equal(t, 654, res.GetSaveMediaListEntry().GetID())
	require.Empty(t, client.updateEntryCalls)

	cached := getCachedMangaCollection(t, cacheLayer)
	entry, found := cached.GetListEntryFromMangaId(202)
	require.True(t, found)
	require.Equal(t, 12, *entry.GetProgress())
	require.Equal(t, float64(85), *entry.GetScore())
	require.Equal(t, mediaapi.MediaListStatusCompleted, *entry.GetStatus())
	require.True(t, mangaListContains(cached, mediaapi.MediaListStatusCompleted, 202))
	require.False(t, mangaListContains(cached, mediaapi.MediaListStatusCurrent, 202))

	// when the api is healthy again, the queued full edit is flushed once and removed
	IsWorking.Store(true)
	cacheLayer.syncQueuedUpdates(context.Background())
	require.Len(t, client.updateEntryCalls, 1)
	require.Equal(t, 202, *client.updateEntryCalls[0].MediaID)
	require.Equal(t, mediaapi.MediaListStatusCompleted, *client.updateEntryCalls[0].Status)
	require.Equal(t, 85, *client.updateEntryCalls[0].ScoreRaw)
	require.Equal(t, 12, *client.updateEntryCalls[0].Progress)
	require.Equal(t, 2025, *client.updateEntryCalls[0].StartedAt.Year)
	require.Equal(t, 3, *client.updateEntryCalls[0].CompletedAt.Day)
	requireNoQueuedUpdate(t, cacheLayer, 202)
}

func TestCacheLayerLiveProgressUpdateClearsQueuedUpdate(t *testing.T) {
	client := &cacheLayerTestClient{
		cacheDir:        t.TempDir(),
		animeCollection: newTestAnimeCollection(101, 321, mediaapi.MediaListStatusCurrent, 2),
	}
	cacheLayer := newTestCacheLayer(t, client)

	_, err := cacheLayer.AnimeCollection(context.Background(), new("user"))
	require.NoError(t, err)

	// first update is queued while the api is marked down
	IsWorking.Store(false)
	_, err = cacheLayer.UpdateMediaListEntryProgress(context.Background(), new(101), new(6), new(mediaapi.MediaListStatusCompleted))
	require.NoError(t, err)
	queued := getQueuedUpdate(t, cacheLayer, 101)
	require.Equal(t, 6, *queued.Progress)

	// a later online update should win and remove the stale queued state
	IsWorking.Store(true)
	_, err = cacheLayer.UpdateMediaListEntryProgress(context.Background(), new(101), new(7), new(mediaapi.MediaListStatusCurrent))
	require.NoError(t, err)
	requireNoQueuedUpdate(t, cacheLayer, 101)

	cacheLayer.syncQueuedUpdates(context.Background())
	require.Len(t, client.updateProgressCalls, 1)
	require.Equal(t, 7, *client.updateProgressCalls[0].Progress)
}

func TestCacheLayerLiveEntryUpdateClearsQueuedUpdate(t *testing.T) {
	client := &cacheLayerTestClient{
		cacheDir:        t.TempDir(),
		mangaCollection: newTestMangaCollection(202, 654, mediaapi.MediaListStatusCurrent, 4),
	}
	cacheLayer := newTestCacheLayer(t, client)

	_, err := cacheLayer.MangaCollection(context.Background(), new("user"))
	require.NoError(t, err)

	// queue an older edit while SIMKL is unavailable
	IsWorking.Store(false)
	_, err = cacheLayer.UpdateMediaListEntry(context.Background(), new(202), new(mediaapi.MediaListStatusCompleted), new(80), new(12), nil, nil)
	require.NoError(t, err)
	queued := getQueuedUpdate(t, cacheLayer, 202)
	require.Equal(t, 80, *queued.ScoreRaw)

	// the successful online edit replaces it and should prevent stale replay
	IsWorking.Store(true)
	_, err = cacheLayer.UpdateMediaListEntry(context.Background(), new(202), new(mediaapi.MediaListStatusCurrent), new(90), new(13), nil, nil)
	require.NoError(t, err)
	requireNoQueuedUpdate(t, cacheLayer, 202)

	cacheLayer.syncQueuedUpdates(context.Background())
	require.Len(t, client.updateEntryCalls, 1)
	require.Equal(t, 90, *client.updateEntryCalls[0].ScoreRaw)
	require.Equal(t, 13, *client.updateEntryCalls[0].Progress)
}

func newTestCacheLayer(t *testing.T, client *cacheLayerTestClient) *CacheLayer {
	t.Helper()
	ShouldCache.Store(true)
	IsWorking.Store(true)
	clearFailureTracking()

	clientRef := util.NewRef[mediaapi.MediaApiClient](client)
	cacheLayer, ok := NewCacheLayer(clientRef).(*CacheLayer)
	require.True(t, ok)

	t.Cleanup(func() {
		ShouldCache.Store(true)
		IsWorking.Store(true)
		clearFailureTracking()
	})

	return cacheLayer
}

func newTestAnimeCollection(mediaID int, entryID int, status mediaapi.MediaListStatus, progress int) *mediaapi.AnimeCollection {
	return &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{
				newTestAnimeList(status, newTestAnimeEntry(mediaID, entryID, status, progress)),
				newTestAnimeList(mediaapi.MediaListStatusCompleted),
			},
		},
	}
}

func newTestAnimeList(status mediaapi.MediaListStatus, entries ...*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries) *mediaapi.AnimeCollection_MediaListCollection_Lists {
	return &mediaapi.AnimeCollection_MediaListCollection_Lists{
		Status:       new(status),
		Name:         new(string(status)),
		IsCustomList: new(false),
		Entries:      entries,
	}
}

func newTestAnimeEntry(mediaID int, entryID int, status mediaapi.MediaListStatus, progress int) *mediaapi.AnimeCollection_MediaListCollection_Lists_Entries {
	return &mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{
		ID:       entryID,
		Media:    &mediaapi.BaseAnime{ID: mediaID, Episodes: new(12)},
		Status:   new(status),
		Progress: new(progress),
		Score:    new(0.0),
	}
}

func newTestMangaCollection(mediaID int, entryID int, status mediaapi.MediaListStatus, progress int) *mediaapi.MangaCollection {
	return &mediaapi.MangaCollection{
		MediaListCollection: &mediaapi.MangaCollection_MediaListCollection{
			Lists: []*mediaapi.MangaCollection_MediaListCollection_Lists{
				newTestMangaList(status, newTestMangaEntry(mediaID, entryID, status, progress)),
				newTestMangaList(mediaapi.MediaListStatusCompleted),
			},
		},
	}
}

func newTestMangaList(status mediaapi.MediaListStatus, entries ...*mediaapi.MangaCollection_MediaListCollection_Lists_Entries) *mediaapi.MangaCollection_MediaListCollection_Lists {
	return &mediaapi.MangaCollection_MediaListCollection_Lists{
		Status:       new(status),
		Name:         new(string(status)),
		IsCustomList: new(false),
		Entries:      entries,
	}
}

func newTestMangaEntry(mediaID int, entryID int, status mediaapi.MediaListStatus, progress int) *mediaapi.MangaCollection_MediaListCollection_Lists_Entries {
	return &mediaapi.MangaCollection_MediaListCollection_Lists_Entries{
		ID:       entryID,
		Media:    &mediaapi.BaseManga{ID: mediaID, Chapters: new(20)},
		Status:   new(status),
		Progress: new(progress),
		Score:    new(0.0),
	}
}

func getCachedAnimeCollection(t *testing.T, cacheLayer *CacheLayer) *mediaapi.AnimeCollection {
	t.Helper()
	var cached mediaapi.AnimeCollection
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[AnimeCollectionBucket], cacheLayer.generateCacheKey("collection", nil), &cached)
	require.NoError(t, err)
	require.True(t, found)
	return &cached
}

func getCachedMangaCollection(t *testing.T, cacheLayer *CacheLayer) *mediaapi.MangaCollection {
	t.Helper()
	var cached mediaapi.MangaCollection
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[MangaCollectionBucket], cacheLayer.generateCacheKey("collection", nil), &cached)
	require.NoError(t, err)
	require.True(t, found)
	return &cached
}

func getQueuedUpdate(t *testing.T, cacheLayer *CacheLayer, mediaID int) queuedMediaListUpdate {
	t.Helper()
	var queued queuedMediaListUpdate
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[PendingMediaListUpdatesBucket], strconv.Itoa(mediaID), &queued)
	require.NoError(t, err)
	require.True(t, found)
	return queued
}

func requireNoQueuedUpdate(t *testing.T, cacheLayer *CacheLayer, mediaID int) {
	t.Helper()
	var queued queuedMediaListUpdate
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[PendingMediaListUpdatesBucket], strconv.Itoa(mediaID), &queued)
	require.NoError(t, err)
	require.False(t, found)
}

func animeListContains(collection *mediaapi.AnimeCollection, status mediaapi.MediaListStatus, mediaID int) bool {
	for _, list := range collection.GetMediaListCollection().GetLists() {
		if list.GetStatus() == nil || *list.GetStatus() != status {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry.GetMedia().GetID() == mediaID {
				return true
			}
		}
	}
	return false
}

func mangaListContains(collection *mediaapi.MangaCollection, status mediaapi.MediaListStatus, mediaID int) bool {
	for _, list := range collection.GetMediaListCollection().GetLists() {
		if list.GetStatus() == nil || *list.GetStatus() != status {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry.GetMedia().GetID() == mediaID {
				return true
			}
		}
	}
	return false
}
