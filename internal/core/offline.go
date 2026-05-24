package core

import (
	"seall/internal/api/metadata_provider"
	"seall/internal/platforms/offline_platform"
	"seall/internal/platforms/simkl_platform"

	"github.com/spf13/viper"
)

// SetOfflineMode changes the offline mode.
// It updates the config and active media platform.
func (a *App) SetOfflineMode(enabled bool) {
	// Update the config
	a.Config.Server.Offline = enabled
	viper.Set("server.offline", enabled)
	err := viper.WriteConfig()
	if err != nil {
		a.Logger.Err(err).Msg("app: Failed to write config after setting offline mode")
	}
	a.Logger.Info().Bool("enabled", enabled).Msg("app: Offline mode set")
	a.isOfflineRef.Set(enabled)

	if a.MediaPlatformRef.IsPresent() {
		a.MediaPlatformRef.Get().Close()
	}
	if a.MetadataProviderRef.IsPresent() {
		a.MetadataProviderRef.Get().Close()
	}

	// Update the platform and metadata provider
	if enabled {
		if a.NakamaManager.IsConnectedToHost() || a.NakamaManager.IsHost() {
			a.NakamaManager.Stop()
		}

		mediaPlatform, _ := offline_platform.NewOfflinePlatform(a.LocalManager, a.MediaClientRef, a.Logger)
		a.MediaPlatformRef.Set(mediaPlatform)
		a.MetadataProviderRef.Set(a.LocalManager.GetOfflineMetadataProvider())
	} else {
		// DEVNOTE: We don't handle local platform since the feature doesn't allow offline mode
		simklPlatform := simkl_platform.NewSimklPlatform(a.SimklClientRef, a.MediaClientRef, a.ExtensionBankRef, a.Logger, a.Database)
		a.MediaPlatformRef.Set(simklPlatform)
		a.MetadataProviderRef.Set(metadata_provider.NewProvider(&metadata_provider.NewProviderImplOptions{
			Logger:           a.Logger,
			FileCacher:       a.FileCacher,
			ExtensionBankRef: a.ExtensionBankRef,
			Database:         a.Database,
			SimklClientRef:   a.SimklClientRef,
		}))
		a.InitOrRefreshMediaData()
	}
	a.AddOnRefreshMediaCollectionFunc("simkl-platform", func() {
		a.MediaPlatformRef.Get().ClearCache()
	})

	a.InitOrRefreshModules()
}
