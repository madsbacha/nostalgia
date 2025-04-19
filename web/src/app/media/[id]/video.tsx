"use client";

import {useParams} from "next/navigation";
import {VideoLarge} from "@/components/video/large";

export function Video() {
    const params = useParams<{ id: string }>()

    return (
        <>
            <VideoLarge mediaId={params.id} />
        </>
    )
}
