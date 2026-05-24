import { useServerMutation, useServerQuery } from "@/api/client/requests"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { Local_QueueState, Local_TrackedMediaItem } from "@/api/generated/types"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import {
    LocalAddTrackedMedia_Variables,
    LocalRemoveTrackedMedia_Variables,
    LocalSetHasLocalChanges_Variables,
    SetOfflineMode_Variables,
} from "../generated/endpoint.types"

export function useLocalGetTrackedMediaItems() {
    return useServerQuery<Array<Local_TrackedMediaItem>>({
        endpoint: API_ENDPOINTS.LOCAL.LocalGetTrackedMediaItems.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalGetTrackedMediaItems.methods[0],
        queryKey: [API_ENDPOINTS.LOCAL.LocalGetTrackedMediaItems.key],
        enabled: true,
        gcTime: 0,
    })
}

export function useLocalAddTrackedMedia() {
    const qc = useQueryClient()
    return useServerMutation<boolean, LocalAddTrackedMedia_Variables>({
        endpoint: API_ENDPOINTS.LOCAL.LocalAddTrackedMedia.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalAddTrackedMedia.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.LocalAddTrackedMedia.key],
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetTrackedMediaItems.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetSyncQueueState.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetIsMediaTracked.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize] })
            toast.success("Added media for offline syncing")
        },
    })
}

export function useLocalRemoveTrackedMedia() {
    const qc = useQueryClient()
    return useServerMutation<boolean, LocalRemoveTrackedMedia_Variables>({
        endpoint: API_ENDPOINTS.LOCAL.LocalRemoveTrackedMedia.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalRemoveTrackedMedia.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.LocalRemoveTrackedMedia.key],
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetTrackedMediaItems.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetSyncQueueState.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetIsMediaTracked.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize] })
            toast.success("Removed offline data")
        },
    })
}

export function useLocalSyncData() {
    return useServerMutation<boolean>({
        endpoint: API_ENDPOINTS.LOCAL.LocalSyncData.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalSyncData.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.LocalSyncData.key],
        onSuccess: async () => {
            toast.info("Syncing local data...")
        },
    })
}

export function useLocalGetSyncQueueData() {
    return useServerQuery<Local_QueueState>({
        endpoint: API_ENDPOINTS.LOCAL.LocalGetSyncQueueState.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalGetSyncQueueState.methods[0],
        queryKey: [API_ENDPOINTS.LOCAL.LocalGetSyncQueueState.key],
        enabled: true,
    })
}

export function useLocalGetIsMediaTracked(id: number, type: string) {
    return useServerQuery<boolean>({
        endpoint: API_ENDPOINTS.LOCAL.LocalGetIsMediaTracked.endpoint.replace("{id}", String(id)).replace("{type}", String(type)),
        method: API_ENDPOINTS.LOCAL.LocalGetIsMediaTracked.methods[0],
        queryKey: [API_ENDPOINTS.LOCAL.LocalGetIsMediaTracked.key, id, type],
        enabled: true,
    })
}

export function useLocalSyncMediaData() {
    const qc = useQueryClient()
    return useServerMutation<boolean>({
        endpoint: API_ENDPOINTS.LOCAL.LocalSyncMediaData.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalSyncMediaData.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.LocalSyncMediaData.key],
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetTrackedMediaItems.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollectionTags.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes] })
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize] })
            toast.success("Updated Simkl data")
        },
    })
}

export function useLocalSetHasLocalChanges() {
    const qc = useQueryClient()
    return useServerMutation<boolean, LocalSetHasLocalChanges_Variables>({
        endpoint: API_ENDPOINTS.LOCAL.LocalSetHasLocalChanges.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalSetHasLocalChanges.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.LocalSetHasLocalChanges.key],
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetHasLocalChanges.key] })
        },
    })
}

export function useLocalGetHasLocalChanges() {
    return useServerQuery<boolean>({
        endpoint: API_ENDPOINTS.LOCAL.LocalGetHasLocalChanges.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalGetHasLocalChanges.methods[0],
        queryKey: [API_ENDPOINTS.LOCAL.LocalGetHasLocalChanges.key],
        enabled: true,
    })
}

export function useLocalGetLocalStorageSize() {
    return useServerQuery<string>({
        endpoint: API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize.methods[0],
        queryKey: [API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize.key],
        enabled: true,
    })
}

export function useLocalSyncSimulatedDataToMedia() {
    const qc = useQueryClient()
    return useServerMutation<boolean>({
        endpoint: API_ENDPOINTS.LOCAL.LocalSyncSimulatedDataToMedia.endpoint,
        method: API_ENDPOINTS.LOCAL.LocalSyncSimulatedDataToMedia.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.LocalSyncSimulatedDataToMedia.key],
        onSuccess: async () => {
            ({ queryKey: [API_ENDPOINTS.LOCAL.LocalGetLocalStorageSize] })
            toast.success("Updated Simkl data")
        },
    })
}

export function useSetOfflineMode() {
    return useServerMutation<boolean, SetOfflineMode_Variables>({
        endpoint: API_ENDPOINTS.LOCAL.SetOfflineMode.endpoint,
        method: API_ENDPOINTS.LOCAL.SetOfflineMode.methods[0],
        mutationKey: [API_ENDPOINTS.LOCAL.SetOfflineMode.key],
        onSuccess: async (data) => {
            if (data) {
                toast.success("Offline mode enabled")
                window.location.href = "/offline"
            } else {
                toast.success("Offline mode disabled")
                window.location.href = "/"
            }
        },
    })
}
