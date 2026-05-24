import { useServerMutation } from "@/api/client/requests"
import { ScanLocalFiles_Variables } from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { useRefreshLibraryExplorerFileTree } from "@/api/generated/library_explorer.hooks"
import { Media_LocalFile } from "@/api/generated/types"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export function useScanLocalFiles(onSuccess?: () => void) {
    const queryClient = useQueryClient()
    const { mutate: refreshLibraryExplorerTree } = useRefreshLibraryExplorerFileTree()

    return useServerMutation<Array<Media_LocalFile>, ScanLocalFiles_Variables>({
        endpoint: API_ENDPOINTS.SCAN.ScanLocalFiles.endpoint,
        method: API_ENDPOINTS.SCAN.ScanLocalFiles.methods[0],
        mutationKey: [API_ENDPOINTS.SCAN.ScanLocalFiles.key],
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            toast.success("Library scanned")
            refreshLibraryExplorerTree()
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.AUTO_DOWNLOADER.GetAutoDownloaderItems.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
            onSuccess?.()
        },
    })
}


