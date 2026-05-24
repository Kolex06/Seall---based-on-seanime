import { Media_MissingEpisodes } from "@/api/generated/types"
import { missingEpisodesAtom, missingSilencedEpisodesAtom } from "@/app/(main)/_atoms/missing-episodes.atoms"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { useSetAtom } from "jotai/react"
import { useEffect } from "react"

/**
 * @description
 * - Sets missing episodes to the atom so that it can be displayed in other components
 * - Filters out adult content if the user has it disabled
 * @param data
 */
export function useHandleMissingEpisodes(data: Media_MissingEpisodes | undefined) {
    const serverStatus = useServerStatus()
    const setAtom = useSetAtom(missingEpisodesAtom)
    const setSilencedAtom = useSetAtom(missingSilencedEpisodesAtom)

    useEffect(() => {
        if (!!data) {
            if (serverStatus?.settings?.simkl?.enableAdultContent) {
                setAtom(data?.episodes ?? [])
            } else {
                setAtom(data?.episodes?.filter(episode => !episode?.baseAnime?.isAdult) ?? [])
            }
            setSilencedAtom(data?.silencedEpisodes ?? [])
        }
    }, [data, serverStatus?.settings?.simkl?.enableAdultContent])

    return {
        missingEpisodes: serverStatus?.settings?.simkl?.enableAdultContent
            ? (data?.episodes ?? [])
            : (data?.episodes?.filter(episode => !episode?.baseAnime?.isAdult) ?? []),
        silencedEpisodes: data?.silencedEpisodes ?? [],
    }

}
