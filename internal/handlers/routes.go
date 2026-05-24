package handlers

import (
	"errors"
	"net/http"
	"path/filepath"
	"seall/internal/core"
	util "seall/internal/util/proxies"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
)

type Handler struct {
	App *core.App
}

func InitRoutes(app *core.App, e *echo.Echo) {
	h := &Handler{App: app}

	e.Use(h.trustedLocalRequestMiddleware)

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return isTrustedCORSOrigin(origin, app.Config.Server.Password, app.Config.Server.AccessAllowlist), nil
		},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Cookie", "Authorization",
			"X-Seall-Token", clientIdHeaderName, clientIdProofHeaderName, clientPlatformHeader,
			"X-Seall-Nakama-Token", "X-Seall-Nakama-Username", "X-Seall-Nakama-Server-Version", "X-Seall-Nakama-Peer-Id"},
		ExposeHeaders:    []string{clientIdHeaderName, clientIdProofHeaderName},
		AllowCredentials: true,
	}))

	lechoLogger := lecho.From(*app.Logger)

	urisToSkip := []string{
		"/internal/metrics",
		"/_next",
		"/icons",
		"/events",
		"/api/v1/image-proxy",
		"/api/v1/mediastream/transcode/",
		"/api/v1/torrent-client/list",
		"/api/v1/proxy",
		"/api/v1/directstream/stream",
	}

	// Logging middleware
	e.Use(lecho.Middleware(lecho.Config{
		Logger: lechoLogger,
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.RequestURI()
			if filepath.Ext(c.Request().URL.Path) == ".txt" ||
				filepath.Ext(c.Request().URL.Path) == ".png" ||
				filepath.Ext(c.Request().URL.Path) == ".ico" {
				return true
			}
			for _, uri := range urisToSkip {
				if uri == path || strings.HasPrefix(path, uri) {
					return true
				}
			}
			return false
		},
		Enricher: func(c echo.Context, logger zerolog.Context) zerolog.Context {
			// Add which file the request came from
			return logger.Str("file", c.Path())
		},
	}))

	// Recovery middleware
	e.Use(middleware.Recover())

	// Client ID middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			cookie, err := c.Cookie(clientIdCookieName)
			cookieClientID := ""
			if err == nil {
				cookieClientID = strings.TrimSpace(cookie.Value)
			}
			clientID := resolveClientIdFromRequest(app, req, cookieClientID)

			if clientID == "" {
				clientID = uuid.New().String()
			}

			if err != nil || cookie == nil || strings.TrimSpace(cookie.Value) != clientID {
				newCookie := new(http.Cookie)
				newCookie.Name = clientIdCookieName
				newCookie.Value = clientID
				newCookie.HttpOnly = true
				newCookie.Expires = time.Now().Add(24 * time.Hour)
				newCookie.Path = "/"
				newCookie.Domain = ""
				newCookie.SameSite = http.SameSiteLaxMode
				newCookie.Secure = requestUsesTrustedHTTPS(req)

				c.SetCookie(newCookie)
			}

			setClientIdentityHeaders(c.Response().Header(), app, clientID)

			c.Set(clientIdCookieName, clientID)
			c.Set(clientPlatformHeader, getClientPlatformFromRequest(req))

			return next(c)
		}
	})

	e.Use(headMethodMiddleware)
	// e.Use(h.controlPlaneBodyLimitMiddleware)
	e.Use(h.controlPlaneMutationRateLimitMiddleware)

	e.GET("/events", h.webSocketEventHandler)

	v1 := e.Group("/api").Group("/v1")

	//
	// Auth middleware
	//
	v1.Use(h.OptionalAuthMiddleware)
	v1.Use(h.FeaturesMiddleware)

	imageProxy := &util.ImageProxy{}
	v1.GET("/image-proxy", imageProxy.ProxyImage)

	v1.GET("/internal/docs", h.HandleGetDocs)

	v1.GET("/proxy", h.VideoProxy)
	v1.HEAD("/proxy", h.VideoProxy)

	v1.GET("/status", h.HandleGetStatus)
	v1.GET("/status/home-items", h.HandleGetHomeItems)
	v1.POST("/status/home-items", h.HandleUpdateHomeItems)

	v1.GET("/log/*", h.HandleGetLogContent)
	v1.GET("/logs/filenames", h.HandleGetLogFilenames)
	v1.DELETE("/logs", h.HandleDeleteLogs)
	v1.GET("/logs/latest", h.HandleGetLatestLogContent)

	v1.GET("/memory/stats", h.HandleGetMemoryStats)
	v1.GET("/memory/profile", h.HandleGetMemoryProfile)
	v1.GET("/memory/goroutine", h.HandleGetGoRoutineProfile)
	v1.GET("/memory/cpu", h.HandleGetCPUProfile)
	v1.POST("/memory/gc", h.HandleForceGC)

	v1.POST("/announcements", h.HandleGetAnnouncements)

	// Auth
	v1.POST("/auth/login", h.HandleLogin)
	v1.POST("/auth/logout", h.HandleLogout)
	v1.PATCH("/auth/simkl/client", h.HandleSaveSimklClientConfig)
	v1.POST("/auth/simkl/pin", h.HandleStartSimklPinLogin)
	v1.POST("/auth/simkl/pin/check", h.HandleCheckSimklPinLogin)

	// Settings
	v1.GET("/settings", h.HandleGetSettings)
	v1.PATCH("/settings", h.HandleSaveSettings)
	v1.POST("/start", h.HandleGettingStarted)
	v1.PATCH("/settings/auto-downloader", h.HandleSaveAutoDownloaderSettings)
	v1.PATCH("/settings/media-player", h.HandleSaveMediaPlayerSettings)

	// Auto Downloader
	v1.POST("/auto-downloader/run", h.HandleRunAutoDownloader)
	v1.POST("/auto-downloader/run/simulation", h.HandleRunAutoDownloaderSimulation)
	v1.GET("/auto-downloader/rule/:id", h.HandleGetAutoDownloaderRule)
	v1.GET("/auto-downloader/rule/media/:id", h.HandleGetAutoDownloaderRulesByAnime)
	v1.GET("/auto-downloader/rules", h.HandleGetAutoDownloaderRules)
	v1.POST("/auto-downloader/rule", h.HandleCreateAutoDownloaderRule)
	v1.PATCH("/auto-downloader/rule", h.HandleUpdateAutoDownloaderRule)
	v1.DELETE("/auto-downloader/rule/:id", h.HandleDeleteAutoDownloaderRule)

	v1.GET("/auto-downloader/items", h.HandleGetAutoDownloaderItems)
	v1.DELETE("/auto-downloader/item", h.HandleDeleteAutoDownloaderItem)

	v1.GET("/auto-downloader/profiles", h.HandleGetAutoDownloaderProfiles)
	v1.GET("/auto-downloader/profile/:id", h.HandleGetAutoDownloaderProfile)
	v1.POST("/auto-downloader/profile", h.HandleCreateAutoDownloaderProfile)
	v1.PATCH("/auto-downloader/profile", h.HandleUpdateAutoDownloaderProfile)
	v1.DELETE("/auto-downloader/profile/:id", h.HandleDeleteAutoDownloaderProfile)

	// Other
	v1.POST("/test-dump", h.HandleTestDump)

	v1.POST("/directory-selector", h.HandleDirectorySelector)

	v1.POST("/open-in-explorer", h.HandleOpenInExplorer)

	v1.POST("/media-player/start", h.HandleStartDefaultMediaPlayer)

	//
	// SIMKL
	//

	v1Simkl := v1.Group("/simkl")

	v1Simkl.GET("/collection", h.HandleGetMediaCollection)
	v1Simkl.POST("/collection", h.HandleGetMediaCollection)

	v1Simkl.GET("/collection/raw", h.HandleGetRawMediaCollection)
	v1Simkl.POST("/collection/raw", h.HandleGetRawMediaCollection)
	v1Simkl.GET("/collection/raw/tags", h.HandleGetRawMediaCollectionTags)

	v1Simkl.GET("/media-details/:id", h.HandleGetMediaDetails)

	v1Simkl.GET("/studio-details/:id", h.HandleGetStudioDetails)

	v1Simkl.POST("/list-entry", h.HandleEditMediaListEntry)

	v1Simkl.DELETE("/list-entry", h.HandleDeleteMediaListEntry)

	v1Simkl.POST("/list-media", h.HandleListMedia)

	v1Simkl.POST("/list-recent-airing-media", h.HandleListRecentAiringMedia)

	v1Simkl.GET("/list-missed-sequels", h.HandleListMissedSequels)

	v1Simkl.GET("/stats", h.HandleGetMediaStats)

	v1Simkl.GET("/cache-layer/status", h.HandleGetMediaCacheLayerStatus)

	v1Simkl.POST("/cache-layer/status", h.HandleToggleMediaCacheLayerStatus)

	//
	// MAL
	//

	v1.POST("/mal/auth", h.HandleMALAuth)

	v1.POST("/mal/logout", h.HandleMALLogout)

	//
	// Library
	//

	v1Library := v1.Group("/library")

	v1Library.POST("/scan", h.HandleScanLocalFiles)

	v1Library.DELETE("/empty-directories", h.HandleRemoveEmptyDirectories)

	v1Library.GET("/local-files", h.HandleGetLocalFiles)
	v1Library.POST("/local-files", h.HandleLocalFileBulkAction)
	v1Library.PATCH("/local-files", h.HandleUpdateLocalFiles)
	v1Library.DELETE("/local-files", h.HandleDeleteLocalFiles)
	v1Library.GET("/local-files/dump", h.HandleDumpLocalFilesToFile)
	v1Library.POST("/local-files/import", h.HandleImportLocalFiles)
	v1Library.PATCH("/local-file", h.HandleUpdateLocalFileData)
	v1Library.PATCH("/local-files/super-update", h.HandleSuperUpdateLocalFiles)

	v1Library.GET("/collection", h.HandleGetLibraryCollection)
	v1Library.GET("/schedule", h.HandleGetMediaCollectionSchedule)

	v1Library.GET("/scan-summaries", h.HandleGetScanSummaries)

	v1Library.GET("/missing-episodes", h.HandleGetMissingEpisodes)
	v1Library.GET("/upcoming-episodes", h.HandleGetUpcomingEpisodes)

	v1Library.GET("/media-entry/:id", h.HandleGetMediaEntry)
	v1Library.POST("/media-entry/suggestions", h.HandleFetchMediaEntrySuggestions)
	v1Library.POST("/media-entry/manual-match", h.HandleMediaEntryManualMatch)
	v1Library.PATCH("/media-entry/bulk-action", h.HandleMediaEntryBulkAction)
	v1Library.POST("/media-entry/open-in-explorer", h.HandleOpenMediaEntryInExplorer)
	v1Library.POST("/media-entry/update-progress", h.HandleUpdateMediaEntryProgress)
	v1Library.POST("/media-entry/update-repeat", h.HandleUpdateMediaEntryRepeat)
	v1Library.GET("/media-entry/silence/:id", h.HandleGetMediaEntrySilenceStatus)
	v1Library.POST("/media-entry/silence", h.HandleToggleMediaEntrySilenceStatus)

	v1Library.POST("/unknown-media", h.HandleAddUnknownMedia)

	//
	// Library Explorer
	//
	v1LibraryExplorer := v1Library.Group("/explorer")

	v1LibraryExplorer.GET("/file-tree", h.HandleGetLibraryExplorerFileTree)
	v1LibraryExplorer.POST("/file-tree/refresh", h.HandleRefreshLibraryExplorerFileTree)
	v1LibraryExplorer.POST("/directory-children", h.HandleLoadLibraryExplorerDirectoryChildren)

	//
	// Anime
	//
	v1.GET("/media/episode-collection/:id", h.HandleGetMediaEpisodeCollection)

	//
	// Torrent / Torrent Client
	//

	v1.POST("/torrent/search", h.HandleSearchTorrent)
	v1.POST("/torrent-client/download", h.HandleTorrentClientDownload)
	v1.GET("/torrent-client/list", h.HandleGetActiveTorrentList)
	v1.POST("/torrent-client/action", h.HandleTorrentClientAction)
	v1.POST("/torrent-client/get-files", h.HandleTorrentClientGetFiles)
	v1.POST("/torrent-client/rule-magnet", h.HandleTorrentClientAddMagnetFromRule)

	//
	// Auto Select
	//

	v1.GET("/auto-select/profile", h.HandleGetAutoSelectProfile)
	v1.POST("/auto-select/profile", h.HandleSaveAutoSelectProfile)
	v1.DELETE("/auto-select/profile", h.HandleDeleteAutoSelectProfile)

	//
	// Download
	//

	v1.POST("/download-torrent-file", h.HandleDownloadTorrentFile)

	//
	// Updates
	//

	v1.GET("/latest-update", h.HandleGetLatestUpdate)
	v1.GET("/changelog", h.HandleGetChangelog)
	v1.POST("/install-update", h.HandleInstallLatestUpdate)
	v1.POST("/download-release", h.HandleDownloadRelease)
	v1.POST("/download-mac-denshi-update", h.HandleDownloadMacDenshiUpdate)
	v1.POST("/check-for-updates", h.HandleCheckForUpdates)

	//
	// Theme
	//

	v1.GET("/theme", h.HandleGetTheme)
	v1.PATCH("/theme", h.HandleUpdateTheme)

	//
	// Playback Manager
	//

	v1.POST("/playback-manager/sync-current-progress", h.HandlePlaybackSyncCurrentProgress)
	v1.POST("/playback-manager/start-playlist", h.HandlePlaybackStartPlaylist)
	v1.POST("/playback-manager/playlist-next", h.HandlePlaybackPlaylistNext)
	v1.POST("/playback-manager/cancel-playlist", h.HandlePlaybackCancelCurrentPlaylist)
	v1.POST("/playback-manager/next-episode", h.HandlePlaybackPlayNextEpisode)
	v1.GET("/playback-manager/next-episode", h.HandlePlaybackGetNextEpisode)
	v1.POST("/playback-manager/autoplay-next-episode", h.HandlePlaybackAutoPlayNextEpisode)
	v1.POST("/playback-manager/play", h.HandlePlaybackPlayVideo)
	v1.POST("/playback-manager/play-random", h.HandlePlaybackPlayRandomVideo)
	//------------
	v1.POST("/playback-manager/manual-tracking/start", h.HandlePlaybackStartManualTracking)
	v1.POST("/playback-manager/manual-tracking/cancel", h.HandlePlaybackCancelManualTracking)

	//
	// Playlists
	//

	v1.GET("/playlists", h.HandleGetPlaylists)
	v1.POST("/playlist", h.HandleCreatePlaylist)
	v1.PATCH("/playlist", h.HandleUpdatePlaylist)
	v1.DELETE("/playlist", h.HandleDeletePlaylist)
	v1.GET("/playlist/episodes/:id", h.HandleGetPlaylistEpisodes)

	//
	// Onlinestream
	//

	v1.POST("/onlinestream/episode-source", h.HandleGetOnlineStreamEpisodeSource)
	v1.POST("/onlinestream/episode-list", h.HandleGetOnlineStreamEpisodeList)
	v1.DELETE("/onlinestream/cache", h.HandleOnlineStreamEmptyCache)

	v1.POST("/onlinestream/search", h.HandleOnlinestreamManualSearch)
	v1.POST("/onlinestream/manual-mapping", h.HandleOnlinestreamManualMapping)
	v1.POST("/onlinestream/get-mapping", h.HandleGetOnlinestreamMapping)
	v1.POST("/onlinestream/remove-mapping", h.HandleRemoveOnlinestreamMapping)

	//
	// Metadata Provider
	//

	v1.POST("/metadata-provider/filler", h.HandlePopulateFillerData)
	v1.DELETE("/metadata-provider/filler", h.HandleRemoveFillerData)
	v1.GET("/metadata/parent/:id", h.HandleGetMediaMetadataParent)
	v1.POST("/metadata/parent", h.HandleSaveMediaMetadataParent)
	v1.DELETE("/metadata/parent", h.HandleDeleteMediaMetadataParent)

	//
	// Reading
	//

	v1Reading := v1.Group("/reading")
	v1Reading.POST("/simkl/collection", h.HandleGetSimklMangaCollection)
	v1Reading.GET("/simkl/collection/raw", h.HandleGetRawMangaCollection)
	v1Reading.POST("/simkl/collection/raw", h.HandleGetRawMangaCollection)
	v1Reading.GET("/simkl/collection/raw/tags", h.HandleGetRawMangaCollectionTags)
	v1Reading.POST("/simkl/list", h.HandleListManga)
	v1Reading.GET("/collection", h.HandleGetMangaCollection)
	v1Reading.GET("/latest-chapter-numbers", h.HandleGetMangaLatestChapterNumbersMap)
	v1Reading.POST("/refetch-chapter-containers", h.HandleRefetchMangaChapterContainers)
	v1Reading.GET("/entry/:id", h.HandleGetMangaEntry)
	v1Reading.GET("/entry/:id/details", h.HandleGetMangaEntryDetails)
	v1Reading.DELETE("/entry/cache", h.HandleEmptyMangaEntryCache)
	v1Reading.POST("/chapters", h.HandleGetMangaEntryChapters)
	v1Reading.POST("/pages", h.HandleGetMangaEntryPages)
	v1Reading.POST("/update-progress", h.HandleUpdateMangaProgress)
	v1Reading.GET("/downloaded-chapters/:id", h.HandleGetMangaEntryDownloadedChapters)
	v1Reading.GET("/downloads", h.HandleGetMangaDownloadsList)
	v1Reading.POST("/download-chapters", h.HandleDownloadMangaChapters)
	v1Reading.POST("/download-data", h.HandleGetMangaDownloadData)
	v1Reading.DELETE("/download-chapter", h.HandleDeleteMangaDownloadedChapters)
	v1Reading.GET("/download-queue", h.HandleGetMangaDownloadQueue)
	v1Reading.POST("/download-queue/start", h.HandleStartMangaDownloadQueue)
	v1Reading.POST("/download-queue/stop", h.HandleStopMangaDownloadQueue)
	v1Reading.DELETE("/download-queue", h.HandleClearAllChapterDownloadQueue)
	v1Reading.POST("/download-queue/reset-errored", h.HandleResetErroredChapterDownloadQueue)
	v1Reading.POST("/search", h.HandleMangaManualSearch)
	v1Reading.POST("/manual-mapping", h.HandleMangaManualMapping)
	v1Reading.POST("/get-mapping", h.HandleGetMangaMapping)
	v1Reading.POST("/remove-mapping", h.HandleRemoveMangaMapping)
	v1Reading.GET("/local-page/:path", h.HandleGetLocalMangaPage)

	//
	// File Cache
	//

	v1FileCache := v1.Group("/filecache")
	v1FileCache.GET("/total-size", h.HandleGetFileCacheTotalSize)
	v1FileCache.DELETE("/bucket", h.HandleRemoveFileCacheBucket)
	v1FileCache.GET("/mediastream/videofiles/total-size", h.HandleGetFileCacheMediastreamVideoFilesTotalSize)
	v1FileCache.DELETE("/mediastream/videofiles", h.HandleClearFileCacheMediastreamVideoFiles)

	//
	// Discord
	//

	v1Discord := v1.Group("/discord")
	v1Discord.POST("/presence/reading", h.HandleSetDiscordMangaActivity)
	v1Discord.POST("/presence/legacy-media", h.HandleSetDiscordLegacyAnimeActivity)
	v1Discord.POST("/presence/media", h.HandleSetDiscordAnimeActivityWithProgress)
	v1Discord.POST("/presence/media-update", h.HandleUpdateDiscordAnimeActivityWithProgress)
	v1Discord.POST("/presence/cancel", h.HandleCancelDiscordActivity)

	//
	// Media Stream
	//
	v1.GET("/mediastream/settings", h.HandleGetMediastreamSettings)
	v1.PATCH("/mediastream/settings", h.HandleSaveMediastreamSettings)
	v1.POST("/mediastream/request", h.HandleRequestMediastreamMediaContainer)
	v1.POST("/mediastream/preload", h.HandlePreloadMediastreamMediaContainer)
	// Transcode
	v1.POST("/mediastream/shutdown-transcode", h.HandleMediastreamShutdownTranscodeStream)
	v1.GET("/mediastream/transcode/*", h.HandleMediastreamTranscode)
	v1.GET("/mediastream/subs/*", h.HandleMediastreamGetSubtitles)
	v1.GET("/mediastream/att/*", h.HandleMediastreamGetAttachments)
	v1.GET("/mediastream/direct", h.HandleMediastreamDirectPlay)
	v1.HEAD("/mediastream/direct", h.HandleMediastreamDirectPlay)
	v1.GET("/mediastream/file", h.HandleMediastreamFile)
	v1.GET("/mediastream/local-subtitles", h.HandleMediastreamLocalSubtitles)

	//
	// Direct Stream
	//
	v1.POST("/directstream/play/localfile", h.HandleDirectstreamPlayLocalFile)
	v1.GET("/directstream/stream", echo.WrapHandler(h.HandleDirectstreamGetStream()))
	v1.HEAD("/directstream/stream", echo.WrapHandler(h.HandleDirectstreamGetStream()))
	v1.GET("/directstream/att/*", h.HandleDirectstreamGetAttachments)
	v1.POST("/directstream/subs/convert-subs", h.HandleDirectstreamConvertSubs)

	//
	// VideoCore
	//
	v1.GET("/videocore/insight/character/:malId", h.HandleVideoCoreInSightGetCharacterDetails)

	//
	// Torrent stream
	//
	v1.GET("/torrentstream/settings", h.HandleGetTorrentstreamSettings)
	v1.PATCH("/torrentstream/settings", h.HandleSaveTorrentstreamSettings)
	v1.POST("/torrentstream/start", h.HandleTorrentstreamStartStream)
	v1.POST("/torrentstream/stop", h.HandleTorrentstreamStopStream)
	v1.POST("/torrentstream/drop", h.HandleTorrentstreamDropTorrent)
	v1.POST("/torrentstream/torrent-file-previews", h.HandleGetTorrentstreamTorrentFilePreviews)
	v1.POST("/torrentstream/batch-history", h.HandleGetTorrentstreamBatchHistory)
	v1.POST("/torrentstream/batch-history/delete", h.HandleDeleteTorrentstreamBatchHistory)
	v1.GET("/torrentstream/stream/*", h.HandleTorrentstreamServeStream)

	//
	// Extensions
	//

	v1Extensions := v1.Group("/extensions")
	v1Extensions.POST("/playground/run", h.HandleRunExtensionPlaygroundCode)
	v1Extensions.POST("/external/fetch", h.HandleFetchExternalExtensionData)
	v1Extensions.POST("/external/install", h.HandleInstallExternalExtension)
	v1Extensions.POST("/external/install-repository", h.HandleInstallExternalExtensionRepository)
	v1Extensions.POST("/external/uninstall", h.HandleUninstallExternalExtension)
	v1Extensions.POST("/external/edit-payload", h.HandleUpdateExtensionCode)
	v1Extensions.POST("/external/reload", h.HandleReloadExternalExtensions)
	v1Extensions.POST("/external/reload", h.HandleReloadExternalExtension)
	v1Extensions.POST("/external/disabled", h.HandleSetExternalExtensionDisabled)
	v1Extensions.POST("/all", h.HandleGetAllExtensions)
	v1Extensions.GET("/updates", h.HandleGetExtensionUpdateData)
	v1Extensions.GET("/list", h.HandleListExtensionData)
	v1Extensions.GET("/payload/:id", h.HandleGetExtensionPayload)
	v1Extensions.GET("/list/development", h.HandleListDevelopmentModeExtensions)
	v1Extensions.GET("/list/reading-provider", h.HandleListMangaProviderExtensions)
	v1Extensions.GET("/list/onlinestream-provider", h.HandleListOnlinestreamProviderExtensions)
	v1Extensions.GET("/list/media-torrent-provider", h.HandleListMediaTorrentProviderExtensions)
	v1Extensions.GET("/list/media-entry-episode-tabs", h.HandleListMediaEntryEpisodeTabExtensions)
	v1Extensions.GET("/list/custom-source", h.HandleListCustomSourceExtensions)
	v1Extensions.GET("/user-config/:id", h.HandleGetExtensionUserConfig)
	v1Extensions.POST("/user-config", h.HandleSaveExtensionUserConfig)
	v1Extensions.GET("/marketplace", h.HandleGetMarketplaceExtensions)
	v1Extensions.GET("/plugin-settings", h.HandleGetPluginSettings)
	v1Extensions.POST("/plugin-settings/pinned-trays", h.HandleSetPluginSettingsPinnedTrays)
	v1Extensions.POST("/plugin-permissions/grant", h.HandleGrantPluginPermissions)

	//
	// Continuity
	//
	v1Continuity := v1.Group("/continuity")
	v1Continuity.PATCH("/item", h.HandleUpdateContinuityWatchHistoryItem)
	v1Continuity.GET("/item/:id", h.HandleGetContinuityWatchHistoryItem)
	v1Continuity.GET("/history", h.HandleGetContinuityWatchHistory)

	//
	// Sync
	//
	v1Local := v1.Group("/local")
	v1Local.GET("/track", h.HandleLocalGetTrackedMediaItems)
	v1Local.POST("/track", h.HandleLocalAddTrackedMedia)
	v1Local.DELETE("/track", h.HandleLocalRemoveTrackedMedia)
	v1Local.GET("/track/:id/:type", h.HandleLocalGetIsMediaTracked)
	v1Local.POST("/local", h.HandleLocalSyncData)
	v1Local.GET("/queue", h.HandleLocalGetSyncQueueState)
	v1Local.POST("/simkl", h.HandleLocalSyncMediaData)
	v1Local.POST("/updated", h.HandleLocalSetHasLocalChanges)
	v1Local.GET("/updated", h.HandleLocalGetHasLocalChanges)
	v1Local.GET("/storage/size", h.HandleLocalGetLocalStorageSize)
	v1Local.POST("/sync-simulated-to-simkl", h.HandleLocalSyncSimulatedDataToMedia)

	v1Local.POST("/offline", h.HandleSetOfflineMode)

	//
	// Debrid
	//

	v1.GET("/debrid/settings", h.HandleGetDebridSettings)
	v1.PATCH("/debrid/settings", h.HandleSaveDebridSettings)
	v1.POST("/debrid/torrents", h.HandleDebridAddTorrents)
	v1.POST("/debrid/torrents/download", h.HandleDebridDownloadTorrent)
	v1.POST("/debrid/torrents/cancel", h.HandleDebridCancelDownload)
	v1.DELETE("/debrid/torrent", h.HandleDebridDeleteTorrent)
	v1.GET("/debrid/torrents", h.HandleDebridGetTorrents)
	v1.POST("/debrid/torrents/info", h.HandleDebridGetTorrentInfo)
	v1.POST("/debrid/torrents/file-previews", h.HandleDebridGetTorrentFilePreviews)
	v1.POST("/debrid/stream/start", h.HandleDebridStartStream)
	v1.POST("/debrid/stream/cancel", h.HandleDebridCancelStream)

	//
	// Report
	//

	v1.POST("/report/issue", h.HandleSaveIssueReport)
	v1.GET("/report/issue/download", h.HandleDownloadIssueReport)
	v1.POST("/report/issue/decompress", h.HandleDecompressIssueReport)

	//
	// Nakama
	//

	v1Nakama := v1.Group("/nakama")
	v1Nakama.GET("/ws", h.HandleNakamaWebSocket)
	v1Nakama.POST("/message", h.HandleSendNakamaMessage)
	v1Nakama.POST("/reconnect", h.HandleNakamaReconnectToHost)
	v1Nakama.POST("/cleanup", h.HandleNakamaRemoveStaleConnections)
	v1Nakama.GET("/room/available", h.HandleNakamaRoomsAvailable)
	v1Nakama.POST("/room/create", h.HandleNakamaCreateAndJoinRoom)
	v1Nakama.POST("/room/disconnect", h.HandleNakamaDisconnectFromRoom)
	v1Nakama.GET("/host/media/library", h.HandleGetNakamaAnimeLibrary)
	v1Nakama.GET("/host/media/library/shared", h.HandleGetNakamaAnimeLibraryShared)
	v1Nakama.GET("/host/media/library/files/:id", h.HandleGetNakamaAnimeLibraryFiles)
	v1Nakama.GET("/host/media/library/files", h.HandleGetNakamaAnimeAllLibraryFiles)
	v1Nakama.POST("/play", h.HandleNakamaPlayVideo)
	v1Nakama.GET("/host/torrentstream/stream", h.HandleNakamaHostTorrentstreamServeStream)
	v1Nakama.HEAD("/host/torrentstream/stream", h.HandleNakamaHostTorrentstreamServeStream)
	v1Nakama.GET("/host/media/library/stream", h.HandleNakamaHostAnimeLibraryServeStream)
	v1Nakama.HEAD("/host/media/library/stream", h.HandleNakamaHostAnimeLibraryServeStream)
	v1Nakama.GET("/host/debridstream/stream", h.HandleNakamaHostDebridstreamServeStream)
	v1Nakama.HEAD("/host/debridstream/stream", h.HandleNakamaHostDebridstreamServeStream)
	v1Nakama.GET("/host/debridstream/url", h.HandleNakamaHostGetDebridstreamURL)
	v1Nakama.GET("/stream", h.HandleNakamaProxyStream)
	v1Nakama.HEAD("/stream", h.HandleNakamaProxyStream)
	v1Nakama.POST("/watch-party/create", h.HandleNakamaCreateWatchParty)
	v1Nakama.POST("/watch-party/join", h.HandleNakamaJoinWatchParty)
	v1Nakama.POST("/watch-party/leave", h.HandleNakamaLeaveWatchParty)
	v1Nakama.POST("/watch-party/chat", h.HandleNakamaSendChatMessage)

	//
	// Custom Source
	//
	v1CustomSource := v1.Group("/custom-source")
	v1CustomSource.POST("/provider/list/media", h.HandleCustomSourceListMedia)
	v1CustomSource.POST("/provider/list/reading", h.HandleCustomSourceListManga)

}

func (h *Handler) JSON(c echo.Context, code int, i interface{}) error {
	return c.JSON(code, i)
}

func (h *Handler) RespondWithData(c echo.Context, data interface{}) error {
	return c.JSON(200, NewDataResponse(data))
}

func (h *Handler) RespondWithError(c echo.Context, err error) error {
	return c.JSON(statusCodeForError(err), NewErrorResponse(err))
}

func (h *Handler) RespondWithStatusError(c echo.Context, code int, err error) error {
	return c.JSON(code, NewErrorResponse(err))
}

func statusCodeForError(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	if echoErr, ok := errors.AsType[*echo.HTTPError](err); ok && echoErr.Code >= 400 && echoErr.Code < 600 {
		return echoErr.Code
	}

	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return http.StatusRequestEntityTooLarge
	}

	if strings.EqualFold(strings.TrimSpace(err.Error()), "UNAUTHENTICATED") {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}

func headMethodMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Skip stream routes
		if strings.Contains(c.Request().URL.Path, "/directstream/stream") ||
			strings.Contains(c.Request().URL.Path, "/nakama") {
			return next(c)
		}

		if c.Request().Method == http.MethodHead {
			// Set the method to GET temporarily to reuse the handler
			c.Request().Method = http.MethodGet

			defer func() {
				c.Request().Method = http.MethodHead
			}() // Restore method after

			// Call the next handler and then clear the response body
			if err := next(c); err != nil {
				if err.Error() == echo.ErrMethodNotAllowed.Error() {
					return c.NoContent(http.StatusOK)
				}

				return err
			}
		}

		return next(c)
	}
}
