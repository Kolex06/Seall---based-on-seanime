import { MediaAPI_MediaFormat } from "@/api/generated/types"

export const ADVANCED_SEARCH_MEDIA_GENRES = [
    "Action",
    "Adventure",
    "Animation",
    "Biography",
    "Comedy",
    "Crime",
    "Documentary",
    "Drama",
    "Ecchi",
    "Educational",
    "Family",
    "Fantasy",
    "Gag Humor",
    "Gore",
    "Harem",
    "History",
    "Historical",
    "Horror",
    "Idol",
    "Isekai",
    "Josei",
    "Kids",
    "Magic",
    "Martial Arts",
    "Mecha",
    "Military",
    "Music",
    "Musical",
    "Mystery",
    "Mythology",
    "Parody",
    "Psychological",
    "Racing",
    "Reality",
    "Reincarnation",
    "Romance",
    "Samurai",
    "School",
    "Science Fiction",
    "Science-Fiction",
    "Sci-Fi",
    "Seinen",
    "Shoujo",
    "Shoujo Ai",
    "Shounen",
    "Shounen Ai",
    "Slice of Life",
    "Space",
    "Sports",
    "Strategy Game",
    "Super Power",
    "Supernatural",
    "Talk Show",
    "Thriller",
    "Vampire",
    "War",
    "Western",
    "Yaoi",
    "Yuri",
]

export const ADVANCED_SEARCH_SEASONS = [
    "Winter",
    "Spring",
    "Summer",
    "Fall",
]

export const ADVANCED_SEARCH_FORMATS: { value: MediaAPI_MediaFormat, label: string }[] = [
    { value: "TV", label: "Series" },
    { value: "MOVIE", label: "Movie" },
    { value: "SPECIAL", label: "Special" },
]

export const ADVANCED_SEARCH_FORMATS_MANGA: { value: MediaAPI_MediaFormat, label: string }[] = [
    { value: "MANGA", label: "Manga" },
    { value: "ONE_SHOT", label: "One Shot" },
]


export const ADVANCED_SEARCH_COUNTRIES_MANGA: { value: string, label: string }[] = [
    { value: "JP", label: "Japan" },
    { value: "KR", label: "South Korea" },
    { value: "CN", label: "China" },
    { value: "TW", label: "Taiwan" },
]

export const ADVANCED_SEARCH_STATUS = [
    { value: "FINISHED", label: "Finished" },
    { value: "RELEASING", label: "Releasing" },
    { value: "NOT_YET_RELEASED", label: "Upcoming" },
    { value: "HIATUS", label: "Hiatus" },
    { value: "CANCELLED", label: "Cancelled" },
]

export const ADVANCED_SEARCH_SORTING = [
    { value: "TRENDING_DESC", label: "Trending" },
    { value: "START_DATE_DESC", label: "Release date" },
    { value: "SCORE_DESC", label: "Highest score" },
    { value: "POPULARITY_DESC", label: "Most popular" },
    { value: "EPISODES_DESC", label: "Number of episodes" },
]

export const ADVANCED_SEARCH_SORTING_MANGA = [
    { value: "TRENDING_DESC", label: "Trending" },
    { value: "START_DATE_DESC", label: "Release date" },
    { value: "SCORE_DESC", label: "Highest score" },
    { value: "POPULARITY_DESC", label: "Most popular" },
    { value: "CHAPTERS_DESC", label: "Number of chapters" },
]

export const ADVANCED_SEARCH_TYPE = [
    { value: "anime", label: "Media" },
]
