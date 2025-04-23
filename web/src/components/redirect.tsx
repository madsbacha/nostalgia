"use client"

import {useRouter} from "next/navigation";

export function Redirect({ url }: { url: string }) {
    const router = useRouter()
    router.push(url)
    return <>
    </>
}
