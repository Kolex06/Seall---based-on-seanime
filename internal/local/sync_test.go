package local

import (
	"errors"
	"fmt"
	"seall/internal/api/mediaapi"
	"seall/internal/extension"
	"seall/internal/platforms/media_platform"
	"seall/internal/platforms/platform"
	"seall/internal/testmocks"
	"seall/internal/testutil"
	"seall/internal/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testSetupManager(t *testing.T) (Manager, *mediaapi.AnimeCollection, *mediaapi.MangaCollection) {
	env := testutil.NewTestEnv(t)
	logger := env.Logger()

	database := env.MustNewDatabase(logger)
	simklClient := mediaapi.NewTestMediaApiClient()
	extensionBankRef := util.NewRef(extension.NewUnifiedBank())
	mediaPlatform := media_platform.NewMediaPlatform(util.NewRef[mediaapi.MediaApiClient](simklClient), extensionBankRef, logger, database)
	animeCollection, err := mediaPlatform.GetMediaCollection(t.Context(), true)
	require.NoError(t, err)
	mangaCollection, err := mediaPlatform.GetMangaCollection(t.Context(), true)
	require.NoError(t, err)

	manager := NewTestManager(t, database)

	manager.SetAnimeCollection(animeCollection)
	manager.SetMangaCollection(mangaCollection)

	return manager, animeCollection, mangaCollection
}

func TestSync2(t *testing.T) {
	manager, animeCollection, _ := testSetupManager(t)

	err := manager.TrackAnime(130003) // Bocchi the rock
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}
	err = manager.TrackAnime(10800) // Chihayafuru
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}
	err = manager.TrackAnime(171457) // Make Heroine ga Oosugiru!
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}
	err = manager.TrackManga(101517) // JJK
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}

	err = manager.SynchronizeLocal()
	require.NoError(t, err)

	select {
	case <-manager.GetSyncer().doneUpdatingLocalCollections:
		util.Spew(manager.GetLocalAnimeCollection().MustGet())
		util.Spew(manager.GetLocalMangaCollection().MustGet())
		break
	case <-time.After(10 * time.Second):
		t.Log("Timeout")
		break
	}

	mediaapi.PatchAnimeCollectionEntry(animeCollection, 130003, mediaapi.AnimeCollectionEntryPatch{
		Status:   new(mediaapi.MediaListStatusCompleted),
		Progress: new(12), // Mock progress
	})

	fmt.Println("================================================================================================")
	fmt.Println("================================================================================================")

	err = manager.SynchronizeLocal()
	require.NoError(t, err)

	select {
	case <-manager.GetSyncer().doneUpdatingLocalCollections:
		util.Spew(manager.GetLocalAnimeCollection().MustGet())
		util.Spew(manager.GetLocalMangaCollection().MustGet())
		break
	case <-time.After(10 * time.Second):
		t.Log("Timeout")
		break
	}

}

func TestSync(t *testing.T) {
	manager, _, _ := testSetupManager(t)

	err := manager.TrackAnime(130003) // Bocchi the rock
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}
	err = manager.TrackAnime(10800) // Chihayafuru
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}
	err = manager.TrackAnime(171457) // Make Heroine ga Oosugiru!
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}
	err = manager.TrackManga(101517) // JJK
	if err != nil && !errors.Is(err, ErrAlreadyTracked) {
		require.NoError(t, err)
	}

	err = manager.SynchronizeLocal()
	require.NoError(t, err)

	select {
	case <-manager.GetSyncer().doneUpdatingLocalCollections:
		util.Spew(manager.GetLocalAnimeCollection().MustGet())
		util.Spew(manager.GetLocalMangaCollection().MustGet())
		break
	case <-time.After(10 * time.Second):
		t.Log("Timeout")
		break
	}

}

func TestSynchronizeSimklDoesNotPanicWithoutLocalCollections(t *testing.T) {
	manager, _, _ := testSetupManager(t)

	require.NotPanics(t, func() {
		require.NoError(t, manager.SynchronizeMediaApi())
	})
}

func TestSynchronizeSimulatedCollectionToSimklCreatesMissingEntries(t *testing.T) {
	manager, animeCollection, mangaCollection := testSetupManager(t)
	managerImpl := manager.(*ManagerImpl)

	animeEntry, found := animeCollection.GetListEntryFromAnimeId(130003)
	require.True(t, found)
	mangaEntry, found := mangaCollection.GetListEntryFromMangaId(101517)
	require.True(t, found)

	manager.SaveSimulatedAnimeCollection(newSingleAnimeCollection(animeEntry))
	manager.SaveSimulatedMangaCollection(newSingleMangaCollection(mangaEntry))

	manager.SetAnimeCollection(newEmptyAnimeCollection())
	manager.SetMangaCollection(newEmptyMangaCollection())

	fakePlatform := testmocks.NewFakePlatformBuilder().Build()
	managerImpl.mediaPlatformRef = util.NewRef[platform.Platform](fakePlatform)

	require.NoError(t, manager.SynchronizeSimulatedCollectionToMediaApi())

	updateCalls := fakePlatform.UpdateEntryCalls()
	require.Len(t, updateCalls, 2)
	require.Contains(t, []int{updateCalls[0].MediaID, updateCalls[1].MediaID}, 130003)
	require.Contains(t, []int{updateCalls[0].MediaID, updateCalls[1].MediaID}, 101517)

	for _, call := range updateCalls {
		switch call.MediaID {
		case 130003:
			require.NotNil(t, call.Status)
			require.Equal(t, animeEntry.GetStatus(), call.Status)
		case 101517:
			require.NotNil(t, call.Status)
			require.Equal(t, mangaEntry.GetStatus(), call.Status)
		}
	}
}

func newEmptyAnimeCollection() *mediaapi.AnimeCollection {
	return &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{},
		},
	}
}

func newEmptyMangaCollection() *mediaapi.MangaCollection {
	return &mediaapi.MangaCollection{
		MediaListCollection: &mediaapi.MangaCollection_MediaListCollection{
			Lists: []*mediaapi.MangaCollection_MediaListCollection_Lists{},
		},
	}
}

func newSingleAnimeCollection(entry *mediaapi.AnimeListEntry) *mediaapi.AnimeCollection {
	return &mediaapi.AnimeCollection{
		MediaListCollection: &mediaapi.AnimeCollection_MediaListCollection{
			Lists: []*mediaapi.AnimeCollection_MediaListCollection_Lists{
				{
					Status:  entry.Status,
					Entries: []*mediaapi.AnimeCollection_MediaListCollection_Lists_Entries{entry},
				},
			},
		},
	}
}

func newSingleMangaCollection(entry *mediaapi.MangaListEntry) *mediaapi.MangaCollection {
	return &mediaapi.MangaCollection{
		MediaListCollection: &mediaapi.MangaCollection_MediaListCollection{
			Lists: []*mediaapi.MangaCollection_MediaListCollection_Lists{
				{
					Status:  entry.Status,
					Entries: []*mediaapi.MangaCollection_MediaListCollection_Lists_Entries{entry},
				},
			},
		},
	}
}
