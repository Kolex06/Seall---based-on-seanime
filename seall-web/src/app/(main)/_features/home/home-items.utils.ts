import { Models_HomeItem, Nullish } from "@/api/generated/types"
import { ADVANCED_SEARCH_COUNTRIES_MANGA, ADVANCED_SEARCH_MEDIA_GENRES } from "@/app/(main)/search/_lib/advanced-search-constants"

export const DEFAULT_HOME_ITEMS: Models_HomeItem[] = [
    {
        id: "discover-header",
        type: "discover-header",
        schemaVersion: 1,
    },
    {
        id: "anime-continue-watching",
        type: "anime-continue-watching",
        schemaVersion: 1,
    },
    {
        id: "anime-library",
        type: "anime-library",
        schemaVersion: 1,
        options: {
            statuses: ["CURRENT", "PAUSED", "PLANNING", "COMPLETED", "DROPPED"],
            layout: "grid",
        },
    },
]

export function isAnimeLibraryItemsOnly(items: Nullish<Models_HomeItem[]>) {
    if (!items) return true

    for (const item of items) {
        if (![
            "anime-continue-watching",
            "anime-library",
            "anime-continue-watching-header",
            "local-anime-library",
            "local-anime-library-stats",
            "library-upcoming-episodes",
        ].includes(item.type)) {
            return false
        }
    }
    return true
}

type HomeItemSchema = {
    name: string
    kind: ("row" | "header")[]
    options?: { label: string, name: string, type: string, options?: any[] }[]
    schemaVersion: number
    description?: string
}

const _carouselOptions = [
    {
        label: "Name",
        type: "text",
        name: "name",
    },
    {
        label: "Sorting",
        type: "select",
        name: "sorting",
        options: [
            {
                label: "Popular",
                value: "POPULARITY_DESC",
            },
            {
                label: "Trending",
                value: "TRENDING_DESC",
            },
            {
                label: "Romaji Title (A-Z)",
                value: "TITLE_ROMAJI_ASC",
            },
            {
                label: "Romaji Title (Z-A)",
                value: "TITLE_ROMAJI_DESC",
            },
            {
                label: "English title (A-Z)",
                value: "TITLE_ENGLISH_ASC",
            },
            {
                label: "English title (Z-A)",
                value: "TITLE_ENGLISH_DESC",
            },
            {
                label: "Score (0-10)",
                value: "SCORE",
            },
            {
                label: "Score (10-0)",
                value: "SCORE_DESC",
            },
        ],
    },
    {
        label: "Status",
        type: "multi-select",
        name: "status",
        options: [
            {
                label: "Releasing",
                value: "RELEASING",
            },
            {
                label: "Finished",
                value: "FINISHED",
            },
            {
                label: "Not yet released",
                value: "NOT_YET_RELEASED",
            },
        ],
    },
    {
        label: "Format",
        type: "select",
        name: "format",
        options: [
            {
                label: "TV",
                value: "TV",
            },
            {
                label: "Movie",
                value: "MOVIE",
            },
            {
                label: "OVA",
                value: "OVA",
            },
            {
                label: "ONA",
                value: "ONA",
            },
            {
                label: "Special",
                value: "SPECIAL",
            },
        ],
    },
    {
        label: "Genres",
        type: "multi-select",
        options: ADVANCED_SEARCH_MEDIA_GENRES.map(n => ({ value: n, label: n })),
        name: "genres",
    },
    {
        label: "Season",
        type: "select",
        name: "season",
        options: [
            { value: "WINTER", label: "Winter" },
            { value: "SPRING", label: "Spring" },
            { value: "SUMMER", label: "Summer" },
            { value: "FALL", label: "Fall" },
        ],
    },
    {
        label: "Year",
        type: "number",
        name: "year",
        min: 0,
        max: 2100,
    },
    {
        label: "Country of Origin",
        type: "select",
        name: "countryOfOrigin",
        options: ADVANCED_SEARCH_COUNTRIES_MANGA,
    },
]

export const HOME_ITEMS = {
    "centered-title": {
        name: "Centered title",
        kind: ["row"],
        schemaVersion: 1,
        description: "Display a centered title text.",
        options: [{
            label: "Text",
            type: "text",
            name: "text",
        }],
    },
    "anime-continue-watching": {
        name: "Continue Watching",
        kind: ["row", "header"],
        schemaVersion: 1,
        description: "Display a list of episodes you are currently watching.",
    },
    "anime-continue-watching-header": {
        name: "Continue Watching Header",
        kind: ["header"],
        schemaVersion: 1,
        description: "Display a header with a carousel of media you are currently watching.",
    },
    "anime-library": {
        name: "Media Library",
        kind: ["row"],
        schemaVersion: 2,
        description: "Display media you have downloaded / you are currently watching by status.",
        options: [
            {
                label: "Statuses",
                name: "statuses",
                type: "multi-select",
                options: [
                    {
                        value: "CURRENT",
                        label: "Currently Watching",
                    },
                    {
                        value: "PAUSED",
                        label: "Paused",
                    },
                    {
                        value: "PLANNING",
                        label: "Planning",
                    },
                    {
                        value: "COMPLETED",
                        label: "Completed",
                    },
                    {
                        value: "DROPPED",
                        label: "Dropped",
                    },
                ],
            },
            {
                label: "Layout",
                name: "layout",
                type: "select",
                options: [
                    {
                        label: "Grid",
                        value: "grid",
                    },
                    {
                        label: "Carousel",
                        value: "carousel",
                    },
                ],
            },
        ],
    },
    "my-lists": {
        name: "Watchlist",
        kind: ["row"],
        schemaVersion: 1,
        description: "Display media from your lists by status.",
        options: [
            {
                label: "Statuses",
                name: "statuses",
                type: "multi-select",
                options: [
                    {
                        value: "CURRENT",
                        label: "Current",
                    },
                    {
                        value: "REPEATING",
                        label: "Repeating",
                    },
                    {
                        value: "PAUSED",
                        label: "Paused",
                    },
                    {
                        value: "PLANNING",
                        label: "Planning",
                    },
                    {
                        value: "COMPLETED",
                        label: "Completed",
                    },
                    {
                        value: "DROPPED",
                        label: "Dropped",
                    },
                ],
            },
            {
                label: "Layout",
                name: "layout",
                type: "select",
                options: [
                    {
                        label: "Grid",
                        value: "grid",
                    },
                    {
                        label: "Carousel",
                        value: "carousel",
                    },
                ],
            },
            {
                label: "Type",
                name: "type",
                type: "select",
                options: [
                    {
                        label: "Watchlist",
                        value: "anime",
                    },
                ],
            },
            {
                label: "Custom list name (Optional)",
                type: "text",
                name: "customListName",
            },
        ],
    },
    "local-anime-library": {
        name: "Local Media Library",
        kind: ["row"],
        schemaVersion: 2,
        description: "Display a complete grid of media you have in your local library.",
        options: [
            {
                label: "Layout",
                name: "layout",
                type: "select",
                options: [
                    {
                        label: "Grid",
                        value: "grid",
                    },
                    {
                        label: "Carousel",
                        value: "carousel",
                    },
                ],
            },
        ],
    },
    "library-upcoming-episodes": {
        name: "Upcoming Library Episodes",
        kind: ["row"],
        schemaVersion: 1,
        description: "Display a carousel of upcoming releases from media you have in your library.",
    },
    "aired-recently": {
        name: "Released Recently (Global)",
        kind: ["row"],
        schemaVersion: 1,
        description: "Display a carousel of media released recently.",
    },
    "missed-sequels": {
        name: "Missed Sequels",
        kind: ["row"],
        schemaVersion: 1,
        description: "Display a carousel of sequels that aren't in your collection.",
    },
    "anime-schedule-calendar": {
        name: "Schedule Calendar",
        kind: ["row"],
        schemaVersion: 2,
        description: "Display a calendar of releases based on the SIMKL schedule.",
        options: [
            {
                label: "Type",
                name: "type",
                type: "select",
                options: [
                    {
                        label: "Watchlist",
                        value: "my-lists",
                    },
                    {
                        label: "Global",
                        value: "global",
                    },
                ],
            },
        ],
    },
    "local-anime-library-stats": {
        name: "Local Media Library Stats",
        kind: ["row"],
        schemaVersion: 1,
        description: "Display the stats for your local media library.",
    },
    "discover-header": {
        name: "Discover Header",
        kind: ["header"],
        schemaVersion: 1,
        description: "Display a header with a carousel of media that are trending.",
    },
    "anime-carousel": {
        name: "Media Carousel",
        kind: ["row"],
        schemaVersion: 3,
        options: _carouselOptions,
        description: "Display a carousel of media based on the selected options.",
    },
} as Record<string, HomeItemSchema>

export const HOME_ITEM_IDS = Object.keys(HOME_ITEMS) as (keyof typeof HOME_ITEMS)[]

// export type HomeItemID = (keyof typeof HOME_ITEMS)
