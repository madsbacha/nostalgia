"use client"

import {AspectRatio} from "@/components/ui/aspect-ratio";
import {useMedia} from "@/lib/api/client";
import {SkeletonVideoCard} from "@/components/skeleton_video_card";

export function VideoLarge({ mediaId }: { mediaId: string }) {
    const { media, isLoading, isError } = useMedia(mediaId);

    if (isLoading) {
        return <SkeletonVideoCard />
    }
    if (isError || !media) {
        return <div>Failed to load</div>
    }

    return (
        <>
            <div className="w-full">
                <AspectRatio ratio={16 / 9}>
                    <video
                        id="player"
                        playsInline
                        autoPlay
                        controls
                        className="h-full w-full rounded-md"
                    >
                        <source src={media.source} type={media.mimeType} />
                    </video>
                </AspectRatio>
            </div>
            <div className="space-y-2">
                <h1 className="text-xl">{media.title}</h1>
                <p className="mb-4">{media.description}</p>
            </div>
        </>
    );
}
