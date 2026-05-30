import {
    Media_Episode,
    Media_EpisodeCollection,
    Metadata_AnimeMetadata,
    Onlinestream_Episode,
} from "@/api/generated/types"

export const ALL_SEASONS_VALUE = "all"

type SeasonOption = {
    label: string
    value: string
}

export type EpisodeSeasonGroup<T> = {
    seasonNumber: number | null
    label: string
    episodes: T[]
}

function validSeasonNumber(value: unknown): number | null {
    const season = typeof value === "number" ? value : Number(value)
    if (!Number.isFinite(season) || season <= 0) return null
    return season
}

function uniqueSortedSeasonNumbers(seasons: Array<number | null | undefined>) {
    return Array.from(new Set(seasons.filter((season): season is number => !!validSeasonNumber(season))))
        .sort((a, b) => a - b)
}

export function seasonOptionsFromNumbers(seasons: Array<number | null | undefined>): SeasonOption[] {
    const numbers = uniqueSortedSeasonNumbers(seasons)
    if (!numbers.length) return []

    return [
        { label: "All seasons", value: ALL_SEASONS_VALUE },
        ...numbers.map(season => ({ label: `Season ${season}`, value: String(season) })),
    ]
}

function seasonGroupLabel(seasonNumber: number | null) {
    return seasonNumber ? `Season ${seasonNumber}` : "Other episodes"
}

export function groupEpisodesBySeason<T>(
    episodes: Array<T> | undefined,
    getSeasonNumber: (episode: T) => number | null | undefined,
): EpisodeSeasonGroup<T>[] {
    const groups = new Map<string, EpisodeSeasonGroup<T>>()

    for (const episode of episodes ?? []) {
        const seasonNumber = validSeasonNumber(getSeasonNumber(episode))
        const key = seasonNumber ? String(seasonNumber) : "other"
        const existing = groups.get(key)

        if (existing) {
            existing.episodes.push(episode)
            continue
        }

        groups.set(key, {
            seasonNumber,
            label: seasonGroupLabel(seasonNumber),
            episodes: [episode],
        })
    }

    return Array.from(groups.values())
        .sort((a, b) => {
            if (a.seasonNumber && b.seasonNumber) return a.seasonNumber - b.seasonNumber
            if (a.seasonNumber) return -1
            if (b.seasonNumber) return 1
            return 0
        })
}

function metadataEpisodeCandidates(episode: Media_Episode) {
    return [
        episode.aniDBEpisode,
        String(episode.episodeNumber || ""),
        String(episode.progressNumber || ""),
        String(episode.absoluteEpisodeNumber || ""),
    ].filter(Boolean)
}

export function getMediaEpisodeSeasonNumber(episode: Media_Episode | null | undefined, metadata?: Metadata_AnimeMetadata | null) {
    if (!episode) return null

    const directSeason = validSeasonNumber((episode as any).seasonNumber)
    if (directSeason) return directSeason

    const parsedFileSeason = validSeasonNumber(episode.fileMetadata?.type === "main" ? episode.localFile?.parsedInfo?.season : undefined)
    if (parsedFileSeason) return parsedFileSeason

    const episodeMetadata = metadata?.episodes
    if (!episodeMetadata) return null

    for (const candidate of metadataEpisodeCandidates(episode)) {
        const season = validSeasonNumber(episodeMetadata[candidate]?.seasonNumber)
        if (season) return season
    }

    const metadataMatch = Object.values(episodeMetadata).find(metadataEpisode => {
        return metadataEpisode?.episodeNumber === episode.episodeNumber
            || metadataEpisode?.absoluteEpisodeNumber === episode.absoluteEpisodeNumber
            || metadataEpisode?.episode === episode.aniDBEpisode
    })

    return validSeasonNumber(metadataMatch?.seasonNumber)
}

export function getTorrentEpisodeSeasonOptions(episodeCollection: Media_EpisodeCollection | undefined) {
    return seasonOptionsFromNumbers((episodeCollection?.episodes ?? [])
        .map(episode => getMediaEpisodeSeasonNumber(episode, episodeCollection?.metadata)))
}

export function getTorrentEpisodeSeasonGroups(episodeCollection: Media_EpisodeCollection | undefined) {
    return groupEpisodesBySeason(episodeCollection?.episodes, episode => getMediaEpisodeSeasonNumber(episode, episodeCollection?.metadata))
}

export function filterEpisodeCollectionBySeason(episodeCollection: Media_EpisodeCollection | undefined, selectedSeason: string) {
    if (!episodeCollection || selectedSeason === ALL_SEASONS_VALUE) return episodeCollection

    const season = validSeasonNumber(selectedSeason)
    if (!season) return episodeCollection

    return {
        ...episodeCollection,
        episodes: episodeCollection.episodes?.filter(episode => getMediaEpisodeSeasonNumber(episode, episodeCollection.metadata) === season),
    }
}

function parseSeasonNumberFromText(value: string | null | undefined) {
    if (!value) return null
    const match = value.match(/\b(?:season|s)\s*0?([1-9]\d?)\b/i)
    return validSeasonNumber(match?.[1])
}

export function getOnlineStreamEpisodeSeasonNumber(episode: Onlinestream_Episode | null | undefined) {
    if (!episode) return null

    const directSeason = validSeasonNumber(episode.seasonNumber)
    if (directSeason) return directSeason

    return getMediaEpisodeSeasonNumber(episode.metadata, undefined) ?? parseSeasonNumberFromText(episode.title)
}

export function getOnlineStreamSeasonOptions(episodes: Array<Onlinestream_Episode> | undefined) {
    return seasonOptionsFromNumbers((episodes ?? []).map(getOnlineStreamEpisodeSeasonNumber))
}

export function getOnlineStreamEpisodeSeasonGroups(episodes: Array<Onlinestream_Episode> | undefined) {
    return groupEpisodesBySeason(episodes, getOnlineStreamEpisodeSeasonNumber)
}

export function filterOnlineStreamEpisodesBySeason(episodes: Array<Onlinestream_Episode> | undefined, selectedSeason: string) {
    if (!episodes || selectedSeason === ALL_SEASONS_VALUE) return episodes

    const season = validSeasonNumber(selectedSeason)
    if (!season) return episodes

    return episodes.filter(episode => getOnlineStreamEpisodeSeasonNumber(episode) === season)
}
