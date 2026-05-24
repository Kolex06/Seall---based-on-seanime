package core

import (
	"seall/internal/constants"
	"seall/internal/extension"
	"seall/internal/extension_repo"
	manga_providers "seall/internal/manga/providers"

	"github.com/rs/zerolog"
)

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

	// Load external extensions
	//extensionRepository.ReloadExternalExtensions()
	extensionRepository.LoadOnlyWrapper([]extension.Type{extension.TypeMangaProvider, extension.TypeOnlinestreamProvider, extension.TypeMediaTorrentProvider, extension.TypePlugin}, func() {
		extensionRepository.ReloadExternalExtensions()
	})
}
