import { CustomLibraryBanner } from "@/app/(main)/_features/anime-library/_containers/custom-library-banner"
import { MediaCollectionLists } from "@/app/(main)/lists/_containers/media-collection-lists"
import { PageWrapper } from "@/components/shared/page-wrapper"
import React from "react"


export default function Home() {

    return (
        <>
            <CustomLibraryBanner discrete />
            <PageWrapper
                className="p-4 sm:p-8 pt-4 relative"
                data-simkl-page
            >
                <MediaCollectionLists />
            </PageWrapper>
        </>
    )
}
