import { useServerMutation, useServerQuery } from "@/api/client/requests"
import {
    MediaEntryBulkAction_Variables,
    MediaEntryManualMatch_Variables,
    FetchMediaEntrySuggestions_Variables,
    OpenMediaEntryInExplorer_Variables,
    ToggleMediaEntrySilenceStatus_Variables,
    UpdateMediaEntryProgress_Variables,
    UpdateMediaEntryRepeat_Variables,
} from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { MediaAPI_BaseMedia, Media_Entry, Media_LocalFile, Media_MissingEpisodes, Media_UpcomingEpisodes, Nullish } from "@/api/generated/types"
import { getEntryPreloadStaleTime } from "@/lib/entry-preloader"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export function useGetMediaEntry(id: Nullish<string | number>) {
    return useServerQuery<Media_Entry>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.endpoint.replace("{id}", String(id)),
        method: API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key, String(id)],
        enabled: !!id,
        staleTime: getEntryPreloadStaleTime("anime", id),
    })
}

export function useMediaEntryBulkAction(id?: Nullish<number>, onSuccess?: () => void) {
    const queryClient = useQueryClient()

    return useServerMutation<Array<Media_LocalFile>, MediaEntryBulkAction_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.MediaEntryBulkAction.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.MediaEntryBulkAction.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.MediaEntryBulkAction.key, String(id)],
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key, String(id)] })
            queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.LIBRARY_EXPLORER.GetLibraryExplorerFileTree.key] })
            onSuccess?.()
        },
    })
}

export function useOpenMediaEntryInExplorer() {
    return useServerMutation<boolean, OpenMediaEntryInExplorer_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.OpenMediaEntryInExplorer.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.OpenMediaEntryInExplorer.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.OpenMediaEntryInExplorer.key],
        onSuccess: async () => {

        },
    })
}

export function useFetchMediaEntrySuggestions() {
    return useServerMutation<Array<MediaAPI_BaseMedia>, FetchMediaEntrySuggestions_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.FetchMediaEntrySuggestions.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.FetchMediaEntrySuggestions.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.FetchMediaEntrySuggestions.key],
        onSuccess: async () => {

        },
    })
}

export function useMediaEntryManualMatch() {
    const queryClient = useQueryClient()

    return useServerMutation<Array<Media_LocalFile>, MediaEntryManualMatch_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.MediaEntryManualMatch.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.MediaEntryManualMatch.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.MediaEntryManualMatch.key],
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
            queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.LIBRARY_EXPLORER.GetLibraryExplorerFileTree.key] })
            toast.success("Files matched")
        },
    })
}

export function useGetMissingEpisodes(enabled?: boolean) {
    return useServerQuery<Media_MissingEpisodes>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key],
        enabled: enabled ?? true, // Default to true if not provided
    })
}

export function useGetMediaEntrySilenceStatus(id: Nullish<string | number>) {
    const { data, ...rest } = useServerQuery({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntrySilenceStatus.endpoint.replace("{id}", String(id)),
        method: API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntrySilenceStatus.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntrySilenceStatus.key],
        enabled: !!id,
    })

    return { isSilenced: !!data, ...rest }
}

export function useToggleMediaEntrySilenceStatus() {
    const queryClient = useQueryClient()

    return useServerMutation<boolean, ToggleMediaEntrySilenceStatus_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.ToggleMediaEntrySilenceStatus.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.ToggleMediaEntrySilenceStatus.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.ToggleMediaEntrySilenceStatus.key],
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntrySilenceStatus.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key] })
        },
    })
}

export function useUpdateMediaEntryProgress(id: Nullish<string | number>, episodeNumber: number, showToast: boolean = true) {
    const queryClient = useQueryClient()

    return useServerMutation<boolean, UpdateMediaEntryProgress_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.UpdateMediaEntryProgress.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.UpdateMediaEntryProgress.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.UpdateMediaEntryProgress.key, id, episodeNumber],
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            if (id) {
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key, String(id)] })
            }
            if (showToast) {
                toast.success("Progress updated successfully")
            }
        },
    })
}

export function useUpdateMediaEntryRepeat(id: Nullish<string | number>) {
    const queryClient = useQueryClient()

    return useServerMutation<boolean, UpdateMediaEntryRepeat_Variables>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.UpdateMediaEntryRepeat.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.UpdateMediaEntryRepeat.methods[0],
        mutationKey: [API_ENDPOINTS.MEDIA_ENTRIES.UpdateMediaEntryRepeat.key, id],
        onSuccess: async () => {
            // if (id) {
            //     await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key, String(id)] })
            // }
            // toast.success("Updated successfully")
        },
    })
}

export function useGetUpcomingEpisodes() {
    return useServerQuery<Media_UpcomingEpisodes>({
        endpoint: API_ENDPOINTS.MEDIA_ENTRIES.GetUpcomingEpisodes.endpoint,
        method: API_ENDPOINTS.MEDIA_ENTRIES.GetUpcomingEpisodes.methods[0],
        queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetUpcomingEpisodes.key],
        enabled: true,
    })
}
