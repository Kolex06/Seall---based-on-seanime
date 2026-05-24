import { useGetMangaCollection } from "@/api/hooks/manga.hooks"
import { MediaCardLazyGrid } from "@/app/(main)/_features/media/_components/media-card-grid"
import { MediaEntryCard } from "@/app/(main)/_features/media/_components/media-entry-card"
import { LuffyError } from "@/components/shared/luffy-error"
import { PageWrapper } from "@/components/shared/page-wrapper"
import { SeaLink } from "@/components/shared/sea-link"
import { Button } from "@/components/ui/button"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import React from "react"

const READING_LIST_LABELS: Record<string, string> = {
    CURRENT: "Reading",
    PLANNING: "Plan to read",
    COMPLETED: "Completed",
    PAUSED: "Paused",
    DROPPED: "Dropped",
    REPEATING: "Rereading",
}

export default function Page() {
    const { data, isLoading } = useGetMangaCollection()

    const lists = React.useMemo(() => {
        return data?.lists?.filter(list => !!list.entries?.length) || []
    }, [data])

    const itemCount = React.useMemo(() => {
        return lists.reduce((acc, list) => acc + (list.entries?.length || 0), 0)
    }, [lists])

    if (isLoading) {
        return <LoadingSpinner className="mt-20" />
    }

    if (!itemCount) {
        return <PageWrapper className="p-4">
            <LuffyError title="No reading found">
                <div className="space-y-2">
                    <p>No reading items have been added to your library yet.</p>
                    <div className="!mt-4">
                        <SeaLink href="/discover">
                            <Button intent="white-outline" rounded>
                                Browse media
                            </Button>
                        </SeaLink>
                    </div>
                </div>
            </LuffyError>
        </PageWrapper>
    }

    return (
        <PageWrapper className="p-4 space-y-10">
            <div>
                <h2>Reading</h2>
            </div>

            {lists.map(list => (
                <section key={list.type || list.status || "reading"} className="space-y-4">
                    <h3>{READING_LIST_LABELS[list.type || ""] || list.status || "Reading"}</h3>
                    <MediaCardLazyGrid itemCount={list.entries?.length || 0}>
                        {list.entries?.map(entry => (
                            entry?.media ? <MediaEntryCard
                                key={entry.media.id}
                                type="manga"
                                media={entry.media}
                                listData={entry.listData}
                            /> : null
                        ))}
                    </MediaCardLazyGrid>
                </section>
            ))}
        </PageWrapper>
    )
}
