package local

import (
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata_provider"
	"seall/internal/database/db"
	"seall/internal/database/db_bridge"
	"seall/internal/database/models"
	"seall/internal/events"
	"seall/internal/extension"
	"seall/internal/library/anime"
	"seall/internal/manga"
	"seall/internal/platforms/media_platform"
	"seall/internal/platforms/platform"
	"seall/internal/testutil"
	"seall/internal/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTestManager(t *testing.T, db *db.Database) Manager {
	env := testutil.NewTestEnv(t)

	logger := env.Logger()
	metadataProvider := metadata_provider.NewTestProviderWithEnv(env, db)
	metadataProviderRef := util.NewRef(metadataProvider)
	mangaRepository := manga.NewTestRepositoryWithEnv(env, db)

	wsEventManager := events.NewMockWSEventManager(logger)
	simklClient := mediaapi.NewFixtureMediaApiClient()
	mediaClientRef := util.NewRef[mediaapi.MediaApiClient](simklClient)
	extensionBankRef := util.NewRef(extension.NewUnifiedBank())
	mediaPlatform := media_platform.NewMediaPlatform(mediaClientRef, extensionBankRef, logger, db)
	mediaPlatformRef := util.NewRef[platform.Platform](mediaPlatform)

	localDir := env.MustMkdirData("offline")
	assetsDir := env.MustMkdirData("offline", "assets")

	var localFilesCount int64
	err := db.Gorm().Model(&models.LocalFiles{}).Count(&localFilesCount).Error
	require.NoError(t, err)
	if localFilesCount == 0 {
		_, err = db_bridge.InsertLocalFiles(db, make([]*anime.LocalFile, 0))
		require.NoError(t, err)
	}

	m, err := NewManager(&NewManagerOptions{
		LocalDir:            localDir,
		AssetDir:            assetsDir,
		Logger:              logger,
		MetadataProviderRef: metadataProviderRef,
		MangaRepository:     mangaRepository,
		Database:            db,
		WSEventManager:      wsEventManager,
		MediaPlatformRef:    mediaPlatformRef,
		IsOffline:           false,
	})
	require.NoError(t, err)

	return m
}
