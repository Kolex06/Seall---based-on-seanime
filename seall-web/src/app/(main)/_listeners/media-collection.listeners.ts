import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { useWebsocketMessageListener } from "@/app/(main)/_hooks/handle-websockets"
import { WSEvents } from "@/lib/server/ws-events"
import { useQueryClient } from "@tanstack/react-query"

/**
 * @description
 * - Listens to REFRESHED_MEDIA_COLLECTION events and re-fetches queries associated with SIMKL collection.
 */
export function useAnimeCollectionListener() {

    const qc = useQueryClient()

    useWebsocketMessageListener({
        type: WSEvents.REFRESHED_MEDIA_COLLECTION,
        onMessage: data => {
            (async () => {
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetMediaCollectionSchedule.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LIBRARY_EXPLORER.GetLibraryExplorerFileTree.key] })
            })()
        },
    })

    useWebsocketMessageListener({
        type: WSEvents.REFRESHED_MANGA_COLLECTION,
        onMessage: data => {
            (async () => {
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
                await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key] })
            })()
        },
    })

}

