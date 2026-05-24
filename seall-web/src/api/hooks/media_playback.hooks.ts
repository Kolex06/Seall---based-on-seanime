import { useServerQuery } from "@/api/client/requests"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { Media_EpisodeCollection, Nullish } from "@/api/generated/types"

export function useGetMediaEpisodeCollection(id: Nullish<number>) {
    return useServerQuery<Media_EpisodeCollection>({
        endpoint: API_ENDPOINTS.MEDIA_PLAYBACK.GetMediaEpisodeCollection.endpoint.replace("{id}", String(id)),
        method: API_ENDPOINTS.MEDIA_PLAYBACK.GetMediaEpisodeCollection.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_PLAYBACK.GetMediaEpisodeCollection.key, String(id)],
        enabled: !!id,
    })
}
