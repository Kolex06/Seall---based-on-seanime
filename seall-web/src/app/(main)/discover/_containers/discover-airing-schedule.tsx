import { useGetMediaCollectionSchedule } from "@/api/hooks/media_library.hooks"
import { SeaContextMenu } from "@/app/(main)/_features/context-menu/sea-context-menu"
import { useMediaPreviewModal } from "@/app/(main)/_features/media/_containers/media-preview-modal"
import { SeaImage } from "@/components/shared/sea-image"
import { SeaLink } from "@/components/shared/sea-link"
import { ContextMenuGroup, ContextMenuItem, ContextMenuLabel, ContextMenuTrigger } from "@/components/ui/context-menu"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Separator } from "@/components/ui/separator"
import { useRouter } from "@/lib/navigation"
import { format, isSameMonth, isToday, startOfDay } from "date-fns"
import { addDays } from "date-fns/addDays"
import { isSameDay } from "date-fns/isSameDay"
import React from "react"
import { LuDock, LuEye } from "react-icons/lu"


export function DiscoverAiringSchedule() {
    const { data, isLoading } = useGetMediaCollectionSchedule({ source: "all" })
    const scheduleStart = React.useMemo(() => startOfDay(new Date()), [])
    const scheduleEnd = React.useMemo(() => addDays(scheduleStart, 14), [scheduleStart])

    const media = React.useMemo(() => (data ?? []).filter(item => {
        if (!item?.dateTime) return false
        const dateTime = new Date(item.dateTime)
        return dateTime >= scheduleStart && dateTime <= scheduleEnd
    }), [data, scheduleEnd, scheduleStart])

    const router = useRouter()
    const { setPreviewModalMediaId } = useMediaPreviewModal()

    const currentDate = scheduleStart

    const days = React.useMemo(() => {

        const daysArray = []
        let day = scheduleStart

        while (day <= scheduleEnd) {
            const upcomingMedia = media.filter((item) => !!item?.dateTime && isSameDay(new Date(item.dateTime), day)).map((item) => {
                return {
                    id: `${item.mediaId}-${item.episodeNumber}-${item.dateTime}`,
                    mediaId: item.mediaId,
                    name: item.title,
                    time: format(new Date(item.dateTime!), "h:mm a"),
                    datetime: format(new Date(item.dateTime!), "yyyy-MM-dd'T'HH:mm"),
                    href: `/entry?id=${item.mediaId}`,
                    image: item.image,
                    episode: item.episodeNumber || 1,
                    isMovie: item.isMovie,
                }
            })

            daysArray.push({
                date: format(day, "yyyy-MM-dd'T'HH:mm"),
                isCurrentMonth: isSameMonth(day, currentDate),
                isToday: isToday(day),
                isSelected: false,
                events: upcomingMedia,
            })
            day = addDays(day, 1)
        }
        return daysArray
    }, [media, currentDate, scheduleEnd, scheduleStart])

    if (isLoading) return <LoadingSpinner />

    if (!media.length) return null

    return (
        <div className="space-y-4 z-[5] relative" data-discover-airing-schedule-container>
            <h2 className="text-center">Airing Schedule</h2>
            <div className="space-y-6">
                {days.map((day, index) => {
                    if (day.events.length === 0) return null
                    return (
                        <React.Fragment key={day.date}>
                            <div className="flex flex-col gap-2">
                                <div className="flex items-center gap-2">
                                    <h3 className="font-semibold">{format(new Date(day.date), "EEEE, PP")}</h3>
                                    {day.isToday && <span className="text-[--muted]">Today</span>}
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
                                    {day.events?.toSorted((a, b) => a.datetime.localeCompare(b.datetime))?.map((event, index) => {
                                        return (
                                            <React.Fragment key={event.id}>
                                                <SeaContextMenu
                                                    content={<ContextMenuGroup>
                                                        <ContextMenuLabel className="text-[--muted] line-clamp-2 py-0 my-2">
                                                            {event.name}
                                                        </ContextMenuLabel>
                                                        <ContextMenuItem
                                                            onClick={() => {
                                                                setPreviewModalMediaId(event.mediaId || 0, "anime")
                                                            }}
                                                        >
                                                            <LuEye /> Preview
                                                        </ContextMenuItem>
                                                        <ContextMenuItem
                                                            onClick={() => {
                                                                router.push(`/entry?id=${event.mediaId}`)
                                                            }}
                                                        >
                                                            <LuDock /> Open page
                                                        </ContextMenuItem>
                                                    </ContextMenuGroup>}
                                                >
                                                    <ContextMenuTrigger>
                                                        <div
                                                            key={String(`${event.id}${index}`)}
                                                            className="flex gap-3 bg-[--background] rounded-[--radius-md] p-2"
                                                        >
                                                            <div
                                                                className="w-[5rem] h-[5rem] rounded-[--radius] flex-none object-cover object-center overflow-hidden relative"
                                                            >
                                                                <SeaImage
                                                                    src={event.image || "/no-cover.png"}
                                                                    alt="banner"
                                                                    fill
                                                                    quality={80}
                                                                    priority
                                                                    sizes="20rem"
                                                                    className="object-cover object-center"
                                                                />
                                                            </div>

                                                            <div className="space-y-1">
                                                                <SeaLink
                                                                    href={event.href}
                                                                    className="font-medium tracking-wide line-clamp-1"
                                                                >{event.name}</SeaLink>

                                                                <p className="text-[--muted]">
                                                                    {event.isMovie ? "Movie" : `Ep ${event.episode}`} at {event.time}
                                                                </p>
                                                            </div>
                                                        </div>
                                                    </ContextMenuTrigger>
                                                </SeaContextMenu>
                                            </React.Fragment>
                                        )
                                    })}
                                </div>
                            </div>
                            {!!days[index + 1]?.events?.length && <Separator />}
                        </React.Fragment>
                    )
                })}
            </div>
        </div>
    )
}
