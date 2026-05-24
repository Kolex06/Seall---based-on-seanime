import { useServerMutation, useServerQuery } from "@/api/client/requests"
import {
    ListMedia_Variables,
    ListRecentAiringMedia_Variables,
    DeleteMediaListEntry_Variables,
    EditMediaListEntry_Variables,
} from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import {
    MediaAPI_MediaCollection,
    MediaAPI_MediaDetailsById_Media,
    MediaAPI_BaseMedia,
    MediaAPI_ListMedia,
    MediaAPI_ListRecentMedia,
    MediaAPI_Stats,
    MediaAPI_StudioDetails,
    Nullish,
} from "@/api/generated/types"
import { getEntryPreloadStaleTime } from "@/lib/entry-preloader"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export function useGetMediaCollection() {
    return useServerQuery<MediaAPI_MediaCollection>({
        endpoint: API_ENDPOINTS.SIMKL.GetMediaCollection.endpoint,
        method: API_ENDPOINTS.SIMKL.GetMediaCollection.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key],
        enabled: true,
    })
}

export function useGetRawMediaCollection() {
    return useServerQuery<MediaAPI_MediaCollection>({
        endpoint: API_ENDPOINTS.SIMKL.GetRawMediaCollection.endpoint,
        method: API_ENDPOINTS.SIMKL.GetRawMediaCollection.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key],
        enabled: true,
    })
}

export function useGetRawMediaCollectionTags() {
    return useServerQuery<Record<number, Array<string>>>({
        endpoint: API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.endpoint,
        method: API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key],
        enabled: true,
    })
}

export function useRefreshMediaCollection() {
    const queryClient = useQueryClient()

    return useServerMutation<MediaAPI_MediaCollection>({
        endpoint: API_ENDPOINTS.SIMKL.GetMediaCollection.endpoint,
        method: API_ENDPOINTS.SIMKL.GetMediaCollection.methods[1],
        mutationKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key],
        onSuccess: async () => {
            toast.success("SIMKL data is up-to-date")
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollectionTags.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetMediaCollectionSchedule.key] })
        },
    })
}

export function useEditMediaListEntry(id: Nullish<string | number>, type: "anime" | "manga") {
    const queryClient = useQueryClient()

    return useServerMutation<boolean, EditMediaListEntry_Variables>({
        endpoint: API_ENDPOINTS.SIMKL.EditMediaListEntry.endpoint,
        method: API_ENDPOINTS.SIMKL.EditMediaListEntry.methods[0],
        mutationKey: [API_ENDPOINTS.SIMKL.EditMediaListEntry.key, String(id)],
        onSuccess: async () => {
            toast.success("Entry updated")
            if (type === "anime") {
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key, String(id)] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key] })
            } else if (type === "manga") {
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key, String(id)] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollectionTags.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
            }
        },
    })
}

export function useGetMediaDetails(id: Nullish<number | string>) {
    return useServerQuery<MediaAPI_MediaDetailsById_Media>({
        endpoint: API_ENDPOINTS.SIMKL.GetMediaDetails.endpoint.replace("{id}", String(id)),
        method: API_ENDPOINTS.SIMKL.GetMediaDetails.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetMediaDetails.key, String(id)],
        enabled: !!id,
        staleTime: getEntryPreloadStaleTime("anime", id),
    })
}

export function useDeleteMediaListEntry(id: Nullish<string | number>, type: "anime" | "manga", onSuccess: () => void) {
    const queryClient = useQueryClient()

    return useServerMutation<boolean, DeleteMediaListEntry_Variables>({
        endpoint: API_ENDPOINTS.SIMKL.DeleteMediaListEntry.endpoint,
        method: API_ENDPOINTS.SIMKL.DeleteMediaListEntry.methods[0],
        mutationKey: [API_ENDPOINTS.SIMKL.DeleteMediaListEntry.key],
        onSuccess: async () => {
            toast.success("Entry deleted")
            if (type === "anime") {
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key, String(id)] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key] })
            } else if (type === "manga") {
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key, String(id)] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollectionTags.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
            }
            onSuccess()
        },
    })
}

export function useListMedia(variables: ListMedia_Variables, enabled: boolean) {
    return useServerQuery<MediaAPI_ListMedia, ListMedia_Variables>({
        endpoint: API_ENDPOINTS.SIMKL.ListMedia.endpoint,
        method: API_ENDPOINTS.SIMKL.ListMedia.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.ListMedia.key, variables],
        data: variables,
        enabled: enabled ?? true,
    })
}

export function useListRecentAiringMedia(variables: ListRecentAiringMedia_Variables, enabled: boolean = true) {
    return useServerQuery<MediaAPI_ListRecentMedia, ListRecentAiringMedia_Variables>({
        endpoint: API_ENDPOINTS.SIMKL.ListRecentAiringMedia.endpoint,
        method: API_ENDPOINTS.SIMKL.ListRecentAiringMedia.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.ListRecentAiringMedia.key, JSON.stringify(variables)],
        data: variables,
        enabled: enabled,
    })
}

export function useGetStudioDetails(id: number) {
    return useServerQuery<MediaAPI_StudioDetails>({
        endpoint: API_ENDPOINTS.SIMKL.GetStudioDetails.endpoint.replace("{id}", String(id)),
        method: API_ENDPOINTS.SIMKL.GetStudioDetails.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetStudioDetails.key, String(id)],
        enabled: true,
    })
}

export function useGetMediaStats(enabled: boolean = true) {
    return useServerQuery<MediaAPI_Stats>({
        endpoint: API_ENDPOINTS.SIMKL.GetMediaStats.endpoint,
        method: API_ENDPOINTS.SIMKL.GetMediaStats.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetMediaStats.key],
        enabled: enabled,
    })
}

export function useListMissedSequels(enabled: boolean) {
    return useServerQuery<Array<MediaAPI_BaseMedia>>({
        endpoint: API_ENDPOINTS.SIMKL.ListMissedSequels.endpoint,
        method: API_ENDPOINTS.SIMKL.ListMissedSequels.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.ListMissedSequels.key],
        enabled: enabled,
    })
}

export function useGetMediaCacheLayerStatus() {
    return useServerQuery<boolean>({
        endpoint: API_ENDPOINTS.SIMKL.GetMediaCacheLayerStatus.endpoint,
        method: API_ENDPOINTS.SIMKL.GetMediaCacheLayerStatus.methods[0],
        queryKey: [API_ENDPOINTS.SIMKL.GetMediaCacheLayerStatus.key],
        gcTime: 0,
        enabled: true,
    })
}

export function useToggleMediaCacheLayerStatus() {
    const queryClient = useQueryClient()
    return useServerMutation<boolean>({
        endpoint: API_ENDPOINTS.SIMKL.ToggleMediaCacheLayerStatus.endpoint,
        method: API_ENDPOINTS.SIMKL.ToggleMediaCacheLayerStatus.methods[0],
        mutationKey: [API_ENDPOINTS.SIMKL.ToggleMediaCacheLayerStatus.key],
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCacheLayerStatus.key] })
        },
    })
}
