import { Media_Entry, Media_Episode, Media_EpisodeCollection } from "@/api/generated/types"
import { getEpisodeMinutesRemaining, getEpisodePercentageComplete, useGetContinuityWatchHistory } from "@/api/hooks/continuity.hooks"
import { EpisodeCard } from "@/app/(main)/_features/anime/_components/episode-card"
import { EpisodeGridItem } from "@/app/(main)/_features/anime/_components/episode-grid-item"
import { MediaEpisodeInfoModal } from "@/app/(main)/_features/media/_components/media-episode-info-modal"
import { PluginEpisodeGridItemMenuItems } from "@/app/(main)/_features/plugin/actions/plugin-actions"
import { EpisodeListGrid, EpisodeListPaginatedGrid } from "@/app/(main)/entry/_components/episode-list-grid"
import { usePlayNextVideoOnMount } from "@/app/(main)/entry/_lib/handle-play-on-mount"
import { episodeCardCarouselItemClass } from "@/components/shared/classnames"
import { IconButton } from "@/components/ui/button"
import { Carousel, CarouselContent, CarouselDotButtons, CarouselItem } from "@/components/ui/carousel"
import { ContextMenuItem } from "@/components/ui/context-menu"
import { DropdownMenu, DropdownMenuItem } from "@/components/ui/dropdown-menu"
import { getTorrentEpisodeSeasonGroups } from "@/lib/helpers/episode-seasons"
import { useThemeSettings } from "@/lib/theme/theme-hooks"
import React, { useMemo } from "react"
import { BiDotsHorizontal } from "react-icons/bi"
import { LuTvMinimalPlay } from "react-icons/lu"

type TorrentStreamEpisodeSectionProps = {
    entry: Media_Entry
    episodeCollection: Media_EpisodeCollection | undefined
    onEpisodeClick: (episode: Media_Episode) => void
    onPlayExternallyEpisodeClick?: (episode: Media_Episode) => void
    onPlayNextEpisodeOnMount: (episode: Media_Episode) => void
    bottomSection?: React.ReactNode
    contextType: "torrentstream" | "debridstream" | string // used for plugin context menu item filtering
}

export function TorrentStreamEpisodeSection(props: TorrentStreamEpisodeSectionProps) {
    const ts = useThemeSettings()

    const {
        entry,
        episodeCollection,
        onEpisodeClick,
        onPlayNextEpisodeOnMount,
        bottomSection,
        onPlayExternallyEpisodeClick,
        contextType,
        ...rest
    } = props

    const { data: watchHistory } = useGetContinuityWatchHistory()


    /**
     * Organize episodes to watch
     */
    const episodesToWatch = useMemo(() => {
        if (!episodeCollection?.episodes) return []
        let ret = [...episodeCollection?.episodes]
        ret = ((!!entry.listData?.progress && !!entry.media?.episodes && entry.listData?.progress === entry.media?.episodes)
                ? ret?.reverse()
                : ret?.slice(entry.listData?.progress || 0)
        )?.slice(0, 30) || []
        return ret
    }, [episodeCollection?.episodes, entry.nextEpisode, entry.listData?.progress])
    const episodeSeasonGroups = useMemo(() => {
        return getTorrentEpisodeSeasonGroups(episodeCollection)
            .map(group => ({
                ...group,
                episodes: [...group.episodes].sort((a, b) => (a.progressNumber || a.episodeNumber || 0) - (b.progressNumber || b.episodeNumber || 0)),
            }))
    }, [episodeCollection])
    const shouldSplitEpisodeGridBySeason = episodeSeasonGroups.length > 1

    /**
     * Play next episode on mount if requested
     */
    usePlayNextVideoOnMount({
        onPlay: () => {
            onPlayNextEpisodeOnMount(episodesToWatch[0])
        },
    }, !!episodesToWatch[0])

    if (!entry || !episodeCollection) return null

    function renderGridEpisodeItem(episode: Media_Episode | undefined, keyPrefix = "") {
        if (!episode) return null

        return (
            <EpisodeGridItem
                key={`${keyPrefix}${episode?.episodeNumber}-${episode?.displayTitle || ""}`}
                media={episode?.baseAnime as any}
                title={episode?.displayTitle || episode?.baseAnime?.title?.userPreferred || ""}
                image={episode?.episodeMetadata?.image || episode?.baseAnime?.coverImage?.large}
                episodeTitle={episode?.episodeTitle}
                onClick={() => {
                    onEpisodeClick(episode)
                }}
                description={episode?.episodeMetadata?.overview}
                isFiller={episode?.episodeMetadata?.isFiller}
                length={episode?.episodeMetadata?.length}
                isWatched={!!entry.listData?.progress && entry.listData.progress >= (episode?.progressNumber || 0)}
                className="flex-none w-full"
                episodeNumber={episode?.episodeNumber}
                watchedProgress={entry.listData?.progress}
                progressNumber={episode?.progressNumber}
                action={<>
                    <MediaEpisodeInfoModal
                        title={episode?.displayTitle}
                        image={episode?.episodeMetadata?.image}
                        episodeTitle={episode?.episodeTitle}
                        airDate={episode?.episodeMetadata?.airDate}
                        length={episode?.episodeMetadata?.length}
                        summary={episode?.episodeMetadata?.overview}
                        isInvalid={episode?.isInvalid}
                    />

                    {!!onPlayExternallyEpisodeClick ? <DropdownMenu
                        trigger={
                            <IconButton
                                icon={<BiDotsHorizontal />}
                                intent="gray-basic"
                                size="xs"
                            />
                        }
                    >

                        {onPlayExternallyEpisodeClick && <DropdownMenuItem
                            onClick={() => {
                                onPlayExternallyEpisodeClick(episode)
                            }}
                        >
                            <LuTvMinimalPlay />
                            Play externally
                        </DropdownMenuItem>}
                        <PluginEpisodeGridItemMenuItems isDropdownMenu={false} type={contextType} episode={episode} />
                    </DropdownMenu> : (
                        <PluginEpisodeGridItemMenuItems isDropdownMenu={true} type={contextType} episode={episode} />
                    )}

                </>}
            />
        )
    }

    return (
        <>
            <Carousel
                className="w-full max-w-full"
                gap="md"
                opts={{
                    align: "start",
                }}
            >
                <CarouselDotButtons />
                <CarouselContent>
                    {episodesToWatch.map((episode, idx) => (
                        <CarouselItem
                            key={episode?.localFile?.path || idx}
                            className={episodeCardCarouselItemClass(ts.smallerEpisodeCarouselSize)}
                        >
                            <EpisodeCard
                                key={episode.localFile?.path || ""}
                                contextType={contextType}
                                episode={episode}
                                image={episode.episodeMetadata?.image || episode.baseAnime?.bannerImage || episode.baseAnime?.coverImage?.extraLarge}
                                topTitle={episode.episodeTitle || episode?.baseAnime?.title?.userPreferred}
                                title={episode.displayTitle}
                                // meta={episode.episodeMetadata?.airDate ?? undefined}
                                isInvalid={episode.isInvalid}
                                progressTotal={episode.baseAnime?.episodes}
                                watchedProgress={entry.listData?.progress}
                                progressNumber={episode.progressNumber}
                                episodeNumber={episode.episodeNumber}
                                length={episode.episodeMetadata?.length}
                                percentageComplete={getEpisodePercentageComplete(watchHistory, entry.mediaId, episode.episodeNumber)}
                                minutesRemaining={getEpisodeMinutesRemaining(watchHistory, entry.mediaId, episode.episodeNumber)}
                                hasDiscrepancy={episodeCollection?.episodes?.findIndex(e => e.type === "special") !== -1}
                                fallbackImage={[episode.baseAnime?.bannerImage, episode.baseAnime?.coverImage?.large,
                                    episode.baseAnime?.coverImage?.extraLarge]}
                                onClick={() => {
                                    onEpisodeClick(episode)
                                }}
                                anime={{
                                    id: entry.mediaId,
                                    image: episode.baseAnime?.coverImage?.medium,
                                    title: episode?.baseAnime?.title?.userPreferred,
                                }}
                                additionalContextMenuItems={<>
                                    {onPlayExternallyEpisodeClick && <ContextMenuItem
                                        onClick={() => onPlayExternallyEpisodeClick(episode)}
                                    >
                                        <LuTvMinimalPlay /> Play externally
                                    </ContextMenuItem>}
                                </>}
                            />
                        </CarouselItem>
                    ))}
                </CarouselContent>
            </Carousel>

            {shouldSplitEpisodeGridBySeason ? (
                <div className="space-y-8">
                    {episodeSeasonGroups.map(group => (
                        <section key={`season-${group.seasonNumber ?? "other"}`} className="space-y-3">
                            <div className="flex items-center justify-between gap-3">
                                <h3 className="text-sm font-semibold uppercase tracking-wide text-[--muted]">{group.label}</h3>
                                <span className="text-xs text-[--muted]">{group.episodes.length} episodes</span>
                            </div>
                            <EpisodeListGrid>
                                {group.episodes.map(episode => renderGridEpisodeItem(episode, `s${group.seasonNumber ?? "other"}-`))}
                            </EpisodeListGrid>
                        </section>
                    ))}
                </div>
            ) : (
                <EpisodeListPaginatedGrid
                    length={episodeCollection?.episodes?.length || 0}
                    shouldDefaultToPageWithEpisode={entry.listData?.progress ? entry.listData?.progress + 1 : undefined}
                    renderItem={(index) => {
                        const episode = episodeCollection?.episodes?.[index]
                        return renderGridEpisodeItem(episode)
                    }}
                />
            )}

            {bottomSection}
        </>
    )
}
