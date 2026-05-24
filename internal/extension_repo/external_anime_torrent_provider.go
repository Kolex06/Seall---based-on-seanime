package extension_repo

import (
	"fmt"
	"seall/internal/extension"
	"seall/internal/util"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Anime Torrent provider
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *Repository) loadExternalMediaTorrentProviderExtension(ext *extension.Extension) (err error) {
	defer util.HandlePanicInModuleWithError("extension_repo/loadExternalMediaTorrentProviderExtension", &err)

	switch ext.Language {
	case extension.LanguageJavascript, extension.LanguageTypescript:
		err = r.loadExternalMediaTorrentProviderExtensionJS(ext, ext.Language)
	default:
		err = fmt.Errorf("unsupported language: %v", ext.Language)
	}

	if err != nil {
		return
	}

	return
}

func (r *Repository) loadExternalMediaTorrentProviderExtensionJS(ext *extension.Extension, language extension.Language) error {
	provider, gojaExt, err := NewGojaMediaTorrentProvider(ext, language, r.logger, r.gojaRuntimeManager, r.wsEventManager)
	if err != nil {
		return err
	}

	// Add the extension to the map
	retExt := extension.NewMediaTorrentProviderExtension(ext, provider)
	r.extensionBankRef.Get().Set(ext.ID, retExt)
	r.gojaExtensions.Set(ext.ID, gojaExt)
	return nil
}
