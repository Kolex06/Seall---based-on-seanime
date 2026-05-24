package scanner

import (
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata_provider"
	"seall/internal/database/db"
	"seall/internal/events"
	"seall/internal/extension"
	"seall/internal/library/anime"
	"seall/internal/platforms/media_platform"
	"seall/internal/testutil"
	"seall/internal/util"
	"seall/internal/util/limiter"
	"testing"

	"github.com/rs/zerolog"
)

const scannerTestLibraryDir = "E:/Anime"

type scannerTestWrapper struct {
	Env                 *testutil.TestEnv
	Config              *testutil.Config
	Logger              *zerolog.Logger
	Database            *db.Database
	MediaApiClient      mediaapi.MediaApiClient
	Platform            *media_platform.MediaPlatform
	MetadataProvider    metadata_provider.Provider
	CompleteAnimeCache  *mediaapi.CompleteAnimeCache
	MediaApiRateLimiter *limiter.Limiter
	WSEventManager      events.WSEventManagerInterface
	LibraryDir          string
}

func newScannerFixtureWrapper(t testing.TB) *scannerTestWrapper {
	t.Helper()

	env := testutil.NewTestEnv(t)
	return newScannerWrapper(t, env, mediaapi.NewTestMediaApiClient(), "")
}

func newScannerLiveWrapper(t testing.TB) *scannerTestWrapper {
	t.Helper()

	env := testutil.NewTestEnv(t, testutil.MediaApi())
	cfg := env.Config()

	return newScannerWrapper(t, env, mediaapi.NewMediaApiClient(cfg.Provider.MediaApiJwt, ""), cfg.Provider.MediaApiUsername)
}

func newScannerWrapper(t testing.TB, env *testutil.TestEnv, client mediaapi.MediaApiClient, username string) *scannerTestWrapper {
	t.Helper()

	logger := env.Logger()
	database := env.MustNewDatabase(logger)
	mediaClientRef := util.NewRef(client)
	extensionBankRef := util.NewRef(extension.NewUnifiedBank())
	platform := media_platform.NewMediaPlatform(mediaClientRef, extensionBankRef, logger, database).(*media_platform.MediaPlatform)
	if username != "" {
		platform.SetUsername(username)
	}

	return &scannerTestWrapper{
		Env:                 env,
		Config:              env.Config(),
		Logger:              logger,
		Database:            database,
		MediaApiClient:      client,
		Platform:            platform,
		MetadataProvider:    metadata_provider.NewTestProviderWithEnv(env, database),
		CompleteAnimeCache:  mediaapi.NewCompleteAnimeCache(),
		MediaApiRateLimiter: limiter.NewMediaApiLimiter(),
		WSEventManager:      events.NewMockWSEventManager(logger),
		LibraryDir:          scannerTestLibraryDir,
	}
}

func (h *scannerTestWrapper) LocalFiles(paths ...string) []*anime.LocalFile {
	localFiles := make([]*anime.LocalFile, 0, len(paths))
	for _, path := range paths {
		localFiles = append(localFiles, anime.NewLocalFile(path, h.LibraryDir))
	}

	return localFiles
}
