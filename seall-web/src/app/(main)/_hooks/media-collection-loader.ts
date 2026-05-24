import { Media_EntryListData, Nullish } from "@/api/generated/types"
import { useGetMediaCollection } from "@/api/hooks/simkl.hooks"
import { __simkl_userAnimeListDataAtom, __media_userMediaAtom } from "@/app/(main)/_atoms/media.atoms"
import { atom } from "jotai"
import { useAtomValue, useSetAtom } from "jotai/react"
import { selectAtom } from "jotai/utils"
import React from "react"

const emptyAnimeListDataAtom = atom<Media_EntryListData | undefined>(undefined)

/**
 * @description
 * - Fetches the Simkl collection
 */
export function useAnimeCollectionLoader() {
    const setSimklUserMedia = useSetAtom(__media_userMediaAtom)

    const setSimklUserMediaListData = useSetAtom(__simkl_userAnimeListDataAtom)

    const { data } = useGetMediaCollection()

    // Store the user's media in `userMediaAtom`
    React.useEffect(() => {
        if (!!data) {
            const allMedia = data.MediaListCollection?.lists?.flatMap(n => n?.entries)?.filter(Boolean)?.map(n => n.media)?.filter(Boolean) ?? []
            setSimklUserMedia(allMedia)

            const listData = data.MediaListCollection?.lists?.flatMap(n => n?.entries)?.filter(Boolean)?.reduce((acc, n) => {
                acc[String(n.media?.id!)] = {
                    status: n.status,
                    progress: n.progress || 0,
                    score: n.score || 0,
                    startedAt: (n.startedAt?.year && n.startedAt?.month) ? new Date(n.startedAt.year || 0,
                        (n.startedAt.month || 1) - 1,
                        n.startedAt.day || 1).toISOString() : undefined,
                    completedAt: (n.completedAt?.year && n.completedAt?.month) ? new Date(n.completedAt.year || 0,
                        (n.completedAt.month || 1) - 1,
                        n.completedAt.day || 1).toISOString() : undefined,
                }
                return acc
            }, {} as Record<string, Media_EntryListData>)
            setSimklUserMediaListData(listData || {})
        }
    }, [data])

    return null
}

export function useSimklUserAnime() {
    return useAtomValue(__media_userMediaAtom)
}

export function useMediaUserListData(mId: Nullish<number | string>, enabled: boolean = true): Media_EntryListData | undefined {
    const mediaId = String(mId)
    const listDataAtom = React.useMemo(() => {
        if (!enabled) {
            return emptyAnimeListDataAtom
        }

        return selectAtom(__simkl_userAnimeListDataAtom, data => data[mediaId])
    }, [enabled, mediaId])

    return useAtomValue(listDataAtom)
}
