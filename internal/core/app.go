package core

import (
	"os"
	"path/filepath"
	"runtime"
	"seall/internal/api/mediaapi"
	"seall/internal/api/metadata_provider"
	"seall/internal/api/simkl"
	"seall/internal/constants"
	"seall/internal/continuity"
	"seall/internal/database/db"
	"seall/internal/database/models"
	debrid_client "seall/internal/debrid/client"
	"seall/internal/directstream"
	discordrpc_presence "seall/internal/discordrpc/presence"
	"seall/internal/doh"
	"seall/internal/events"
	"seall/internal/extension"
	"seall/internal/extension_playground"
	"seall/internal/extension_repo"
	"seall/internal/hook"
	"seall/internal/library/autodownloader"
	"seall/internal/library/autoscanner"
	"seall/internal/library/fillermanager"
	"seall/internal/library/playbackmanager"
	"seall/internal/library/scanner"
	"seall/internal/library_explorer"
	"seall/internal/local"
	"seall/internal/manga"
	"seall/internal/mediaplayers/iina"
	"seall/internal/mediaplayers/mediaplayer"
	"seall/internal/mediaplayers/mpchc"
	"seall/internal/mediaplayers/mpv"
	"seall/internal/mediaplayers/vlc"
	"seall/internal/mediastream"
	"seall/internal/nakama"
	"seall/internal/nativeplayer"
	"seall/internal/onlinestream"
	"seall/internal/platforms/offline_platform"
	"seall/internal/platforms/platform"
	"seall/internal/platforms/simkl_platform"
	"seall/internal/platforms/simulated_platform"
	"seall/internal/playlist"
	"seall/internal/plugin"
	"seall/internal/report"
	"seall/internal/torrent_clients/torrent_client"
	"seall/internal/torrents/torrent"
	"seall/internal/torrentstream"
	"seall/internal/updater"
	"seall/internal/user"
	"seall/internal/util"
	"seall/internal/util/filecache"
	"seall/internal/util/result"
	"seall/internal/videocore"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

type (
	App struct {
		// Core
		Config   *Config
		Database *db.Database
		Logger   *zerolog.Logger

		// Torrent and debrid services
		TorrentClientRepository *torrent_client.Repository
		TorrentRepository       *torrent.Repository
		DebridClientRepository  *debrid_client.Repository

		// File system monitoring
		Watcher *scanner.Watcher

		// API clients and providers
		MediaClientRef      *util.Ref[mediaapi.MediaApiClient]
		SimklClientRef      *util.Ref[*simkl.Client]
		MediaPlatformRef    *util.Ref[platform.Platform]
		OfflinePlatformRef  *util.Ref[platform.Platform]
		MetadataProviderRef *util.Ref[metadata_provider.Provider]

		// Library
		FillerManager   *fillermanager.FillerManager
		AutoDownloader  *autodownloader.AutoDownloader
		AutoScanner     *autoscanner.AutoScanner
		PlaybackManager *playbackmanager.PlaybackManager

		// Real-time communication
		WSEventManager *events.WSEventManager

		// Extensions
		ExtensionRepository           *extension_repo.Repository
		ExtensionBankRef              *util.Ref[*extension.UnifiedBank]
		ExtensionPlaygroundRepository *extension_playground.PlaygroundRepository

		// Streaming
		DirectStreamManager     *directstream.Manager
		OnlinestreamRepository  *onlinestream.Repository
		MediastreamRepository   *mediastream.Repository
		TorrentstreamRepository *torrentstream.Repository

		// Players
		NativePlayer *nativeplayer.NativePlayer
		VideoCore    *videocore.VideoCore
		MediaPlayer  struct {
			VLC   *vlc.VLC
			MpcHc *mpchc.MpcHc
			Mpv   *mpv.Mpv
			Iina  *iina.Iina
		}
		MediaPlayerRepository *mediaplayer.Repository

		// Manga services
		MangaRepository *manga.Repository
		MangaDownloader *manga.Downloader

		// Offline and local account
		LocalManager local.Manager

		// Utilities
		FileCacher       *filecache.Cacher
		Updater          *updater.Updater
		SelfUpdater      *updater.SelfUpdater
		ReportRepository *report.Repository

		// Integrations
		DiscordPresence *discordrpc_presence.Presence

		// Continuity and sync
		ContinuityManager *continuity.Manager

		// Lifecycle management
		Cleanups                      []func()
		OnRefreshMediaCollectionFuncs *result.Map[string, func()]
		OnFlushLogs                   func()

		// Configuration and feature flags
		FeatureFlags      FeatureFlags
		FeatureManager    *FeatureManager
		Settings          *models.Settings
		SecondarySettings struct {
			Mediastream   *models.MediastreamSettings
			Torrentstream *models.TorrentstreamSettings
			Debrid        *models.DebridSettings
		}

		// Metadata
		Version          string
		TotalLibrarySize uint64
		LibraryDir       string
		MediaApiCacheDir string
		SimklCacheDir    string
		IsDesktopSidecar bool
		Flags            SeallFlags

		// Internal state
		user                 *user.User
		previousVersion      string
		moduleMu             sync.Mutex
		ServerReady          bool
		isOfflineRef         *util.Ref[bool]
		ServerPasswordHash   string
		ClientIdentitySecret string
		logoutInProgress     atomic.Bool

		// Plugin system
		HookManager hook.Manager

		// Features
		PlaylistManager *playlist.Manager
		LibraryExplorer *library_explorer.LibraryExplorer
		NakamaManager   *nakama.Manager

		// Show this version's tour on the frontend
		// Hydrated by migrations.go when there's a version change
		ShowTour string
	}
)

// NewApp creates a new server instance
func NewApp(configOpts *ConfigOptions, selfupdater *updater.SelfUpdater) *App {

	var app *App

	// Initialize logger with predefined format
	logger := util.NewLogger()

	// Log application version, OS, architecture and system info
	logger.Info().Msgf("app: %s %s-%s", constants.AppName, constants.Version, constants.VersionName)
	logger.Info().Msgf("app: OS: %s", runtime.GOOS)
	logger.Info().Msgf("app: Arch: %s", runtime.GOARCH)
	logger.Info().Msgf("app: Processor count: %d", runtime.NumCPU())

	// Initialize hook manager for plugin event system
	hookManager := hook.NewHookManager(hook.NewHookManagerOptions{Logger: logger})
	hook.SetGlobalHookManager(hookManager)
	plugin.GlobalAppContext.SetLogger(logger)

	// Store current version to detect version changes
	previousVersion := constants.Version

	// Add callback to track version changes
	configOpts.OnVersionChange = append(configOpts.OnVersionChange, func(oldVersion string, newVersion string) {
		logger.Info().Str("prev", oldVersion).Str("current", newVersion).Msg("app: Version change detected")
		previousVersion = oldVersion
	})

	// Initialize configuration with provided options
	// Creates config directory if it doesn't exist
	cfg, err := NewConfig(configOpts, logger)
	if err != nil {
		logger.Fatal().Err(err).Msgf("app: Failed to initialize config")
	}

	// Compute SHA-256 hash of the server password
	serverPasswordHash := ""
	if cfg.Server.Password != "" {
		serverPasswordHash = util.HashSHA256Hex(cfg.Server.Password)
	}

	// Create logs directory if it doesn't exist
	_ = os.MkdirAll(cfg.Logs.Dir, 0755)

	// Start background process to trim log files
	go TrimLogEntries(cfg.Logs.Dir, logger)

	logger.Info().Msgf("app: Data directory: %s", cfg.Data.AppDataDir)
	logger.Info().Msgf("app: Working directory: %s", cfg.Data.WorkingDir)

	// Log if running in desktop sidecar mode
	if configOpts.Flags.IsDesktopSidecar {
		logger.Info().Msg("app: Desktop sidecar mode enabled")
	}

	// Initialize database connection
	database, err := db.NewDatabase(cfg.Data.AppDataDir, cfg.Database.Name, logger)
	if err != nil {
		logger.Fatal().Err(err).Msgf("app: Failed to initialize database")
	}

	HandleNewDatabaseEntries(database, logger)

	// Clean up old database entries using the cleanup manager to prevent concurrent access issues
	database.RunDatabaseCleanup() // Remove old entries from all tables sequentially

	// Get anime library paths for plugin context
	animeLibraryPaths, _ := database.GetAllLibraryPathsFromSettings()
	plugin.GlobalAppContext.SetModulesPartial(plugin.AppContextModules{
		Database:          database,
		AnimeLibraryPaths: &animeLibraryPaths,
	})

	simklToken := database.GetSimklToken()

	mediaApiCacheDir := filepath.Join(cfg.Cache.Dir, "simkl")
	simklCacheDir := filepath.Join(cfg.Cache.Dir, "simkl")

	mediaClient := mediaapi.NewMediaApiClient("", mediaApiCacheDir)
	mediaClientRef := util.NewRef[mediaapi.MediaApiClient](mediaClient)

	simklClient := simkl.NewClient(simkl.NewClientOptions{
		Token:        simklToken,
		ClientID:     cfg.Simkl.ClientID,
		ClientSecret: cfg.Simkl.ClientSecret,
		RedirectURI:  cfg.Simkl.RedirectURI,
		CacheDir:     simklCacheDir,
		Logger:       logger,
	})
	simklClientRef := util.NewRef(simklClient)

	// Initialize WebSocket event manager for real-time communication
	wsEventManager := events.NewWSEventManager(logger)

	// Exit if no WebSocket connections in desktop sidecar mode
	if configOpts.Flags.IsDesktopSidecar {
		wsEventManager.ExitIfNoConnsAsDesktopSidecar()
	}

	// Initialize DNS-over-HTTPS service in background
	go doh.HandleDoH(cfg.Server.DoHUrl, logger)

	// Initialize file cache system for media and metadata
	fileCacher, err := filecache.NewCacher(cfg.Cache.Dir)
	if err != nil {
		logger.Fatal().Err(err).Msgf("app: Failed to initialize file cacher")
	}

	// Initialize the extension bank that will be shared across modules
	extensionBankRef := util.NewRef(extension.NewUnifiedBank())

	// Initialize extension repository
	extensionRepository := extension_repo.NewRepository(&extension_repo.NewRepositoryOptions{
		Logger:           logger,
		ExtensionDir:     cfg.Extensions.Dir,
		WSEventManager:   wsEventManager,
		FileCacher:       fileCacher,
		HookManager:      hookManager,
		ExtensionBankRef: extensionBankRef,
	})

	// Initialize metadata provider for media information
	metadataProvider := metadata_provider.NewProvider(&metadata_provider.NewProviderImplOptions{
		Logger:           logger,
		FileCacher:       fileCacher,
		Database:         database,
		ExtensionBankRef: extensionBankRef,
		SimklClientRef:   simklClientRef,
	})

	// Set initial metadata provider (will change if offline mode is enabled)
	activeMetadataProvider := metadataProvider

	// Initialize manga repository
	mangaRepository := manga.NewRepository(&manga.NewRepositoryOptions{
		Logger:           logger,
		FileCacher:       fileCacher,
		CacheDir:         cfg.Cache.Dir,
		ServerURI:        cfg.GetServerURI(),
		WsEventManager:   wsEventManager,
		DownloadDir:      cfg.Manga.DownloadDir,
		Database:         database,
		ExtensionBankRef: extensionBankRef,
	})

	// SIMKL is the primary online platform for Seall.
	activePlatform := simkl_platform.NewSimklPlatform(simklClientRef, mediaClientRef, extensionBankRef, logger, database)
	activePlatformRef := util.NewRef[platform.Platform](activePlatform)
	metadataProviderRef := util.NewRef[metadata_provider.Provider](activeMetadataProvider)

	// Initialize sync manager for offline/online synchronization
	localManager, err := local.NewManager(&local.NewManagerOptions{
		LocalDir:            cfg.Offline.Dir,
		AssetDir:            cfg.Offline.AssetDir,
		Logger:              logger,
		MetadataProviderRef: metadataProviderRef,
		MangaRepository:     mangaRepository,
		Database:            database,
		WSEventManager:      wsEventManager,
		IsOffline:           cfg.Server.Offline,
		MediaPlatformRef:    activePlatformRef,
	})
	if err != nil {
		logger.Fatal().Err(err).Msgf("app: Failed to initialize sync manager")
	}

	// Use local metadata provider if in offline mode
	if cfg.Server.Offline {
		activeMetadataProvider = localManager.GetOfflineMetadataProvider()
	}

	// Initialize local platform for offline operations
	offlinePlatform, err := offline_platform.NewOfflinePlatform(localManager, mediaClientRef, logger)
	if err != nil {
		logger.Fatal().Err(err).Msgf("app: Failed to initialize local platform")
	}

	// Initialize simulated platform for unauthenticated operations
	simulatedPlatform, err := simulated_platform.NewSimulatedPlatform(localManager, mediaClientRef, extensionBankRef, logger, database)
	if err != nil {
		logger.Fatal().Err(err).Msgf("app: Failed to initialize simulated platform")
	}

	// Change active platform if offline mode is enabled
	if cfg.Server.Offline {
		logger.Warn().Msg("app: Offline mode is active, using offline platform")
		activePlatformRef.Set(offlinePlatform)
	} else if !simklClientRef.Get().IsAuthenticated() {
		logger.Warn().Msg("app: SIMKL client is not authenticated, using simulated platform")
		activePlatformRef.Set(simulatedPlatform)
	}

	isOfflineRef := util.NewRef(cfg.Server.Offline)
	offlinePlatformRef := util.NewRef[platform.Platform](offlinePlatform)

	// Update plugin context with new modules
	plugin.GlobalAppContext.SetModulesPartial(plugin.AppContextModules{
		IsOfflineRef:        isOfflineRef,
		MediaPlatformRef:    activePlatformRef,
		WSEventManager:      wsEventManager,
		MetadataProviderRef: metadataProviderRef,
	})

	// Initialize online streaming repository
	onlinestreamRepository := onlinestream.NewRepository(&onlinestream.NewRepositoryOptions{
		Logger:              logger,
		FileCacher:          fileCacher,
		MetadataProviderRef: metadataProviderRef,
		PlatformRef:         activePlatformRef,
		Database:            database,
		ExtensionBankRef:    extensionBankRef,
	})

	// Initialize extension playground for testing extensions
	extensionPlaygroundRepository := extension_playground.NewPlaygroundRepository(logger, activePlatformRef, metadataProviderRef)

	// Create the main app instance with initialized components
	app = &App{
		Config:                        cfg,
		Flags:                         configOpts.Flags,
		FeatureManager:                NewFeatureManager(logger, configOpts.Flags),
		Database:                      database,
		MediaClientRef:                mediaClientRef,
		SimklClientRef:                simklClientRef,
		MediaPlatformRef:              activePlatformRef,
		OfflinePlatformRef:            offlinePlatformRef,
		LocalManager:                  localManager,
		WSEventManager:                wsEventManager,
		MediaApiCacheDir:              mediaApiCacheDir,
		SimklCacheDir:                 simklCacheDir,
		Logger:                        logger,
		Version:                       constants.Version,
		Updater:                       updater.New(constants.Version, logger, wsEventManager),
		FileCacher:                    fileCacher,
		OnlinestreamRepository:        onlinestreamRepository,
		MetadataProviderRef:           metadataProviderRef,
		MangaRepository:               mangaRepository,
		ExtensionRepository:           extensionRepository,
		ExtensionBankRef:              extensionBankRef,
		ExtensionPlaygroundRepository: extensionPlaygroundRepository,
		ReportRepository:              report.NewRepository(logger),
		TorrentRepository:             nil, // Initialized in App.initModulesOnce
		FillerManager:                 nil, // Initialized in App.initModulesOnce
		MangaDownloader:               nil, // Initialized in App.initModulesOnce
		PlaybackManager:               nil, // Initialized in App.initModulesOnce
		AutoDownloader:                nil, // Initialized in App.initModulesOnce
		AutoScanner:                   nil, // Initialized in App.initModulesOnce
		MediastreamRepository:         nil, // Initialized in App.initModulesOnce
		TorrentstreamRepository:       nil, // Initialized in App.initModulesOnce
		ContinuityManager:             nil, // Initialized in App.initModulesOnce
		DebridClientRepository:        nil, // Initialized in App.initModulesOnce
		DirectStreamManager:           nil, // Initialized in App.initModulesOnce
		NativePlayer:                  nil, // Initialized in App.initModulesOnce
		VideoCore:                     nil, // Initialized in App.initModulesOnce
		NakamaManager:                 nil, // Initialized in App.initModulesOnce
		LibraryExplorer:               nil, // Initialized in App.initModulesOnce
		TorrentClientRepository:       nil, // Initialized in App.InitOrRefreshModules
		MediaPlayerRepository:         nil, // Initialized in App.InitOrRefreshModules
		DiscordPresence:               nil, // Initialized in App.InitOrRefreshModules
		previousVersion:               previousVersion,
		FeatureFlags:                  NewFeatureFlags(cfg, logger),
		IsDesktopSidecar:              configOpts.Flags.IsDesktopSidecar,
		SecondarySettings: struct {
			Mediastream   *models.MediastreamSettings
			Torrentstream *models.TorrentstreamSettings
			Debrid        *models.DebridSettings
		}{Mediastream: nil, Torrentstream: nil},
		SelfUpdater:                   selfupdater,
		moduleMu:                      sync.Mutex{},
		OnRefreshMediaCollectionFuncs: result.NewMap[string, func()](),
		HookManager:                   hookManager,
		isOfflineRef:                  isOfflineRef,
		ServerPasswordHash:            serverPasswordHash,
		ClientIdentitySecret:          util.GenerateCryptoID(),
	}

	plugin.GlobalAppContext.SetModulesPartial(plugin.AppContextModules{
		PromptManager: extensionRepository.PromptManager(),
		Auth: plugin.AuthActions{
			Login: app.LoginToSimkl,
			Logout: func() error {
				app.LogoutFromSimkl()
				return nil
			},
		},
		Settings: plugin.SettingsActions{
			OnSaved: func(settings *models.Settings) {
				app.WSEventManager.SendEvent("settings", settings)
				app.InitOrRefreshModules()
			},
		},
		Extensions: plugin.ExtensionActions{
			SetDisabled: app.ExtensionRepository.SetExternalExtensionDisabled,
			GetName:     app.ExtensionRepository.GetExtensionName,
		},
	})

	app.SyncSecurityConfig()

	if cfg.Server.SecureMode != "" {
		app.SetSecureMode(cfg.Server.SecureMode, false)
		logger.Warn().Str("mode", cfg.Server.SecureMode).Msg("app: Secure mode configured")
	}

	// Run database migrations if version has changed
	app.runMigrations()

	// Initialize modules that only need to be initialized once
	app.initModulesOnce()

	plugin.GlobalAppContext.SetModulesPartial(plugin.AppContextModules{
		ContinuityManager:       app.ContinuityManager,
		AutoScanner:             app.AutoScanner,
		AutoDownloader:          app.AutoDownloader,
		FileCacher:              app.FileCacher,
		TorrentRepository:       app.TorrentRepository,
		DebridClientRepository:  app.DebridClientRepository,
		OnlinestreamRepository:  app.OnlinestreamRepository,
		MediastreamRepository:   app.MediastreamRepository,
		TorrentstreamRepository: app.TorrentstreamRepository,
	})

	if !app.IsOffline() {
		go app.Updater.FetchAnnouncements()
	}

	// Initialize all modules that depend on settings
	app.InitOrRefreshModules()

	// Load custom source extensions before fetching SIMKL data
	LoadCustomSourceExtensions(extensionRepository)

	// Initialize Simkl data if not in offline mode
	if !app.IsOffline() {
		app.InitOrRefreshMediaData()
	} else {
		app.ServerReady = true
	}

	// Load the other extensions asynchronously
	go LoadExtensions(extensionRepository, logger, cfg)

	// Initialize mediastream settings (for streaming media)
	app.InitOrRefreshMediastreamSettings()

	// Initialize torrentstream settings (for torrent streaming)
	app.InitOrRefreshTorrentstreamSettings()

	// Initialize debrid settings (for debrid services)
	app.InitOrRefreshDebridSettings()

	// Register Nakama manager cleanup
	app.AddCleanupFunction(app.NakamaManager.Cleanup)

	// Run one-time initialization actions
	app.performActionsOnce()

	return app
}

func (a *App) IsOffline() bool {
	return a.isOfflineRef.Get()
}

func (a *App) IsOfflineRef() *util.Ref[bool] {
	return a.isOfflineRef
}

func (a *App) AddCleanupFunction(f func()) {
	a.Cleanups = append(a.Cleanups, f)
}
func (a *App) AddOnRefreshMediaCollectionFunc(key string, f func()) {
	if key == "" {
		return
	}
	a.OnRefreshMediaCollectionFuncs.Set(key, f)
}

func (a *App) Cleanup() {
	for _, f := range a.Cleanups {
		f()
	}
}
