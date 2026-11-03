export function simklMediaCategory(format?: string | null) {
    switch (format) {
        case "MOVIE": return { mediaType: "movies" as const, format: undefined }
        case "TV": return { mediaType: "shows" as const, format: undefined }
        case "ANIME": return { mediaType: "anime" as const, format: undefined }
        default: return { format: format || undefined }
    }
}
