import { useGetMediaCollection } from "@/api/hooks/simkl.hooks"
import { useGetMangaCollection } from "@/api/hooks/manga.hooks"
import { useLibraryCollection } from "@/app/(main)/_hooks/anime-library-collection-loader.ts"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { CommandGroup, CommandItem, CommandShortcut } from "@/components/ui/command"
import { useRouter } from "@/lib/navigation"
import React from "react"
import { BiArrowBack } from "react-icons/bi"
import { CommandHelperText, CommandItemMedia } from "./_components/command-utils"
import { useSeaCommandContext } from "./sea-command"
import { seaCommand_compareMediaTitles } from "./utils"

// only rendered when typing "/media", "/library" or "/reading"
export function SeaCommandUserMediaNavigation() {

    const { input, select, command: { isCommand, command, args }, scrollToTop } = useSeaCommandContext()
    const { data: animeCollection, isLoading: isAnimeLoading } = useGetMediaCollection() // should be available instantly
    const { data: mangaCollection, isLoading: isMangaLoading } = useGetMangaCollection()
    const animeLibraryCollection = useLibraryCollection()

    const anime = animeCollection?.MediaListCollection?.lists?.flatMap(n => n?.entries)?.filter(Boolean)?.map(n => n.media)?.filter(Boolean) ?? []
    const manga = mangaCollection?.lists?.flatMap(n => n?.entries)?.filter(Boolean)?.map(n => n.media)?.filter(Boolean) ?? []

    const router = useRouter()

    const query = args.join(" ")
    const filteredMedia = (command === "media" && query.length > 0) ? anime.filter(n => seaCommand_compareMediaTitles(n.title, query)) : []
    const filteredReading = (command === "reading" && query.length > 0) ? manga.filter(n => seaCommand_compareMediaTitles(n.title, query)) : []
    const filteredAnimeLibrary = (command === "library" && query.length > 0) ? animeLibraryCollection?.lists?.flatMap(l => l.entries)
        ?.filter(n => seaCommand_compareMediaTitles(n?.media?.title, query))
        ?.map(n => n?.media)
        ?.filter(Boolean) ?? [] : []

    return (
        <>
            {query.length === 0 && (
                <>
                    <CommandHelperText
                        command="/media [title]"
                        description="Find media in your collection"
                        show={command === "media"}
                    />
                    <CommandHelperText
                        command="/reading [title]"
                        description="Find reading items in your collection"
                        show={command === "reading"}
                    />
                    <CommandHelperText
                        command="/library [title]"
                        description="Find media in your library"
                        show={command === "library"}
                    />
                </>
            )}

            {command === "media" && filteredMedia.length > 0 && (
                <CommandGroup heading="My media">
                    {filteredMedia.map(n => (
                        <CommandItem
                            key={n.id}
                            onSelect={() => {
                                select(() => {
                                    router.push(`/entry?id=${n.id}`)
                                })
                            }}
                        >
                            <CommandItemMedia media={n} type="anime" />
                        </CommandItem>
                    ))}
                </CommandGroup>
            )}

            {command === "library" && filteredAnimeLibrary.length > 0 && (
                <CommandGroup heading="Library media">
                    {filteredAnimeLibrary.map(n => (
                        <CommandItem
                            key={n.id}
                            onSelect={() => {
                                select(() => {
                                    router.push(`/entry?id=${n.id}`)
                                })
                            }}
                        >
                            <CommandItemMedia media={n} type="anime" />
                        </CommandItem>
                    ))}
                </CommandGroup>
            )}
            {command === "reading" && filteredReading.length > 0 && (
                <CommandGroup heading="My reading">
                    {filteredReading.map(n => (
                        <CommandItem
                            key={n.id}
                            onSelect={() => {
                                select(() => {
                                    router.push(`/reading/entry?id=${n.id}`)
                                })
                            }}
                        >
                            <CommandItemMedia media={n} type="manga" />
                        </CommandItem>
                    ))}
                </CommandGroup>
            )}
        </>
    )
}

export function SeaCommandNavigation() {

    const serverStatus = useServerStatus()

    const { input, select, command: { isCommand, command, args } } = useSeaCommandContext()

    const router = useRouter()

    const pages = [
        {
            name: "Home",
            href: "/",
            flag: "home",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Schedule",
            href: "/schedule",
            flag: "schedule",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Settings",
            href: "/settings",
            flag: "settings",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Discover",
            href: "/discover",
            flag: "discover",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Watchlist",
            href: "/lists",
            flag: "lists",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Auto Downloader",
            href: "/auto-downloader",
            flag: "auto-downloader",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Torrent list",
            href: "/torrent-list",
            flag: "torrent-list",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Scan summaries",
            href: "/scan-summaries",
            flag: "scan-summaries",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Extensions",
            href: "/extensions",
            flag: "extensions",
            show: !serverStatus?.isOffline,
        },
        {
            name: "Advanced search",
            href: "/search",
            flag: "search",
            show: !serverStatus?.isOffline,
        },
    ]

    // If no args, show all pages
    // If args, show pages that match the args
    const filteredPages = pages.filter(page => page.flag.startsWith(command))


    // if (!input.startsWith("/")) return null


    return (
        <>
            {command.startsWith("ba") && (
                <CommandGroup heading="Navigation">
                    <CommandItem
                        onSelect={() => {
                            select(() => {
                                router.back()
                            })
                        }}
                    >
                        <BiArrowBack className="mr-2 h-4 w-4" />
                        <span>Go back</span>
                    </CommandItem>
                </CommandGroup>
            )}
            {command.startsWith("fo") && (
                <CommandGroup heading="Navigation">
                    <CommandItem
                        onSelect={() => {
                            select(() => {
                                router.forward()
                            })
                        }}
                    >
                        <BiArrowBack className="mr-2 h-4 w-4 rotate-180" />
                        <span>Go forward</span>
                    </CommandItem>
                </CommandGroup>
            )}

            {/*Typing `/library`, `/schedule`, etc. without args*/}
            {isCommand && filteredPages.length > 0 && args.length === 0 && (
                <CommandGroup heading="Screens">
                    <>
                        {filteredPages.filter(page => page.show).map(page => (
                            <CommandItem
                                key={page.flag}
                                onSelect={() => {
                                    select(() => {
                                        router.push(page.href)
                                    })
                                }}
                            >
                                <span className="text-sm tracking-wide font-bold text-[--muted]">Go to:&nbsp;</span>{" "}{page.name}
                                {command === page.flag ? <CommandShortcut>Enter</CommandShortcut> : <CommandShortcut>/{page.flag}</CommandShortcut>}
                            </CommandItem>
                        ))}
                    </>
                </CommandGroup>
            )}
            {(command !== "back" && command !== "forward") && (
                <CommandGroup heading="Navigation">
                    {/* {command === "" && ( */}
                    <>
                        <CommandItem
                            onSelect={() => {
                                select(() => {
                                    router.back()
                                })
                            }}
                        >
                            <BiArrowBack className="mr-2 h-4 w-4" />
                            <span>Go back</span>
                        </CommandItem>
                        <CommandItem
                            onSelect={() => {
                                select(() => {
                                    router.forward()
                                })
                            }}
                        >
                            <BiArrowBack className="mr-2 h-4 w-4 rotate-180" />
                            <span>Go forward</span>
                        </CommandItem>
                    </>
                    {/* )} */}
                </CommandGroup>
            )}
        </>
    )
}
