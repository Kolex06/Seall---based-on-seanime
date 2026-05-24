package handlers

import (
	"errors"
	"testing"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"seall/internal/api/metadata_provider"
	"seall/internal/database/models"
	"seall/internal/extension"
	hibiketorrent "seall/internal/extension/hibike/torrent"
	"seall/internal/library/autodownloader"
	torrentrepo "seall/internal/torrents/torrent"
	"seall/internal/util"
)

type ruleMagnetTestProvider struct {
	magnet string
	err    error
	calls  int
}

func (p *ruleMagnetTestProvider) Search(hibiketorrent.AnimeSearchOptions) ([]*hibiketorrent.MediaTorrent, error) {
	return nil, nil
}

func (p *ruleMagnetTestProvider) SmartSearch(hibiketorrent.AnimeSmartSearchOptions) ([]*hibiketorrent.MediaTorrent, error) {
	return nil, nil
}

func (p *ruleMagnetTestProvider) GetTorrentInfoHash(torrent *hibiketorrent.MediaTorrent) (string, error) {
	if torrent == nil {
		return "", nil
	}
	return torrent.InfoHash, nil
}

func (p *ruleMagnetTestProvider) GetTorrentMagnetLink(*hibiketorrent.MediaTorrent) (string, error) {
	p.calls++
	if p.err != nil {
		return "", p.err
	}
	return p.magnet, nil
}

func (p *ruleMagnetTestProvider) GetLatest() ([]*hibiketorrent.MediaTorrent, error) {
	return nil, nil
}

func (p *ruleMagnetTestProvider) GetSettings() hibiketorrent.AnimeProviderSettings {
	return hibiketorrent.AnimeProviderSettings{Type: hibiketorrent.AnimeProviderTypeMain}
}

func TestResolveAutoDownloaderItemMagnetUsesStoredTorrentExtension(t *testing.T) {
	provider := &ruleMagnetTestProvider{magnet: "magnet:?xt=urn:btih:resolved-from-provider"}
	repo := newTorrentRepositoryForRuleMagnetTests(map[string]*ruleMagnetTestProvider{"fake": provider})

	torrentData, err := json.Marshal(&autodownloader.NormalizedTorrent{
		MediaTorrent: &hibiketorrent.MediaTorrent{
			Name:     "Example torrent",
			InfoHash: "hash-from-torrent",
		},
		ExtensionID: "fake",
	})
	require.NoError(t, err)

	item := &models.AutoDownloaderItem{
		Hash:        "hash-from-item",
		TorrentData: torrentData,
	}

	magnet, err := resolveAutoDownloaderItemMagnet(item, repo)
	require.NoError(t, err)
	assert.Equal(t, "magnet:?xt=urn:btih:resolved-from-provider", magnet)
	assert.Equal(t, 1, provider.calls)
}

func TestResolveAutoDownloaderItemMagnetFallsBackToHash(t *testing.T) {
	provider := &ruleMagnetTestProvider{err: errors.New("provider failed")}
	repo := newTorrentRepositoryForRuleMagnetTests(map[string]*ruleMagnetTestProvider{"fake": provider})

	torrentData, err := json.Marshal(&autodownloader.NormalizedTorrent{
		MediaTorrent: &hibiketorrent.MediaTorrent{
			Name:     "Example torrent",
			InfoHash: "hash-from-torrent",
		},
		ExtensionID: "fake",
	})
	require.NoError(t, err)

	item := &models.AutoDownloaderItem{
		Hash:        "hash-from-item",
		TorrentData: torrentData,
	}

	magnet, err := resolveAutoDownloaderItemMagnet(item, repo)
	require.NoError(t, err)
	assert.Equal(t, "magnet:?xt=urn:btih:hash-from-item", magnet)
	assert.Equal(t, 1, provider.calls)
}

func newTorrentRepositoryForRuleMagnetTests(providers map[string]*ruleMagnetTestProvider) *torrentrepo.Repository {
	logger := zerolog.Nop()
	bank := extension.NewUnifiedBank()
	for id, provider := range providers {
		bank.Set(id, extension.NewMediaTorrentProviderExtension(&extension.Extension{
			ID:          id,
			Name:        id,
			Version:     "1.0.0",
			ManifestURI: "builtin",
			Language:    extension.LanguageGo,
			Type:        extension.TypeMediaTorrentProvider,
		}, provider))
	}

	var metadata metadata_provider.Provider
	repo := torrentrepo.NewRepository(&torrentrepo.NewRepositoryOptions{
		Logger:              &logger,
		MetadataProviderRef: util.NewRef[metadata_provider.Provider](metadata),
		ExtensionBankRef:    util.NewRef(bank),
	})
	repo.SetSettings(&torrentrepo.RepositorySettings{DefaultAnimeProvider: "fake"})

	return repo
}
