import { MediaAPI_BaseMedia, Media_EntryListData, Manga_EntryListData } from "@/api/generated/types"
import { atom } from "jotai"

export const __media_userMediaAtom = atom<MediaAPI_BaseMedia[] | undefined>(undefined)

// e.g. { "123": { ... } }
export const __simkl_userAnimeListDataAtom = atom<Record<string, Media_EntryListData>>({})

export const __simkl_userMangaListDataAtom = atom<Record<string, Manga_EntryListData>>({})
