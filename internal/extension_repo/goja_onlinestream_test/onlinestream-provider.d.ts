declare type SearchResult = {
    id: string
    title: string
    url: string
    subOrDub: SubOrDub
}

declare type SubOrDub = "sub" | "dub" | "both"

declare type EpisodeDetails = {
    id: string
    number: number
    url: string
    title?: string
}

declare type EpisodeServer = {
    server: string
    headers: { [key: string]: string }
    videoSources: VideoSource[]
}

declare type VideoSourceType = "mp4" | "m3u8" | "unknown"

declare type VideoSource = {
    url: string
    type: VideoSourceType
    // Quality or label of the video source, should be unique (e.g. "1080p", "1080p - English")
    quality: string
    // Secondary label of the video source (e.g. "English")
    label?: string
    subtitles: VideoSubtitle[]
}

declare type VideoSubtitle = {
    id: string
    url: string
    language: string
    isDefault: boolean
}

declare interface Media {
    id: number
    idMal?: number
    siteUrl?: string
    mediaType?: "movies" | "shows" | "anime" | "all"
    status?: string
    format?: string
    englishTitle?: string
    romajiTitle?: string
    episodeCount?: number
    absoluteSeasonOffset?: number
    synonyms: string[]
    isAdult: boolean
    startDate?: FuzzyDate
}

declare interface FuzzyDate {
    year: number
    month?: number
    day?: number
}

declare type SearchOptions = {
    media: Media
    query: string
    dub: boolean
    year?: number
    mediaType?: "movies" | "shows" | "anime" | "all"
}

declare type Settings = {
    episodeServers: string[]
    supportsDub: boolean
    supportedMediaTypes?: Array<"movies" | "shows" | "anime" | "all">
}

declare abstract class StreamProvider {
    search(opts: SearchOptions): Promise<SearchResult[]>

    findMediaItems(id: string): Promise<EpisodeDetails[]>

    findMediaItemServer(episode: EpisodeDetails, server: string): Promise<EpisodeServer>

    getSettings(): Settings
}
