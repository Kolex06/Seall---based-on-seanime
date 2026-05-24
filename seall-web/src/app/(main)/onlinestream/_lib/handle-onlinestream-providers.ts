import { useListOnlinestreamProviderExtensions } from "@/api/hooks/extensions.hooks"
import { ExtensionRepo_OnlinestreamProviderExtensionItem } from "@/api/generated/types"
import { __onlinestream_selectedProviderAtom } from "@/app/(main)/onlinestream/_lib/onlinestream.atoms"
import { logger } from "@/lib/helpers/debug"
import { useAtom } from "jotai/react"
import React from "react"

function normalizeMediaType(mediaType: string | null | undefined) {
    if (!mediaType) return "anime"
    const value = mediaType.trim().toLowerCase()
    if (value === "movie") return "movies"
    if (["show", "tv", "series"].includes(value)) return "shows"
    return value
}

function supportsMediaType(extension: ExtensionRepo_OnlinestreamProviderExtensionItem, mediaType: string | null | undefined) {
    const targetType = normalizeMediaType(mediaType)
    const supported = extension.supportedMediaTypes?.length ? extension.supportedMediaTypes.map(normalizeMediaType) : ["anime"]
    return supported.includes("all") || supported.includes(targetType)
}

export function useHandleOnlinestreamProviderExtensions(mediaType?: string | null) {

    const { data: providerExtensions } = useListOnlinestreamProviderExtensions()

    const [provider, setProvider] = useAtom(__onlinestream_selectedProviderAtom)
    const filteredProviderExtensions = React.useMemo(() => {
        return (providerExtensions ?? []).filter(extension => supportsMediaType(extension, mediaType))
    }, [providerExtensions, mediaType])

    /**
     * Override the selected provider if it is not available
     */
    React.useLayoutEffect(() => {
        logger("ONLINESTREAM").info("extensions", filteredProviderExtensions)

        if (!providerExtensions) return

        if (provider === null || !filteredProviderExtensions.find(p => p.id === provider)) {
            if (filteredProviderExtensions.length > 0) {
                setProvider(filteredProviderExtensions[0].id)
            } else {
                setProvider(null)
            }
        }
    }, [providerExtensions, filteredProviderExtensions, provider])

    return {
        providerExtensions: filteredProviderExtensions,
        allProviderExtensions: providerExtensions ?? [],
        providerExtensionOptions: filteredProviderExtensions.map(provider => ({
            label: provider.name,
            value: provider.id,
        })).sort((a, b) => a.label.localeCompare(b.label)),
    }

}
