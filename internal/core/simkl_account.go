package core

import (
	"context"
	"errors"
	"seall/internal/api/mediaapi"
	"seall/internal/api/simkl"
	"seall/internal/database/models"
	"seall/internal/events"
	"seall/internal/platforms/platform"
	"seall/internal/platforms/simkl_platform"
	"seall/internal/platforms/simulated_platform"
	"seall/internal/user"
	"seall/internal/util"
	"strconv"
	"time"

	"github.com/goccy/go-json"
)

// GetUser returns the currently logged-in user or a simulated one.
func (a *App) GetUser() *user.User {
	if a.user == nil {
		return user.NewSimulatedUser()
	}
	return a.user
}

// GetUsername returns the username of the currently logged-in user
func (a *App) GetUsername() string {
	if a.user == nil {
		return ""
	}
	if a.user.Viewer == nil {
		return ""
	}
	return a.user.Viewer.GetName()
}

func (a *App) GetUserSimklToken() string {
	if a.user == nil || a.user.Token == user.SimulatedUserToken {
		return ""
	}

	return a.user.Token
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// UpdatePlatform changes the current platform to the provided one.
func (a *App) UpdatePlatform(platform platform.Platform) {
	if a.MediaPlatformRef.IsPresent() {
		a.MediaPlatformRef.Get().Close()
	}
	a.MediaPlatformRef.Set(platform)
	a.AddOnRefreshMediaCollectionFunc("simkl-platform", func() {
		a.MediaPlatformRef.Get().ClearCache()
	})
}

// UpdateMediaClientToken updates the generated media API compatibility client token.
// This function should be called when a user logs in
func (a *App) UpdateMediaClientToken(token string) {
	ac := mediaapi.NewMediaApiClient(token, a.MediaApiCacheDir)
	a.MediaClientRef.Set(ac)
}

func (a *App) UpdateSimklClientToken(token string) {
	sc := simkl.NewClient(simkl.NewClientOptions{
		Token:        token,
		ClientID:     a.Config.Simkl.ClientID,
		ClientSecret: a.Config.Simkl.ClientSecret,
		RedirectURI:  a.Config.Simkl.RedirectURI,
		CacheDir:     a.SimklCacheDir,
		Logger:       a.Logger,
	})
	a.SimklClientRef.Set(sc)
}

func (a *App) LoginToSimkl(token string) error {
	if token == "" {
		return errors.New("token is empty")
	}

	a.UpdateSimklClientToken(token)

	settings, err := a.SimklClientRef.Get().Settings(context.Background())
	if err != nil {
		a.Logger.Error().Msg("Could not authenticate to SIMKL")
		return err
	}

	viewer := simklUserToViewer(settings)
	if len(viewer.Name) == 0 {
		return errors.New("could not find user")
	}

	bytes, err := json.Marshal(viewer)
	if err != nil {
		a.Logger.Err(err).Msg("scan: could not save local files")
	}

	_, err = a.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			ID:        1,
			UpdatedAt: time.Now(),
		},
		Username: viewer.Name,
		Token:    token,
		Viewer:   bytes,
	})
	if err != nil {
		return err
	}

	a.Logger.Info().Msg("app: Authenticated to SIMKL")

	simklPlatform := simkl_platform.NewSimklPlatform(a.SimklClientRef, a.MediaClientRef, a.ExtensionBankRef, a.Logger, a.Database)
	a.UpdatePlatform(simklPlatform)

	a.InitOrRefreshMediaData()
	a.InitOrRefreshModules()

	go func() {
		defer util.HandlePanicThen(func() {})
		a.InitOrRefreshTorrentstreamSettings()
		a.InitOrRefreshMediastreamSettings()
		a.InitOrRefreshDebridSettings()
	}()

	return nil
}

func simklUserToViewer(settings *simkl.UserSettings) *mediaapi.GetViewer_Viewer {
	if settings == nil {
		return &mediaapi.GetViewer_Viewer{}
	}
	name := settings.User.Name
	if name == "" && settings.Account.ID != 0 {
		name = "SIMKL " + strconv.Itoa(settings.Account.ID)
	}
	avatar := settings.User.Avatar
	return &mediaapi.GetViewer_Viewer{
		Name: name,
		Avatar: &mediaapi.GetViewer_Viewer_Avatar{
			Large:  emptyStringToNil(avatar),
			Medium: emptyStringToNil(avatar),
		},
	}
}

func emptyStringToNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// LogoutFromSimkl clears the SIMKL token and switches to the simulated platform.
// This is called internally when the token is detected as invalid.
func (a *App) LogoutFromSimkl() {
	// prevent multiple concurrent calls (e.g. from parallel failing requests)
	if !a.logoutInProgress.CompareAndSwap(false, true) {
		return
	}
	defer a.logoutInProgress.Store(false)

	a.UpdateMediaClientToken("")
	a.UpdateSimklClientToken("")

	simulatedPlatform, err := simulated_platform.NewSimulatedPlatform(a.LocalManager, a.MediaClientRef, a.ExtensionBankRef, a.Logger, a.Database)
	if err != nil {
		a.Logger.Error().Err(err).Msg("app: Failed to create simulated platform during auto-logout")
	} else {
		a.UpdatePlatform(simulatedPlatform)
	}

	_, _ = a.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			ID:        1,
			UpdatedAt: time.Now(),
		},
		Username: "",
		Token:    "",
		Viewer:   nil,
	})

	a.Logger.Debug().Msg("app: Logged out from SIMKL, switched to simulated platform")

	a.InitOrRefreshModules()
	a.InitOrRefreshMediaData()
}

// GetMediaCollection returns the user's SIMKL media collection.
// When bypassCache is true, it will always query SIMKL for the user's collection.
func (a *App) GetMediaCollection(bypassCache bool) (*mediaapi.AnimeCollection, error) {
	return a.MediaPlatformRef.Get().GetMediaCollection(context.Background(), bypassCache)
}

// GetRawMediaCollection is the same as GetMediaCollection but returns the raw collection that includes custom lists
func (a *App) GetRawMediaCollection(bypassCache bool) (*mediaapi.AnimeCollection, error) {
	return a.MediaPlatformRef.Get().GetRawMediaCollection(context.Background(), bypassCache)
}

func (a *App) SyncMediaToSimulatedCollection() {
	if a.LocalManager != nil &&
		!a.GetUser().IsSimulated &&
		a.Settings != nil &&
		a.Settings.Library != nil &&
		a.Settings.Library.AutoSyncToLocalAccount {
		_ = a.LocalManager.SynchronizeMediaToSimulatedCollection()
	}
}

// RefreshAnimeCollection queries SIMKL for the user's media collection.
func (a *App) RefreshAnimeCollection() (*mediaapi.AnimeCollection, error) {
	go func() {
		a.OnRefreshMediaCollectionFuncs.Range(func(key string, f func()) bool {
			go f()
			return true
		})
	}()

	ret, err := a.MediaPlatformRef.Get().RefreshAnimeCollection(context.Background())

	if err != nil {
		return nil, err
	}

	// Save the collection to PlaybackManager
	a.PlaybackManager.SetAnimeCollection(ret)

	// Save the collection to AutoDownloader
	a.AutoDownloader.SetAnimeCollection(ret)

	// Save the collection to LocalManager
	a.LocalManager.SetAnimeCollection(ret)

	// Save the collection to DirectStreamManager
	a.DirectStreamManager.SetAnimeCollection(ret)

	// Save the collection to LibraryExplorer
	a.LibraryExplorer.SetAnimeCollection(ret)

	a.AutoScanner.SetAnimeCollection(ret)

	//a.SyncMediaToSimulatedCollection()

	a.WSEventManager.SendEvent(events.RefreshedMediaCollection, nil)

	return ret, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetMangaCollection is the same as GetMediaCollection but for manga
func (a *App) GetMangaCollection(bypassCache bool) (*mediaapi.MangaCollection, error) {
	return a.MediaPlatformRef.Get().GetMangaCollection(context.Background(), bypassCache)
}

// GetRawMangaCollection does not exclude custom lists
func (a *App) GetRawMangaCollection(bypassCache bool) (*mediaapi.MangaCollection, error) {
	return a.MediaPlatformRef.Get().GetRawMangaCollection(context.Background(), bypassCache)
}

// RefreshMangaCollection queries Simkl for the user's manga collection
func (a *App) RefreshMangaCollection() (*mediaapi.MangaCollection, error) {
	mc, err := a.MediaPlatformRef.Get().RefreshMangaCollection(context.Background())

	if err != nil {
		return nil, err
	}

	a.LocalManager.SetMangaCollection(mc)

	a.WSEventManager.SendEvent(events.RefreshedMangaCollection, nil)

	return mc, nil
}
