import { useServerMutation } from "@/api/client/requests"
import { Login_Variables } from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { Status } from "@/api/generated/types"
import { useSetServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { useRouter } from "@/lib/navigation"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export type SimklPinCode = {
    result?: string
    device_code?: string
    user_code?: string
    verification_url?: string
    expires_in?: number
    interval?: number
}

export type SimklPinStatus = {
    result?: string
    message?: string
    access_token?: string
}

export function useLogin() {
    const queryClient = useQueryClient()
    const router = useRouter()
    const setServerStatus = useSetServerStatus()

    return useServerMutation<Status, Login_Variables>({
        endpoint: API_ENDPOINTS.AUTH.Login.endpoint,
        method: API_ENDPOINTS.AUTH.Login.methods[0],
        mutationKey: [API_ENDPOINTS.AUTH.Login.key],
        onSuccess: async data => {
            if (data) {
                toast.success("Successfully authenticated")
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollectionTags.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
                setServerStatus(data)
                router.push("/")
                queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key] })
                queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
                queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key] })
            }
        },
        onError: async error => {
            toast.error(error.message)
            router.push("/")
        },
    })
}

export function useStartSimklPinLogin() {
    return useServerMutation<SimklPinCode>({
        endpoint: "/api/v1/auth/simkl/pin",
        method: "POST",
        mutationKey: ["AUTH-start-simkl-pin-login"],
    })
}

export function useCheckSimklPinLogin() {
    const queryClient = useQueryClient()
    const router = useRouter()
    const setServerStatus = useSetServerStatus()

    return useServerMutation<Status | SimklPinStatus, { userCode: string }>({
        endpoint: "/api/v1/auth/simkl/pin/check",
        method: "POST",
        mutationKey: ["AUTH-check-simkl-pin-login"],
        onSuccess: async data => {
            if (data && "user" in data) {
                toast.success("Successfully authenticated with SIMKL")
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
                await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
                setServerStatus(data as Status)
                router.push("/")
            }
        },
    })
}

export function useSaveSimklClientConfig() {
    const queryClient = useQueryClient()
    const setServerStatus = useSetServerStatus()

    return useServerMutation<Status, { clientId: string }>({
        endpoint: "/api/v1/auth/simkl/client",
        method: "PATCH",
        mutationKey: ["AUTH-save-simkl-client-config"],
        onSuccess: async data => {
            if (data) {
                setServerStatus(data)
            }
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.STATUS.GetStatus.key] })
            toast.success("SIMKL client ID saved")
        },
    })
}

export function useLogout() {
    const queryClient = useQueryClient()
    const router = useRouter()
    const setServerStatus = useSetServerStatus()

    return useServerMutation<Status>({
        endpoint: API_ENDPOINTS.AUTH.Logout.endpoint,
        method: API_ENDPOINTS.AUTH.Logout.methods[0],
        mutationKey: [API_ENDPOINTS.AUTH.Logout.key],
        onSuccess: async data => {
            if (data) {
                setServerStatus(data)
            }
            toast.success("Successfully logged out")
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_LIBRARY.GetLibraryCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetRawMediaCollectionTags.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.SIMKL.GetMediaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollection.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetRawMangaCollectionTags.key] })
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaCollection.key] })
            router.push("/")
            queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMissingEpisodes.key] })
            queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MEDIA_ENTRIES.GetMediaEntry.key] })
            queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.MANGA.GetMangaEntry.key] })
        },
    })
}
