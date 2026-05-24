package simulated_platform

import (
	"context"
	"errors"
	"seall/internal/api/mediaapi"
	"seall/internal/extension"
	"seall/internal/local"
	"seall/internal/testmocks"
	"seall/internal/testutil"
	"seall/internal/util"
	"testing"

	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/stretchr/testify/require"
)

func TestRefreshAnimeCollectionRefreshesMutableEntries(t *testing.T) {
	sp, manager, client := newTestSimulatedPlatform(t, newRefreshingFixtureClient(
		map[int]*mediaapi.BaseAnime{
			101: testmocks.NewBaseAnimeBuilder(101, "anime current fresh").WithStatus(mediaapi.MediaStatusFinished).Build(),
			102: testmocks.NewBaseAnimeBuilder(102, "anime paused fresh").WithStatus(mediaapi.MediaStatusReleasing).Build(),
			103: testmocks.NewBaseAnimeBuilder(103, "anime planning fresh").WithStatus(mediaapi.MediaStatusFinished).Build(),
		},
		nil,
	))

	// keep a mix of mutable and settled entries so refresh only fetches the ones that can change
	manager.SaveSimulatedAnimeCollection(&mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{
				newAnimeCollectionList(mediaapi.MediaListStatusCurrent,
					newAnimeCollectionEntry(testmocks.NewBaseAnimeBuilder(101, "anime current stale").WithStatus(mediaapi.MediaStatusReleasing).Build(), mediaapi.MediaListStatusCurrent),
					newAnimeCollectionEntry(testmocks.NewBaseAnimeBuilder(105, "anime settled stale").WithStatus(mediaapi.MediaStatusFinished).Build(), mediaapi.MediaListStatusCurrent),
				),
				newAnimeCollectionList(mediaapi.MediaListStatusPaused,
					newAnimeCollectionEntry(testmocks.NewBaseAnimeBuilder(102, "anime paused stale").WithStatus(mediaapi.MediaStatusNotYetReleased).Build(), mediaapi.MediaListStatusPaused),
				),
				newAnimeCollectionList(mediaapi.MediaListStatusPlanning,
					newAnimeCollectionEntry(testmocks.NewBaseAnimeBuilder(103, "anime planning stale").WithStatus(mediaapi.MediaStatusReleasing).Build(), mediaapi.MediaListStatusPlanning),
				),
				newAnimeCollectionList(mediaapi.MediaListStatusCompleted,
					newAnimeCollectionEntry(testmocks.NewBaseAnimeBuilder(104, "anime completed stale").WithStatus(mediaapi.MediaStatusReleasing).Build(), mediaapi.MediaListStatusCompleted),
				),
			},
		},
	})

	_, err := sp.RefreshAnimeCollection(context.Background())
	require.NoError(t, err)

	require.ElementsMatch(t, []int{101, 102, 103}, client.animeCalls)

	collection := manager.GetSimulatedAnimeCollection().MustGet()
	currentEntry, found := collection.GetListEntryFromAnimeId(101)
	require.True(t, found)
	require.Equal(t, "anime current fresh", *currentEntry.GetMedia().GetTitle().GetEnglish())
	require.Equal(t, mediaapi.MediaStatusFinished, *currentEntry.GetMedia().GetStatus())

	pausedEntry, found := collection.GetListEntryFromAnimeId(102)
	require.True(t, found)
	require.Equal(t, "anime paused fresh", *pausedEntry.GetMedia().GetTitle().GetEnglish())

	planningEntry, found := collection.GetListEntryFromAnimeId(103)
	require.True(t, found)
	require.Equal(t, "anime planning fresh", *planningEntry.GetMedia().GetTitle().GetEnglish())

	completedEntry, found := collection.GetListEntryFromAnimeId(104)
	require.True(t, found)
	require.Equal(t, "anime completed stale", *completedEntry.GetMedia().GetTitle().GetEnglish())

	settledEntry, found := collection.GetListEntryFromAnimeId(105)
	require.True(t, found)
	require.Equal(t, "anime settled stale", *settledEntry.GetMedia().GetTitle().GetEnglish())
}

func TestRefreshMangaCollectionRefreshesMutableEntries(t *testing.T) {
	sp, manager, client := newTestSimulatedPlatform(t, newRefreshingFixtureClient(
		nil,
		map[int]*mediaapi.BaseManga{
			201: testmocks.NewBaseMangaBuilder(201, "manga current fresh").WithStatus(mediaapi.MediaStatusFinished).Build(),
			202: testmocks.NewBaseMangaBuilder(202, "manga paused fresh").WithStatus(mediaapi.MediaStatusReleasing).Build(),
			203: testmocks.NewBaseMangaBuilder(203, "manga planning fresh").WithStatus(mediaapi.MediaStatusFinished).Build(),
		},
	))

	// refresh should skip dropped or already settled manga entries.
	manager.SaveSimulatedMangaCollection(&mediaapi.MangaCollection{
		MediaListCollection: &mediaapi.MangaCollection_MediaListCollection{
			Lists: []*mediaapi.MangaCollection_MediaListCollection_Lists{
				newMangaCollectionList(mediaapi.MediaListStatusCurrent,
					newMangaCollectionEntry(testmocks.NewBaseMangaBuilder(201, "manga current stale").WithStatus(mediaapi.MediaStatusReleasing).Build(), mediaapi.MediaListStatusCurrent),
					newMangaCollectionEntry(testmocks.NewBaseMangaBuilder(205, "manga settled stale").WithStatus(mediaapi.MediaStatusFinished).Build(), mediaapi.MediaListStatusCurrent),
				),
				newMangaCollectionList(mediaapi.MediaListStatusPaused,
					newMangaCollectionEntry(testmocks.NewBaseMangaBuilder(202, "manga paused stale").WithStatus(mediaapi.MediaStatusNotYetReleased).Build(), mediaapi.MediaListStatusPaused),
				),
				newMangaCollectionList(mediaapi.MediaListStatusPlanning,
					newMangaCollectionEntry(testmocks.NewBaseMangaBuilder(203, "manga planning stale").WithStatus(mediaapi.MediaStatusReleasing).Build(), mediaapi.MediaListStatusPlanning),
				),
				newMangaCollectionList(mediaapi.MediaListStatusDropped,
					newMangaCollectionEntry(testmocks.NewBaseMangaBuilder(204, "manga dropped stale").WithStatus(mediaapi.MediaStatusReleasing).Build(), mediaapi.MediaListStatusDropped),
				),
			},
		},
	})

	_, err := sp.RefreshMangaCollection(context.Background())
	require.NoError(t, err)

	require.ElementsMatch(t, []int{201, 202, 203}, client.mangaCalls)

	collection := manager.GetSimulatedMangaCollection().MustGet()
	currentEntry, found := collection.GetListEntryFromMangaId(201)
	require.True(t, found)
	require.Equal(t, "manga current fresh", *currentEntry.GetMedia().GetTitle().GetEnglish())
	require.Equal(t, mediaapi.MediaStatusFinished, *currentEntry.GetMedia().GetStatus())

	pausedEntry, found := collection.GetListEntryFromMangaId(202)
	require.True(t, found)
	require.Equal(t, "manga paused fresh", *pausedEntry.GetMedia().GetTitle().GetEnglish())

	planningEntry, found := collection.GetListEntryFromMangaId(203)
	require.True(t, found)
	require.Equal(t, "manga planning fresh", *planningEntry.GetMedia().GetTitle().GetEnglish())

	droppedEntry, found := collection.GetListEntryFromMangaId(204)
	require.True(t, found)
	require.Equal(t, "manga dropped stale", *droppedEntry.GetMedia().GetTitle().GetEnglish())

	settledEntry, found := collection.GetListEntryFromMangaId(205)
	require.True(t, found)
	require.Equal(t, "manga settled stale", *settledEntry.GetMedia().GetTitle().GetEnglish())
}

func newTestSimulatedPlatform(t *testing.T, client *refreshingFixtureClient) (*SimulatedPlatform, local.Manager, *refreshingFixtureClient) {
	t.Helper()

	env := testutil.NewTestEnv(t)
	logger := env.Logger()
	database := env.MustNewDatabase(logger)
	manager := local.NewTestManager(t, database)

	platformInstance, err := NewSimulatedPlatform(
		manager,
		util.NewRef[mediaapi.MediaApiClient](client),
		util.NewRef(extension.NewUnifiedBank()),
		logger,
		database,
	)
	require.NoError(t, err)

	sp, ok := platformInstance.(*SimulatedPlatform)
	require.True(t, ok)
	return sp, manager, client
}

type refreshingFixtureClient struct {
	mediaapi.MediaApiClient
	animeByID  map[int]*mediaapi.BaseAnime
	mangaByID  map[int]*mediaapi.BaseManga
	animeCalls []int
	mangaCalls []int
}

func newRefreshingFixtureClient(animeByID map[int]*mediaapi.BaseAnime, mangaByID map[int]*mediaapi.BaseManga) *refreshingFixtureClient {
	return &refreshingFixtureClient{
		MediaApiClient: mediaapi.NewTestMediaApiClient(),
		animeByID:      animeByID,
		mangaByID:      mangaByID,
	}
}

func (c *refreshingFixtureClient) BaseAnimeByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.BaseAnimeByID, error) {
	if id != nil {
		c.animeCalls = append(c.animeCalls, *id)
		if media, ok := c.animeByID[*id]; ok {
			return &mediaapi.BaseAnimeByID{Media: media}, nil
		}
	}

	return nil, errors.New("unexpected anime refresh request")
}

func (c *refreshingFixtureClient) BaseMangaByID(ctx context.Context, id *int, interceptors ...clientv2.RequestInterceptor) (*mediaapi.BaseMangaByID, error) {
	if id != nil {
		c.mangaCalls = append(c.mangaCalls, *id)
		if media, ok := c.mangaByID[*id]; ok {
			return &mediaapi.BaseMangaByID{Media: media}, nil
		}
	}

	return nil, errors.New("unexpected manga refresh request")
}

func newAnimeCollectionList(status mediaapi.MediaListStatus, entries ...*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries) *mediaapi.AnimeCollection_MediaListCollection_Lists {
	return &mediaapi.AnimeCollection_MediaListCollection_Lists{
		Status:       &status,
		Name:         new(string(status)),
		IsCustomList: new(false),
		Entries:      entries,
	}
}

func newAnimeCollectionEntry(media *mediaapi.BaseAnime, status mediaapi.MediaListStatus) *mediaapi.AnimeCollection_MediaListCollection_Lists_Entries {
	return &mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{
		Media:    media,
		Progress: new(0),
		Score:    new(0.0),
		Repeat:   new(0),
		Status:   &status,
	}
}

func newMangaCollectionList(status mediaapi.MediaListStatus, entries ...*mediaapi.MangaCollection_MediaListCollection_Lists_Entries) *mediaapi.MangaCollection_MediaListCollection_Lists {
	return &mediaapi.MangaCollection_MediaListCollection_Lists{
		Status:       &status,
		Name:         new(string(status)),
		IsCustomList: new(false),
		Entries:      entries,
	}
}

func newMangaCollectionEntry(media *mediaapi.BaseManga, status mediaapi.MediaListStatus) *mediaapi.MangaCollection_MediaListCollection_Lists_Entries {
	return &mediaapi.MangaCollection_MediaListCollection_Lists_Entries{
		Media:    media,
		Progress: new(0),
		Score:    new(0.0),
		Repeat:   new(0),
		Status:   &status,
	}
}
