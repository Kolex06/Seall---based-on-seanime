package extension

import (
	hibiketorrent "seall/internal/extension/hibike/torrent"
)

type MediaTorrentProviderExtension interface {
	BaseExtension
	GetProvider() hibiketorrent.AnimeProvider
}

type MediaTorrentProviderExtensionImpl struct {
	ext      *Extension
	provider hibiketorrent.AnimeProvider
}

func NewMediaTorrentProviderExtension(ext *Extension, provider hibiketorrent.AnimeProvider) MediaTorrentProviderExtension {
	return &MediaTorrentProviderExtensionImpl{
		ext:      ext,
		provider: provider,
	}
}

func (m *MediaTorrentProviderExtensionImpl) GetProvider() hibiketorrent.AnimeProvider {
	return m.provider
}

func (m *MediaTorrentProviderExtensionImpl) GetExtension() *Extension {
	return m.ext
}

func (m *MediaTorrentProviderExtensionImpl) GetType() Type {
	return m.ext.Type
}

func (m *MediaTorrentProviderExtensionImpl) GetID() string {
	return m.ext.ID
}

func (m *MediaTorrentProviderExtensionImpl) GetName() string {
	return m.ext.Name
}

func (m *MediaTorrentProviderExtensionImpl) GetVersion() string {
	return m.ext.Version
}

func (m *MediaTorrentProviderExtensionImpl) GetManifestURI() string {
	return m.ext.ManifestURI
}

func (m *MediaTorrentProviderExtensionImpl) GetLanguage() Language {
	return m.ext.Language
}

func (m *MediaTorrentProviderExtensionImpl) GetLang() string {
	return GetExtensionLang(m.ext.Lang)
}

func (m *MediaTorrentProviderExtensionImpl) GetDescription() string {
	return m.ext.Description
}

func (m *MediaTorrentProviderExtensionImpl) GetNotes() string {
	return m.ext.Notes
}

func (m *MediaTorrentProviderExtensionImpl) GetAuthor() string {
	return m.ext.Author
}

func (m *MediaTorrentProviderExtensionImpl) GetPayload() string {
	return m.ext.Payload
}

func (m *MediaTorrentProviderExtensionImpl) GetWebsite() string {
	return m.ext.Website
}

func (m *MediaTorrentProviderExtensionImpl) GetReadme() string {
	return m.ext.Readme
}

func (m *MediaTorrentProviderExtensionImpl) GetIcon() string {
	return m.ext.Icon
}

func (m *MediaTorrentProviderExtensionImpl) GetPermissions() []string {
	return m.ext.Permissions
}

func (m *MediaTorrentProviderExtensionImpl) GetUserConfig() *UserConfig {
	return m.ext.UserConfig
}

func (m *MediaTorrentProviderExtensionImpl) GetSavedUserConfig() *SavedUserConfig {
	return m.ext.SavedUserConfig
}

func (m *MediaTorrentProviderExtensionImpl) GetPayloadURI() string {
	return m.ext.PayloadURI
}

func (m *MediaTorrentProviderExtensionImpl) GetIsDevelopment() bool {
	return m.ext.IsDevelopment
}

func (m *MediaTorrentProviderExtensionImpl) GetPluginManifest() *PluginManifest {
	return m.ext.Plugin
}
