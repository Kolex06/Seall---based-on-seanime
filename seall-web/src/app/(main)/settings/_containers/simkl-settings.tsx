import { useLocalSyncSimulatedDataToMedia } from "@/api/hooks/local.hooks"
import { SettingsPageHeader } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import React from "react"
import { SiSimkl } from "react-icons/si"

type Props = {
    isPending: boolean
    children?: React.ReactNode
}

export function SimklSettings(props: Props) {

    const {
        isPending,
        children,
        ...rest
    } = props

    const { mutate: upload, isPending: isUploading } = useLocalSyncSimulatedDataToMedia()

    const confirmDialog = useConfirmationDialog({
        title: "Upload to SIMKL",
        description: "This will upload your local Seall collection to your SIMKL account. Are you sure you want to proceed?",
        actionText: "Upload",
        actionIntent: "primary",
        onConfirm: async () => {
            if (isUploading) return
            upload()
        },
    })

    return (
        <div className="space-y-4">

            <SettingsPageHeader
                title="SIMKL"
                description="Manage your SIMKL account"
                icon={SiSimkl}
            />


            <SettingsSubmitButton isPending={isPending} />

            <ConfirmationDialog {...confirmDialog} />

        </div>
    )
}
