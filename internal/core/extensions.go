package core

import (
	_ "embed"

	"seall/internal/constants"
	"seall/internal/extension"
	"seall/internal/extension_repo"
	manga_providers "seall/internal/manga/providers"

	"github.com/rs/zerolog"
)

//go:embed builtin_plugins/custom_css_manager/provider.ts
var customCSSManagerPayload string

func LoadCustomSourceExtensions(extensionRepository *extension_repo.Repository) {
	extensionRepository.LoadOnlyWrapper([]extension.Type{extension.TypeCustomSource}, func() {
		extensionRepository.ReloadExternalExtensions()
	})
}

func LoadExtensions(extensionRepository *extension_repo.Repository, logger *zerolog.Logger, config *Config) {
	// Load built-in extensions
	extensionRepository.ReloadBuiltInExtension(extension.Extension{
		ID:          manga_providers.LocalProvider,
		Name:        "Local",
		Version:     "",
		ManifestURI: "builtin",
		Language:    extension.LanguageGo,
		Type:        extension.TypeMangaProvider,
		Author:      constants.AppName,
		Lang:        "multi",
		Icon:        constants.ProjectRawMainUrl + "/seall-denshi/assets/icons/256x256.png",
	}, manga_providers.NewLocal(config.Manga.LocalDir, logger))

	extensionRepository.ReloadBuiltInExtension(extension.Extension{
		ID:          "custom-css-manager",
		Name:        "Custom CSS Manager",
		Version:     "1.0.10",
		ManifestURI: "builtin",
		Language:    extension.LanguageTypescript,
		Type:        extension.TypePlugin,
		Description: "Manage Seall custom CSS presets.",
		Author:      constants.AppName,
		Lang:        "en",
		Icon:        constants.ProjectRawMainUrl + "/seall-denshi/assets/icons/256x256.png",
		Website:     constants.ProjectRepositoryUrl,
		Payload:     customCSSManagerPayload,
		Plugin: &extension.PluginManifest{
			Version: extension.PluginManifestVersion,
			Permissions: extension.PluginPermissions{
				Scopes: []extension.PluginPermissionScope{extension.PluginPermissionStorage},
				Allow: extension.PluginAllowlist{
					NetworkAccess: extension.PluginNetworkAcess{
						AllowedDomains: []string{"raw.githubusercontent.com", "jigsaw.w3.org"},
						Reasoning:      "Used for optional community style refreshes and CSS validation.",
					},
				},
			},
		},
	}, nil)

	// Load external extensions
	//extensionRepository.ReloadExternalExtensions()
	extensionRepository.LoadOnlyWrapper([]extension.Type{extension.TypeMangaProvider, extension.TypeOnlinestreamProvider, extension.TypeMediaTorrentProvider, extension.TypePlugin}, func() {
		extensionRepository.ReloadExternalExtensions()
	})
}
