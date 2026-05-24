import { NavigationMenu, NavigationMenuProps } from "@/components/ui/navigation-menu"
import { usePathname } from "@/lib/navigation"
import React, { useMemo } from "react"

interface OfflineTopMenuProps {
    children?: React.ReactNode
}

export const OfflineTopMenu: React.FC<OfflineTopMenuProps> = (props) => {

    const { children, ...rest } = props

    const pathname = usePathname()

    const navigationItems = useMemo<NavigationMenuProps["items"]>(() => {

        return [
            {
                href: "/offline",
                // icon: IoLibrary,
                isCurrent: pathname === "/offline",
                name: "Media Library",
            },
        ].filter(Boolean)
    }, [pathname])

    return (
        <NavigationMenu
            className="p-0 hidden lg:inline-block"
            itemClass="text-xl"
            items={navigationItems}
        />
    )

}
