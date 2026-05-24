import { NavigationMenu, NavigationMenuProps } from "@/components/ui/navigation-menu"
import { usePathname } from "@/lib/navigation"
import React, { useMemo } from "react"

interface TopMenuProps {
    children?: React.ReactNode
}

export const TopMenu: React.FC<TopMenuProps> = (props) => {

    const { children, ...rest } = props

    const pathname = usePathname()

    const navigationItems = useMemo<NavigationMenuProps["items"]>(() => {

        return [
            {
                href: "/",
                // icon: IoLibrary,
                isCurrent: pathname === "/",
                name: "Home",
            },
            {
                href: "/schedule",
                icon: null,
                isCurrent: pathname.startsWith("/schedule"),
                name: "Schedule",
                // addon: missingEpisodeCount > 0 ? <Badge
                //     className="absolute -top-1 right-2 h-2 w-2 p-0 z-[5]" size="sm"
                //     intent="alert-solid"
                // /> : undefined,
            },
            {
                href: "/lists",
                icon: null,
                isCurrent: pathname.startsWith("/lists"),
                name: "Watchlist",
            },
            {
                href: "/discover",
                icon: null,
                isCurrent: pathname.startsWith("/discover") || pathname.startsWith("/search"),
                name: "Discover",
            },
        ].filter(Boolean)
    }, [pathname])

    return (
        <NavigationMenu
            className="p-0 hidden lg:inline-block"
            items={navigationItems}
            desktopListClass="space-x-0"
            data-top-menu
        />
    )

}
