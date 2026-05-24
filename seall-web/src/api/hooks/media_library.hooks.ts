import { useServerMutation, useServerQuery } from "@/api/client/requests"
import { AddUnknownMedia_Variables } from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { MediaAPI_MediaCollection, Media_LibraryCollection, Media_ScheduleItem } from "@/api/generated/types"
import { useRefreshMediaCollection } from "@/api/hooks/simkl.hooks"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export function useGetLibraryCollection({ enabled }: { enabled?: boolean } = { enabled: true }) {
    return useServerQuery<Media_LibraryCollection>({
        endpoint: API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.endpoint,
        method: API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key],
        enabled: enabled,
    })
}

export function useAddUnknownMedia() {
    const queryClient = useQueryClient()
    const { mutate } = useRefreshMediaCollection()

    return useServerMutation<MediaAPI_MediaCollection, AddUnknownMedia_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_LIBRARY.AddUnknownMedia.endpoint,
        method: API_ENDPOINTS.MEDIA_LIBRARY.AddUnknownMedia.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_LIBRARY.AddUnknownMedia.key],
        onSuccess: async () => {
            toast.success("Media added successfully")
            mutate(undefined, {
                onSuccess: () => {
                    queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.LIBRARY_EXPLORER.GetLibraryExplorerFileTree.key] })
                },
            })
        },
    })
}

export function useGetMediaCollectionSchedule({ enabled, source = "list" }: { enabled?: boolean, source?: "list" | "all" } = { enabled: true, source: "list" }) {
    const params = new URLSearchParams({ source })
    return useServerQuery<Array<Media_ScheduleItem>>({
        endpoint: `${API_ENDPOINTS.MEDIA_LIBRARY.GetMediaCollectionSchedule.endpoint}?${params.toString()}`,
        method: API_ENDPOINTS.MEDIA_LIBRARY.GetMediaCollectionSchedule.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetMediaCollectionSchedule.key, source],
        enabled: enabled,
    })
}
