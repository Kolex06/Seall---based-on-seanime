import { MediaAPI_MediaFormat, MediaAPI_MediaSeason, MediaAPI_MediaSort, MediaAPI_MediaStatus } from "@/api/generated/types"
import { atomWithImmer } from "jotai-immer"

type Params = {
    active: boolean
    title: string | null
    sorting: MediaAPI_MediaSort[] | null
    genre: string[] | null
    tags: string[] | null
    status: MediaAPI_MediaStatus[] | null
    format: MediaAPI_MediaFormat | null
    season: MediaAPI_MediaSeason | null
    year: string | null
    minScore: string | null
    isAdult: boolean
    countryOfOrigin: string | null
    type: "anime" | "manga"
}

export const __advancedSearch_paramsAtom = atomWithImmer<Params>({
    active: true,
    title: null,
    sorting: null,
    status: null,
    genre: null,
    tags: null,
    format: null,
    season: null,
    year: null,
    minScore: null,
    isAdult: false,
    countryOfOrigin: null,
    type: "anime",
})

export function __advancedSearch_getValue<T extends any>(value: T | ""): any {
    if (value === "") return undefined
    if (Array.isArray(value) && value.filter(Boolean).length === 0) return undefined
    if (typeof value === "string" && !isNaN(parseInt(value))) return Number(value)
    if (value === null) return undefined
    return value
}
