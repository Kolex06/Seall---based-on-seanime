package simkl

import "sync"

var discoveryMediaCache sync.Map

func cacheDiscoveryMedia(kind MediaType, media StandardMedia) {
	kind = KindFromStandardMedia(kind, &media)
	id := media.IDs.PrimarySimklID()
	if id == 0 {
		id = simklIDFromURL(media.URL)
	}
	if id == 0 {
		id = stableFallbackID(kind, media.Title, media.Year)
	}
	if id == 0 {
		return
	}
	if media.IDs.Simkl == 0 {
		media.IDs.Simkl = id
	}
	if cached, ok := CachedDiscoveryMedia(id); ok && cached.Kind == kind {
		if media.EnglishName == "" {
			media.EnglishName = cached.Media.EnglishName
		}
		if media.TitleEnglish == "" {
			media.TitleEnglish = cached.Media.TitleEnglish
		}
		if media.Relations == nil {
			media.Relations = cached.Media.Relations
		}
	}
	discoveryMediaCache.Store(id, DiscoveryMedia{
		Kind:  kind,
		Media: media,
	})
}

func CachedDiscoveryMedia(id int) (DiscoveryMedia, bool) {
	if item, ok := discoveryMediaCache.Load(id); ok {
		if media, ok := item.(DiscoveryMedia); ok {
			return media, true
		}
	}
	return DiscoveryMedia{}, false
}
