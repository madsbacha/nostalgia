import { Skeleton } from "@/components/ui/skeleton"
import {AspectRatio} from "@/components/ui/aspect-ratio";

export function SkeletonVideoCard() {
    return (
        <div className="flex flex-col space-y-3">
            <div className="w-full">
                <AspectRatio ratio={16 / 9}>
                    <Skeleton className="h-full w-full rounded-xl" />
                </AspectRatio>
            </div>
            <div className="space-y-2">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-4 w-3/5" />
            </div>
        </div>
    )
}
