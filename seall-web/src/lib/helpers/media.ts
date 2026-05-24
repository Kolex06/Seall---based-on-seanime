import { MediaAPI_MediaListEntry, MediaAPI_BaseMedia, MediaAPI_MangaListEntry, Nullish } from "@/api/generated/types"

export function simkl_getTotalEpisodes(anime: Nullish<MediaAPI_BaseMedia>) {
    if (!anime) return -1
    let maxEp = anime?.episodes ?? -1
    if (maxEp === -1) {
        if (anime.nextAiringEpisode && anime.nextAiringEpisode.episode) {
            maxEp = anime.nextAiringEpisode.episode - 1
        }
    }
    if (maxEp === -1) {
        return 0
    }
    return maxEp
}

export function simkl_getCurrentEpisodes(anime: Nullish<MediaAPI_BaseMedia>) {
    if (!anime) return -1
    let maxEp = -1
    if (anime.nextAiringEpisode && anime.nextAiringEpisode.episode) {
        maxEp = anime.nextAiringEpisode.episode - 1
    }
    if (maxEp === -1) {
        maxEp = anime.episodes ?? 0
    }
    return maxEp
}

export function simkl_getListDataFromEntry(entry: Nullish<MediaAPI_MediaListEntry | MediaAPI_MangaListEntry>) {
    return {
        progress: entry?.progress,
        score: entry?.score,
        status: entry?.status,
        startedAt: new Date(entry?.startedAt?.year || 0,
            entry?.startedAt?.month ? entry?.startedAt?.month - 1 : 0,
            entry?.startedAt?.day || 0).toUTCString(),
        completedAt: new Date(entry?.completedAt?.year || 0,
            entry?.completedAt?.month ? entry?.completedAt?.month - 1 : 0,
            entry?.completedAt?.day || 0).toUTCString(),
    }
}


export function simkl_animeIsMovie(anime: Nullish<MediaAPI_BaseMedia>) {
    if (!anime) return false
    return anime?.format === "MOVIE"

}

export function simkl_animeIsSingleEpisode(anime: Nullish<MediaAPI_BaseMedia>) {
    if (!anime) return false
    return anime?.format === "MOVIE" || anime?.episodes === 1
}


export function simkl_getUnwatchedCount(anime: Nullish<MediaAPI_BaseMedia>, progress: Nullish<number>) {
    if (!anime) return false
    const maxEp = simkl_getCurrentEpisodes(anime)
    return maxEp - (progress ?? 0)
}

