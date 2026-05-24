/// <reference path="../goja_plugin_types/app.d.ts" />

declare type Settings = {
    supportsAnime: boolean
    supportsManga: boolean
}

declare type ListResponse<T extends $app.MediaAPI_BaseMedia | $app.MediaAPI_BaseManga> = {
    media: T[]
    page: number
    totalPages: number
    total: number
}

declare abstract class CustomSource {
    getSettings(): Settings

    async getAnime(ids: number[]): Promise<$app.MediaAPI_BaseMedia[]>

    async getAnimeMetadata(id: number): Promise<$app.Metadata_AnimeMetadata | null>

    async getAnimeWithRelations(id: number): Promise<$app.MediaAPI_CompleteMedia>

    async getAnimeDetails(id: number): Promise<$app.MediaAPI_MediaDetailsById_Media | null>

    async getManga(ids: number[]): Promise<$app.MediaAPI_BaseManga[]>

    async listAnime(search: string, page: number, perPage: number): Promise<ListResponse<$app.MediaAPI_BaseMedia>>

    async getMangaDetails(id: number): Promise<$app.MediaAPI_MangaDetailsById_Media | null>

    async listManga(search: string, page: number, perPage: number): Promise<ListResponse<$app.MediaAPI_BaseManga>>
}
