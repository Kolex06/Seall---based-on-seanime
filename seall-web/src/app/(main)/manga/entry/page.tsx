import { useRouter } from "@/lib/navigation"
import React from "react"

export default function Page() {
    const router = useRouter()

    React.useEffect(() => {
        router.replace("/lists")
    }, [router])

    return null
}
