package onlinestream

import (
	"errors"
	"fmt"
	"seall/internal/api/mediaapi"
	"seall/internal/extension"
	hibikeonlinestream "seall/internal/extension/hibike/onlinestream"
	"seall/internal/util"
	"seall/internal/util/result"
	"strings"
)

var searchResultCache = result.NewCache[string, []*hibikeonlinestream.SearchResult]()

func (r *Repository) ManualSearch(provider string, query string, dub bool, media *mediaapi.BaseAnime) (ret []*hibikeonlinestream.SearchResult, err error) {
	defer util.HandlePanicInModuleWithError("onlinestream/ManualSearch", &err)

	if query == "" {
		return make([]*hibikeonlinestream.SearchResult, 0), nil
	}

	// Get the search results
	providerExtension, ok := extension.GetExtension[extension.OnlinestreamProviderExtension](r.extensionBankRef.Get(), provider)
	if !ok {
		r.logger.Error().Str("provider", provider).Msg("onlinestream: Provider not found")
		return nil, errors.New("onlinestream: Provider not found")
	}

	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	mediaType := MediaTypeForBaseMedia(media)
	settings := providerExtension.GetProvider().GetSettings()
	supportedMediaTypes := hibikeonlinestream.NormalizeSupportedMediaTypes(settings)
	if !hibikeonlinestream.SupportsMediaType(settings, mediaType) {
		return nil, unsupportedProviderMediaTypeError(provider, supportedMediaTypes, mediaType)
	}

	searchRes, found := searchResultCache.Get(provider + normalizedQuery + fmt.Sprintf("%t", dub) + mediaType)
	if found {
		return searchRes, nil
	}

	searchRes, err = providerExtension.GetProvider().Search(hibikeonlinestream.SearchOptions{
		Media:     queryMediaFromBaseMedia(media, mediaType),
		Query:     normalizedQuery,
		Dub:       dub,
		Year:      0,
		MediaType: mediaType,
	})
	if err != nil {
		r.logger.Error().Err(err).Str("query", normalizedQuery).Msg("onlinestream: Search failed")
		return nil, err
	}

	searchResultCache.Set(provider+normalizedQuery+fmt.Sprintf("%t", dub)+mediaType, searchRes)

	return searchRes, nil
}

// ManualMapping is used to manually map media to a provider result.
// After calling this, the client should re-fetch the episode list.
func (r *Repository) ManualMapping(provider string, mediaId int, providerMediaId string) (err error) {
	defer util.HandlePanicInModuleWithError("onlinestream/ManualMapping", &err)

	r.logger.Trace().Msgf("onlinestream: Removing cached bucket for %s, media ID: %d", provider, mediaId)

	// Delete the cached data if any
	epListBucket := r.getFcEpisodeListBucket(provider, mediaId)
	_ = r.fileCacher.Remove(epListBucket.Name())
	epDataBucket := r.getFcEpisodeDataBucket(provider, mediaId)
	_ = r.fileCacher.Remove(epDataBucket.Name())

	r.logger.Trace().
		Str("provider", provider).
		Int("mediaId", mediaId).
		Str("providerMediaId", providerMediaId).
		Msg("onlinestream: Manual mapping")

	// Insert the mapping into the database
	err = r.db.InsertOnlinestreamMapping(provider, mediaId, providerMediaId)
	if err != nil {
		r.logger.Error().Err(err).Msg("onlinestream: Failed to insert mapping")
		return err
	}

	r.logger.Debug().Msg("onlinestream: Manual mapping successful")

	return nil
}

type MappingResponse struct {
	ProviderMediaId *string `json:"providerMediaId"`
}

func (r *Repository) GetMapping(provider string, mediaId int) (ret MappingResponse) {
	defer util.HandlePanicInModuleThen("onlinestream/GetMapping", func() {
		ret = MappingResponse{}
	})

	mapping, found := r.db.GetOnlinestreamMapping(provider, mediaId)
	if !found {
		return MappingResponse{}
	}

	return MappingResponse{
		ProviderMediaId: &mapping.ProviderMediaID,
	}
}

func (r *Repository) RemoveMapping(provider string, mediaId int) (err error) {
	defer util.HandlePanicInModuleWithError("onlinestream/RemoveMapping", &err)

	// Delete the mapping from the database
	err = r.db.DeleteOnlinestreamMapping(provider, mediaId)
	if err != nil {
		r.logger.Error().Err(err).Msg("onlinestream: Failed to delete mapping")
		return err
	}

	r.logger.Debug().Msg("onlinestream: Mapping removed")

	r.logger.Trace().Msgf("onlinestream: Removing cached bucket for %s, media ID: %d", provider, mediaId)
	// Delete the cached data if any
	epListBucket := r.getFcEpisodeListBucket(provider, mediaId)
	_ = r.fileCacher.Remove(epListBucket.Name())
	epDataBucket := r.getFcEpisodeDataBucket(provider, mediaId)
	_ = r.fileCacher.Remove(epDataBucket.Name())

	return nil
}
