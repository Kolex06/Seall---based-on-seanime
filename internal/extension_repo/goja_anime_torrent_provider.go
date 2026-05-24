package extension_repo

import (
	"context"
	"seall/internal/events"
	"seall/internal/extension"
	hibiketorrent "seall/internal/extension/hibike/torrent"
	"seall/internal/goja/goja_runtime"
	"seall/internal/util"

	"github.com/rs/zerolog"
)

type GojaMediaTorrentProvider struct {
	*gojaProviderBase
}

func NewGojaMediaTorrentProvider(ext *extension.Extension, language extension.Language, logger *zerolog.Logger, runtimeManager *goja_runtime.Manager, wsEventManager events.WSEventManagerInterface) (hibiketorrent.AnimeProvider, *GojaMediaTorrentProvider, error) {
	base, err := initializeProviderBase(ext, language, logger, runtimeManager, wsEventManager)
	if err != nil {
		return nil, nil, err
	}

	provider := &GojaMediaTorrentProvider{
		gojaProviderBase: base,
	}
	return provider, provider, nil
}

func (g *GojaMediaTorrentProvider) Search(opts hibiketorrent.AnimeSearchOptions) (ret []*hibiketorrent.MediaTorrent, err error) {
	defer util.HandlePanicInModuleWithError(g.ext.ID+".Search", &err)

	method, err := g.callClassMethod(context.Background(), "search", structToMap(opts))

	promiseRes, err := g.waitForPromise(method)
	if err != nil {
		return nil, err
	}

	err = g.unmarshalValue(promiseRes, &ret)
	if err != nil {
		return nil, err
	}

	for i := range ret {
		ret[i].Provider = g.ext.ID
	}

	return
}

func (g *GojaMediaTorrentProvider) SmartSearch(opts hibiketorrent.AnimeSmartSearchOptions) (ret []*hibiketorrent.MediaTorrent, err error) {
	defer util.HandlePanicInModuleWithError(g.ext.ID+".SmartSearch", &err)

	method, err := g.callClassMethod(context.Background(), "smartSearch", structToMap(opts))

	promiseRes, err := g.waitForPromise(method)
	if err != nil {
		return nil, err
	}

	err = g.unmarshalValue(promiseRes, &ret)
	if err != nil {
		return nil, err
	}

	for i := range ret {
		ret[i].Provider = g.ext.ID
	}

	return
}

func (g *GojaMediaTorrentProvider) GetTorrentInfoHash(torrent *hibiketorrent.MediaTorrent) (ret string, err error) {
	defer util.HandlePanicInModuleWithError(g.ext.ID+".GetTorrentInfoHash", &err)

	res, err := g.callClassMethod(context.Background(), "getTorrentInfoHash", structToMap(torrent))
	if err != nil {
		return "", err
	}

	promiseRes, err := g.waitForPromise(res)
	if err != nil {
		return "", err
	}

	err = g.unmarshalValue(promiseRes, &ret)
	if err != nil {
		return "", err
	}

	return
}

func (g *GojaMediaTorrentProvider) GetTorrentMagnetLink(torrent *hibiketorrent.MediaTorrent) (ret string, err error) {
	defer util.HandlePanicInModuleWithError(g.ext.ID+".GetTorrentMagnetLink", &err)

	res, err := g.callClassMethod(context.Background(), "getTorrentMagnetLink", structToMap(torrent))
	if err != nil {
		return "", err
	}

	promiseRes, err := g.waitForPromise(res)
	if err != nil {
		return "", err
	}

	err = g.unmarshalValue(promiseRes, &ret)
	if err != nil {
		return "", err
	}

	return
}

func (g *GojaMediaTorrentProvider) GetLatest() (ret []*hibiketorrent.MediaTorrent, err error) {
	defer util.HandlePanicInModuleWithError(g.ext.ID+".GetLatest", &err)

	method, err := g.callClassMethod(context.Background(), "getLatest")
	if err != nil {
		return nil, err
	}

	promiseRes, err := g.waitForPromise(method)
	if err != nil {
		return nil, err
	}

	err = g.unmarshalValue(promiseRes, &ret)
	if err != nil {
		return nil, err
	}

	return
}

func (g *GojaMediaTorrentProvider) GetSettings() (ret hibiketorrent.AnimeProviderSettings) {
	defer util.HandlePanicInModuleThen(g.ext.ID+".GetSettings", func() {
		ret = hibiketorrent.AnimeProviderSettings{}
	})

	res, err := g.callClassMethod(context.Background(), "getSettings")
	if err != nil {
		return
	}

	err = g.unmarshalValue(res, &ret)
	if err != nil {
		return
	}

	return
}
