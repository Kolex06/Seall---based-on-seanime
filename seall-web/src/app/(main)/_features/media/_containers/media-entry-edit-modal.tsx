import { MediaAPI_BaseMedia, MediaAPI_BaseManga, MediaAPI_MediaListStatus, Media_EntryListData, Manga_EntryListData } from "@/api/generated/types"
import { useDeleteMediaListEntry, useEditMediaListEntry } from "@/api/hooks/simkl.hooks"
import { useUpdateMediaEntryRepeat } from "@/api/hooks/media_entries.hooks"
import { PluginWebviewSlot } from "@/app/(main)/_features/plugin/webview/plugin-webviews"
import { useCurrentUser } from "@/app/(main)/_hooks/use-server-status"
import { SeaImage } from "@/components/shared/sea-image"
import { Button, IconButton } from "@/components/ui/button"
import { cn } from "@/components/ui/core/styling"
import { Disclosure, DisclosureContent, DisclosureItem, DisclosureTrigger } from "@/components/ui/disclosure"
import { defineSchema, Field, Form } from "@/components/ui/form"
import { Modal, ModalProps } from "@/components/ui/modal"
import { NumberInput } from "@/components/ui/number-input"
import { Popover, PopoverProps } from "@/components/ui/popover"
import { Tooltip } from "@/components/ui/tooltip"
import { normalizeDate } from "@/lib/helpers/date"
import { getImageUrl } from "@/lib/server/assets"
import { useWindowSize } from "@uidotdev/usehooks"
import React, { Fragment } from "react"
import { BiListPlus, BiPlus, BiStar, BiTrash } from "react-icons/bi"
import { TbEdit } from "react-icons/tb"
import { useToggle } from "react-use"
import { z } from "zod"

type MediaEntryEditModalProps = {
    children?: React.ReactNode
    listData?: Media_EntryListData | Manga_EntryListData
    media?: MediaAPI_BaseMedia | MediaAPI_BaseManga
    hideButton?: boolean
    type?: "anime" | "manga"
    forceModal?: boolean
}

export const mediaListDataSchema = defineSchema(({ z, presets }) => z.object({
    status: z.custom<MediaAPI_MediaListStatus>().nullish(),
    score: z.number().min(0).max(100).nullish(),
    progress: z.number().min(0).nullish(),
    startedAt: presets.datePicker.nullish(),
    completedAt: presets.datePicker.nullish(),
}))

function IsomorphicPopover(props: PopoverProps & ModalProps & { media?: MediaAPI_BaseMedia | MediaAPI_BaseManga, forceModal?: boolean }) {
    const { title, children, media, forceModal, ...rest } = props
    const { width } = useWindowSize()

    if ((width && width > 1024) && !forceModal) {
        return <Popover
            {...rest}
            className="max-w-5xl !w-full overflow-hidden bg-gray-950/95 backdrop-blur-sm rounded-xl"
        >
            <p className="mb-4 font-semibold text-center px-6 line-clamp-1">
                {media?.title?.userPreferred}
            </p>
            {children}
        </Popover>
    }

    return <Modal
        {...rest}
        title={title}
        titleClass="text-xl"
        contentClass="max-w-3xl overflow-hidden"
    >
        {media?.bannerImage && <div
            data-media-entry-edit-modal-banner-image-container
            className="h-24 w-full flex-none object-cover object-center overflow-hidden absolute left-0 top-0 z-[0]"
        >
            <SeaImage
                data-media-entry-edit-modal-banner-image
                src={getImageUrl(media?.bannerImage!)}
                alt="banner"
                fill
                quality={80}
                sizes="20rem"
                className="object-cover object-center opacity-5 z-[1]"
            />
            <div
                data-media-entry-edit-modal-banner-image-bottom-gradient
                className="z-[5] absolute bottom-0 w-full h-[60%] bg-gradient-to-t from-[--background] to-transparent"
            />
        </div>}
        {children}
    </Modal>
}


export const MediaEntryEditModal = (props: MediaEntryEditModalProps) => {
    const [open, toggle] = useToggle(false)
    const [repeat, setRepeat] = React.useState(0)

    const { children, media, listData, hideButton, type = "anime", forceModal, ...rest } = props

    const user = useCurrentUser()

    const { mutate, isPending: isEditing, isSuccess, reset } = useEditMediaListEntry(media?.id, type)
    const { mutate: mutateRepeat, isPending: _isPending2 } = useUpdateMediaEntryRepeat(media?.id)
    const { mutate: deleteEntry, isPending: isDeleting } = useDeleteMediaListEntry(media?.id, type, () => {
        toggle(false)
    })

    React.useEffect(() => {
        setRepeat(listData?.repeat || 0)
    }, [listData])

    const handleSubmit = React.useCallback((data: z.infer<typeof mediaListDataSchema>) => {
        if (repeat !== (listData?.repeat ?? 0)) {
            mutateRepeat({
                mediaId: media?.id || 0,
                repeat: repeat,
            })
        }
        mutate({
            mediaId: media?.id || 0,
            status: data.status || "PLANNING",
            score: data.score ? data.score * 10 : 0, // should be 0-100
            progress: data.progress || 0,
            startedAt: data.startedAt ? {
                // @ts-ignore
                day: data.startedAt.getDate(),
                month: data.startedAt.getMonth() + 1,
                year: data.startedAt.getFullYear(),
            } : undefined,
            completedAt: data.completedAt ? {
                // @ts-ignore
                day: data.completedAt.getDate(),
                month: data.completedAt.getMonth() + 1,
                year: data.completedAt.getFullYear(),
            } : undefined,
            type: type,
        })
    }, [repeat, listData?.repeat, media?.id, type, mutate, mutateRepeat])


    if (!user) return null

    return (
        <>
            {!hideButton && <>
                {(!listData) && <Tooltip
                    trigger={<IconButton
                        data-media-entry-edit-modal-add-button
                        intent="gray-subtle"
                        icon={<BiPlus />}
                        rounded
                        size="sm"
                        loading={isEditing || isDeleting}
                        className={cn({ "hidden": isSuccess })} // Hide button when mutation is successful
                        onClick={() => mutate({
                            mediaId: media?.id || 0,
                            status: "PLANNING",
                            score: 0,
                            progress: 0,
                            startedAt: undefined,
                            completedAt: undefined,
                            type: type,
                        })}
                    />}
                >
                    Add to list
                </Tooltip>}
            </>}

            {!!listData && <IsomorphicPopover
                forceModal={forceModal}
                open={open}
                onOpenChange={o => toggle(o)}
                title={media?.title?.userPreferred ?? undefined}
                trigger={<span>
                    {!hideButton && <>
                        {!!listData && <IconButton
                            data-media-entry-edit-modal-edit-button
                            intent="white-subtle"
                            icon={<TbEdit />}
                            rounded
                            size="sm"
                            loading={isEditing || isDeleting}
                            onClick={toggle}
                        />}
                    </>}
                </span>}
                media={media}
            >

                {open && <Content
                    open={open}
                    onToggle={toggle}
                    repeat={repeat}
                    setRepeat={setRepeat}
                    handleSubmit={handleSubmit}
                    deleteEntry={deleteEntry}
                    isEditing={isEditing}
                    isDeleting={isDeleting}
                    {...props}
                />}

            </IsomorphicPopover>}
        </>
    )

}

function Content(props: MediaEntryEditModalProps & {
    open: boolean
    onToggle: (open: boolean) => void
    repeat: number
    setRepeat: (repeat: number) => void
    handleSubmit: (data: z.infer<typeof mediaListDataSchema>) => void
    deleteEntry: any
    isEditing: boolean
    isDeleting: boolean
}) {
    const {
        children,
        media,
        listData,
        hideButton,
        type = "anime",
        forceModal,
        open,
        onToggle,
        repeat,
        setRepeat,
        handleSubmit,
        deleteEntry,
        isEditing,
        isDeleting,
        ...rest
    } = props

    return (
        <>
            {(!!listData) && <Form
                data-media-entry-edit-modal-form
                schema={mediaListDataSchema}
                onSubmit={handleSubmit}
                className={cn(
                    // {
                    //     "mt-8": !!media?.bannerImage,
                    // },
                )}
                onError={console.log}
                defaultValues={{
                    status: listData?.status,
                    score: listData?.score ? listData?.score / 10 : undefined, // Returned score is 0-100
                    progress: listData?.progress,
                    startedAt: listData?.startedAt ? (normalizeDate(listData?.startedAt)) : undefined,
                    completedAt: listData?.completedAt ? (normalizeDate(listData?.completedAt)) : undefined,
                }}
            >
                <div className="flex flex-col sm:flex-row gap-4">
                    <Field.Select
                        label="Status"
                        name="status"
                        options={[
                            media?.status !== "NOT_YET_RELEASED" ? {
                                value: "CURRENT",
                                label: type === "anime" ? "Watching" : "Reading",
                            } : undefined,
                            { value: "PLANNING", label: "Planning" },
                            media?.status !== "NOT_YET_RELEASED" ? {
                                value: "PAUSED",
                                label: "Paused",
                            } : undefined,
                            media?.status !== "NOT_YET_RELEASED" ? {
                                value: "COMPLETED",
                                label: "Completed",
                            } : undefined,
                            media?.status !== "NOT_YET_RELEASED" ? {
                                value: "DROPPED",
                                label: "Dropped",
                            } : undefined,
                            media?.status !== "NOT_YET_RELEASED" ? {
                                value: "REPEATING",
                                label: "Repeating",
                            } : undefined,
                        ].filter(Boolean)}
                    />
                    {media?.status !== "NOT_YET_RELEASED" && <>
                        <Field.Number
                            label="Score"
                            name="score"
                            min={0}
                            max={10}
                            formatOptions={{
                                maximumFractionDigits: 1,
                                minimumFractionDigits: 0,
                                useGrouping: false,
                            }}
                            rightIcon={<BiStar />}
                        />
                        <Field.Number
                            label="Progress"
                            name="progress"
                            min={0}
                            max={type === "anime" ? (!!(media as MediaAPI_BaseMedia)?.nextAiringEpisode?.episode
                                ? (media as MediaAPI_BaseMedia)?.nextAiringEpisode?.episode! - 1
                                : ((media as MediaAPI_BaseMedia)?.episodes
                                    ? (media as MediaAPI_BaseMedia).episodes
                                    : undefined)) : (media as MediaAPI_BaseManga)?.chapters}
                            formatOptions={{
                                maximumFractionDigits: 0,
                                minimumFractionDigits: 0,
                                useGrouping: false,
                            }}
                            rightIcon={<BiListPlus />}
                        />
                    </>}
                </div>
                {media?.status !== "NOT_YET_RELEASED" && <div className="flex flex-col sm:flex-row gap-4">
                    <Field.DatePicker
                        label="Start date"
                        name="startedAt"
                        // defaultValue={(state.startedAt && state.startedAt.year) ? parseAbsoluteToLocal(new Date(state.startedAt.year,
                        // (state.startedAt.month || 1)-1, state.startedAt.day || 1).toISOString()) : undefined}
                    />
                    <Field.DatePicker
                        label="Completion date"
                        name="completedAt"
                        // defaultValue={(state.completedAt && state.completedAt.year) ? parseAbsoluteToLocal(new Date(state.completedAt.year,
                        // (state.completedAt.month || 1)-1, state.completedAt.day || 1).toISOString()) : undefined}
                    />

                    <NumberInput
                        name="repeat"
                        label={type === "anime" ? "Total rewatches" : "Total rereads"}
                        min={0}
                        max={1000}
                        value={repeat}
                        onValueChange={setRepeat}
                        formatOptions={{
                            maximumFractionDigits: 0,
                            minimumFractionDigits: 0,
                            useGrouping: false,
                        }}
                    />
                </div>}

                <div className="flex w-full items-center justify-between mt-4">
                    <div>
                        <Disclosure type="multiple" defaultValue={["item-2"]}>
                            <DisclosureItem value="item-1" className="flex items-center gap-1">
                                <DisclosureTrigger>
                                    <IconButton
                                        intent="alert-subtle"
                                        icon={<BiTrash />}
                                        rounded
                                        size="md"
                                    />
                                </DisclosureTrigger>
                                <DisclosureContent>
                                    <Button
                                        intent="alert-basic"
                                        rounded
                                        size="md"
                                        loading={isDeleting}
                                        onClick={() => deleteEntry({
                                            mediaId: media?.id!,
                                            type: type,
                                        })}
                                    >Confirm</Button>
                                </DisclosureContent>
                            </DisclosureItem>
                        </Disclosure>
                    </div>

                    <Field.Submit role="save" disableIfInvalid={true} loading={isEditing} disabled={isDeleting}>
                        Save
                    </Field.Submit>
                </div>

                <PluginWebviewSlot slot="after-media-entry-form" />
            </Form>}
        </>
    )
}
