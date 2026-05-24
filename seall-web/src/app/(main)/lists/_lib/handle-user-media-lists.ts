import { MediaAPI_MediaCollection_MediaListCollection_Lists } from "@/api/generated/types"
import { useGetRawMediaCollection, useGetRawMediaCollectionTags } from "@/api/hooks/simkl.hooks"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { CollectionParams, CollectionType, DEFAULT_COLLECTION_PARAMS, filterEntriesByTitle, filterListEntries } from "@/lib/helpers/filtering"
import { atomWithImmer } from "jotai-immer"
import { useAtom } from "jotai/react"
import React from "react"
import { useDebounce } from "use-debounce"

export const MYLISTS_DEFAULT_PARAMS: CollectionParams<"anime"> | CollectionParams<"manga"> = {
    ...DEFAULT_COLLECTION_PARAMS,
    sorting: "SCORE_DESC",
    unreadOnly: false,
    continueWatchingOnly: false,
}

export const __myListsSearch_paramsAtom = atomWithImmer<CollectionParams<"anime"> | CollectionParams<"manga">>(MYLISTS_DEFAULT_PARAMS)

export const __myListsSearch_paramsInputAtom = atomWithImmer<CollectionParams<"anime"> | CollectionParams<"manga">>(MYLISTS_DEFAULT_PARAMS)

export const __myLists_selectedTypeAtom = atomWithImmer<"anime" | "manga" | "stats">("anime")

export function useHandleUserMediaLists(debouncedSearchInput: string, type?: "anime" | "manga") {
    void type

    const serverStatus = useServerStatus()
    const [selectedType, setSelectedType] = useAtom(__myLists_selectedTypeAtom)
    const { data: animeData } = useGetRawMediaCollection()
    const { data: animeTagMap } = useGetRawMediaCollectionTags()

    const data = React.useMemo(() => {
        return animeData
    }, [animeData])

    const lists = React.useMemo(() => data?.MediaListCollection?.lists, [data])
    const mediaTagMap = React.useMemo(() => {
        return animeTagMap
    }, [animeTagMap])

    const [params, _setParams] = useAtom(__myListsSearch_paramsAtom)
    const [debouncedParams] = useDebounce(params, 500)

    React.useLayoutEffect(() => {
        if (selectedType === "manga") {
            setSelectedType("anime")
        }
    }, [selectedType, setSelectedType])

    React.useLayoutEffect(() => {
        _setParams(MYLISTS_DEFAULT_PARAMS)
    }, [selectedType])

    const _filteredLists: MediaAPI_MediaCollection_MediaListCollection_Lists[] = React.useMemo(() => {
        return lists?.map(obj => {
            if (!obj) return undefined
            const arr = filterListEntries(
                "anime" as CollectionType,
                obj?.entries,
                params,
                serverStatus?.settings?.simkl?.enableAdultContent,
                mediaTagMap,
            )
            return {
                name: obj?.name,
                isCustomList: obj?.isCustomList,
                status: obj?.status,
                entries: arr,
            }
        }).filter(Boolean) ?? []
    }, [lists, debouncedParams, mediaTagMap, serverStatus?.settings?.simkl?.enableAdultContent])

    const filteredLists: MediaAPI_MediaCollection_MediaListCollection_Lists[] = React.useMemo(() => {
        return _filteredLists?.map(obj => {
            if (!obj) return undefined
            const arr = filterEntriesByTitle(obj?.entries, debouncedSearchInput)
            return {
                name: obj?.name,
                isCustomList: obj?.isCustomList,
                status: obj?.status,
                entries: arr,
            }
        })?.filter(Boolean) ?? []
    }, [_filteredLists, debouncedSearchInput])

    const customLists = React.useMemo(() => {
        return filteredLists?.filter(obj => obj?.isCustomList) ?? []
    }, [filteredLists])

    return {
        currentList: React.useMemo(() => filteredLists?.find(l => l?.status === "CURRENT"), [filteredLists]),
        repeatingList: React.useMemo(() => filteredLists?.find(l => l?.status === "REPEATING"), [filteredLists]),
        planningList: React.useMemo(() => filteredLists?.find(l => l?.status === "PLANNING"), [filteredLists]),
        pausedList: React.useMemo(() => filteredLists?.find(l => l?.status === "PAUSED"), [filteredLists]),
        completedList: React.useMemo(() => filteredLists?.find(l => l?.status === "COMPLETED"), [filteredLists]),
        droppedList: React.useMemo(() => filteredLists?.find(l => l?.status === "DROPPED"), [filteredLists]),
        customLists,
    }
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
