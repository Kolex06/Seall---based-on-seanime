import { Media_Entry, Media_Episode, Media_EpisodeCollection } from "@/api/generated/types"
import { useWebsocketMessageListener } from "@/app/(main)/_hooks/handle-websockets"
import { TorrentStreamEpisodeSection } from "@/app/(main)/entry/_containers/torrent-stream/_components/torrent-stream-episode-section"
import { cn } from "@/components/ui/core/styling"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { WSEvents } from "@/lib/server/ws-events"
import { atom } from "jotai"
import { useAtomValue } from "jotai/react"
import { useAtom } from "jotai/react"
import React, { startTransition } from "react"
import { BiExtension } from "react-icons/bi"
import {
    usePluginListenMediaEntryEpisodeTabEpisodeCollectionEvent,
    usePluginListenMediaEntryEpisodeTabsUpdatedEvent,
    usePluginSendMediaEntryEpisodeTabOpenEvent,
    usePluginSendMediaEntryEpisodeTabSelectEpisodeEvent,
    usePluginSendMediaEntryEpisodeTabsRenderEvent,
    usePluginSendMediaEntryEpisodeTabStateChangedEvent,
} from "./generated/plugin-events"

function sortTabs(tabs: PluginMediaEntryEpisodeTab[]) {
    return tabs.sort((a, b) => {
        const extCmp = a.extensionId.localeCompare(b.extensionId, undefined, { numeric: true })
        if (extCmp !== 0) return extCmp
        return a.name.localeCompare(b.name, undefined, { numeric: true })
    })
}

export function getPluginEpisodeTabViewId(extensionId: string) {
    return `episodeTab:${extensionId}`
}

export function getPluginEpisodeTabExtensionId(viewId: string) {
    return viewId.startsWith("episodeTab:") ? viewId.slice("episodeTab:".length) : ""
}

const __plugin_episodeTabsAtom = atom<PluginMediaEntryEpisodeTab[]>([])
const __plugin_episodeTabCollectionsAtom = atom<Record<string, Media_EpisodeCollection | undefined>>({})
const __plugin_episodeTabRenderedExtensionIdsAtom = atom<string[]>([])

export type PluginMediaEntryEpisodeTab = {
    extensionId: string
    name: string
    icon?: string
    viewId: string
}

export function usePluginMediaEntryEpisodeTabsListener(props: {
    mediaId: number
    currentView: string
    setView: (view: string) => void
}) {
    const { mediaId, currentView, setView } = props

    const [tabs, setTabs] = useAtom(__plugin_episodeTabsAtom)
    const [, setCollections] = useAtom(__plugin_episodeTabCollectionsAtom)
    const [renderedExtensionIds, setRenderedExtensionIds] = useAtom(__plugin_episodeTabRenderedExtensionIdsAtom)

    const { sendMediaEntryEpisodeTabsRenderEvent } = usePluginSendMediaEntryEpisodeTabsRenderEvent()
    const { sendMediaEntryEpisodeTabOpenEvent } = usePluginSendMediaEntryEpisodeTabOpenEvent()
    const { sendMediaEntryEpisodeTabStateChangedEvent } = usePluginSendMediaEntryEpisodeTabStateChangedEvent()

    const renderTabs = React.useEffectEvent(() => {
        if (!mediaId) return
        setRenderedExtensionIds([])
        sendMediaEntryEpisodeTabsRenderEvent({ mediaId }, "")
    })

    React.useEffect(() => {
        setCollections({})
        renderTabs()
    }, [mediaId])

    useWebsocketMessageListener({
        type: WSEvents.PLUGIN_LOADED,
        onMessage: () => {
            renderTabs()
        },
    })

    useWebsocketMessageListener({
        type: WSEvents.PLUGIN_UNLOADED,
        onMessage: (extensionId: string) => {
            startTransition(() => {
                setTabs(prev => prev.filter(tab => tab.extensionId !== extensionId))
                setRenderedExtensionIds(prev => prev.filter(id => id !== extensionId))
                setCollections(prev => {
                    const next = { ...prev }
                    Object.keys(next).forEach(key => {
                        if (key === getPluginEpisodeTabViewId(extensionId)) {
                            delete next[key]
                        }
                    })
                    return next
                })
                if (currentView === getPluginEpisodeTabViewId(extensionId)) {
                    setView("library")
                }
            })
        },
    })

    usePluginListenMediaEntryEpisodeTabsUpdatedEvent((event, extensionId) => {
        startTransition(() => {
            setRenderedExtensionIds(prev => prev.includes(extensionId) ? prev : [...prev, extensionId])
            setTabs(prev => {
                const otherTabs = prev.filter(tab => tab.extensionId !== extensionId)
                const extensionTabs = (event.tabs ?? []).map((tab: Record<string, any>) => ({
                    ...tab,
                    extensionId,
                    viewId: getPluginEpisodeTabViewId(extensionId),
                } as PluginMediaEntryEpisodeTab))
                return sortTabs([...otherTabs, ...extensionTabs])
            })
        })
    }, "")

    usePluginListenMediaEntryEpisodeTabEpisodeCollectionEvent((event, extensionId) => {
        const viewId = getPluginEpisodeTabViewId(extensionId)
        startTransition(() => {
            setCollections(prev => ({
                ...prev,
                [viewId]: event.episodeCollection as Media_EpisodeCollection,
            }))
        })
    }, "")

    const selectedTab = tabs.find(tab => tab.viewId === currentView)

    React.useEffect(() => {
        tabs.forEach(tab => {
            sendMediaEntryEpisodeTabStateChangedEvent({
                isOpen: tab.viewId === currentView,
            }, tab.extensionId)
        })
    }, [currentView, tabs])

    React.useEffect(() => {
        if (!selectedTab || !mediaId) return
        sendMediaEntryEpisodeTabOpenEvent({
            mediaId,
        }, selectedTab.extensionId)
    }, [mediaId, selectedTab, sendMediaEntryEpisodeTabOpenEvent])

    return {
        tabs,
        renderedExtensionIds,
    }
}

export function usePluginMediaEntryEpisodeTabs(props: {
    mediaId: number
    currentView: string
    setView: (view: string) => void
}) {
    const { mediaId, currentView, setView } = props

    const tabs = useAtomValue(__plugin_episodeTabsAtom)
    const collections = useAtomValue(__plugin_episodeTabCollectionsAtom)
    const renderedExtensionIds = useAtomValue(__plugin_episodeTabRenderedExtensionIdsAtom)

    const selectedTab = tabs.find(tab => tab.viewId === currentView)

    const { sendMediaEntryEpisodeTabSelectEpisodeEvent } = usePluginSendMediaEntryEpisodeTabSelectEpisodeEvent()

    const selectEpisode = React.useCallback((episode: Media_Episode) => {
        if (!selectedTab || !mediaId) return

        sendMediaEntryEpisodeTabSelectEpisodeEvent({
            mediaId,
            episodeNumber: episode.episodeNumber,
            aniDbEpisode: episode.aniDBEpisode ?? "",
            episode,
        }, selectedTab.extensionId)
    }, [mediaId, selectedTab])

    return {
        tabs,
        selectedTab,
        selectedEpisodeCollection: selectedTab ? collections[selectedTab.viewId] : undefined,
        renderedExtensionIds,
        selectEpisode,
    }
}

export function PluginMediaEntryEpisodeTabContent(props: {
    entry: Media_Entry
    tab: PluginMediaEntryEpisodeTab
    episodeCollection: Media_EpisodeCollection | undefined
    bottomSection?: React.ReactNode
    onSelectEpisode: (episode: Media_Episode) => void
}) {
    const { entry, tab, episodeCollection, bottomSection, onSelectEpisode } = props

    if (!episodeCollection) {
        return <LoadingSpinner />
    }

    return <>
        <div className="h-10" />

        {/* {episodeCollection.hasMappingError && (
         <div data-plugin-media-entry-episode-tab-no-metadata-message-container>
         <p className="text-red-200 opacity-50">
         No metadata info available for this anime. Episode mapping may be incomplete.
         </p>
         </div>
         )} */}

        <TorrentStreamEpisodeSection
            contextType={`episodeTab:${tab.extensionId}`}
            episodeCollection={episodeCollection}
            entry={entry}
            onEpisodeClick={onSelectEpisode}
            onPlayNextEpisodeOnMount={onSelectEpisode}
            bottomSection={bottomSection}
        />
    </>
}

export function PluginAnimeEntryTabIcon(props: { icon?: string, className?: string }) {
    const { icon, className, ...rest } = props

    if (!icon) {
        return <BiExtension className={className} aria-hidden="true" {...rest} />
    }

    if (icon.startsWith("http://") || icon.startsWith("https://") || icon.startsWith("data:image/") || icon.startsWith("/")) {
        return <img
            src={icon}
            alt=""
            className={cn("inline-block size-4 rounded-sm object-contain", className)}
            aria-hidden="true"
            {...rest}
        />
    }

    return <span
        {...props}
        className={cn("", className)}
        dangerouslySetInnerHTML={{ __html: icon }}
    />
}
