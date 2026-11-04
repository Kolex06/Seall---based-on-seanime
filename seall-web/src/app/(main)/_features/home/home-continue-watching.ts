import { Media_Episode } from "@/api/generated/types"
import { useGetRawMediaCollection } from "@/api/hooks/simkl.hooks"
import React from "react"

export function useHomeContinueWatching(playbackEpisodes: Media_Episode[]) {
    const { data } = useGetRawMediaCollection()
    return React.useMemo(() => {
        const result = [...playbackEpisodes]
        const seen = new Set(result.map(episode => episode.baseAnime?.id))
        for (const list of data?.MediaListCollection?.lists ?? []) {
            if (list.status !== "CURRENT" && list.status !== "REPEATING") continue
            for (const entry of list.entries ?? []) {
                const media = entry.media
                if (!media || seen.has(media.id) || media.status === "NOT_YET_RELEASED") continue
                const next = (entry.progress ?? 0) + 1
                if (media.episodes && next > media.episodes) continue
                seen.add(media.id)
                result.push({
                    type: "main", baseAnime: media,
                    episodeTitle: media.title?.userPreferred ?? "",
                    displayTitle: media.format === "MOVIE" ? "Movie" : `Episode ${next}`,
                    episodeNumber: next, aniDBEpisode: String(next), absoluteEpisodeNumber: next,
                    progressNumber: next, isDownloaded: false, isInvalid: false, _isNakamaEpisode: false,
                })
            }
        }
        return result
    }, [data, playbackEpisodes])
}
