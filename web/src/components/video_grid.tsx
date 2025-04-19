"use client";

import {SkeletonVideoCard} from "@/components/skeleton_video_card";
import {useMediaList, useTags} from "@/lib/api/client";
import {AspectRatio} from "@/components/ui/aspect-ratio";
import Image from "next/image"
import Link from "next/link";
import {MultiSelect} from "@/components/multi-select";
import {useState} from "react";
import {Skeleton} from "@/components/ui/skeleton";

export function VideoGrid() {
    const { mediaList, isLoading, isError } = useMediaList();
    const { tags, isLoading: isLoadingTags, isError: isErrorTags } = useTags();


    const [selectedTags, setselectedTags] = useState<string[]>([]);

    let multiSelect = <Skeleton className="w-full h-10 my-4" />;
    if (isErrorTags) {
        multiSelect = <></>
    } else if (!isLoadingTags && tags) {
        const options = tags.tags.map(tag => ({
            value: tag,
            label: tag,
        }));
        multiSelect = (
            <MultiSelect
                options={options}
                onValueChange={setselectedTags}
                defaultValue={selectedTags}
                placeholder="Select Tags"
                variant="inverted"
                className="w-full my-4"
            />
        );
    }

    console.log(mediaList);

    const videos = []
    if (isLoading) {
        for (let i = 0; i < 20; i++) {
            videos.push(<SkeletonVideoCard key={i} />)
        }
    } else if (isError || !mediaList || !mediaList.media_list) {
        videos.push(<div>Failed to load</div>)
    } else {
        const filteredMediaList = mediaList?.media_list.filter(media => {
            if (selectedTags.length === 0) return true;
            return selectedTags.every(tag => media.tags.includes(tag));
        });
        for (const media of filteredMediaList) {
            videos.push(
                <div className="flex flex-col space-y-3" key={media.id}>
                    <Link href={`/media/${media.id}`}>
                        <div className="w-full mb-1">
                            <AspectRatio ratio={16 / 9}>
                                <Image
                                    src={media.thumbnail_url}
                                    alt={media.title}
                                    fill
                                    className="h-full w-full rounded-md object-cover"
                                    unoptimized
                                />
                            </AspectRatio>
                        </div>
                        <h3>{media.title}</h3>
                    </Link>
                </div>
            )
        }
    }

    return (
        <>
            {multiSelect}
            <div className="grid xl:grid-cols-6 md:grid-cols-4 sm:grid-cols-3 grid-cols-2 gap-4 w-full">
                {...videos}
            </div>
        </>
    )
}