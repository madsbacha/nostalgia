import Image from "next/image";
import {AspectRatio} from "@/components/ui/aspect-ratio";
import {useEffect, useState} from "react";
import {blurHashToDataURL} from "@/lib/blurhash/data_url";
import {PlaceholderValue} from "next/dist/shared/lib/get-img-props";

interface ThumbnailProps {
    src: string;
    alt: string;
    blurhash?: string;
}

export function Thumbnail({ src, alt, blurhash }: ThumbnailProps) {
    const [placeholder, setPlaceholder] = useState<PlaceholderValue | undefined>(undefined);
    const [error, setError] = useState<boolean>(false);

    useEffect(() => {
        if (blurhash === undefined) {
            return;
        }

        const blurhashDataUrl = blurHashToDataURL(blurhash);
        if (blurhashDataUrl === undefined) {
            return
        }
        setPlaceholder(`data:image/${blurhashDataUrl}`)
    }, [blurhash]);

    let actualSrc = src;
    if (error && placeholder !== undefined) {
        actualSrc = placeholder;
    }

    return <>
        <AspectRatio ratio={16 / 9}>
            <Image
                src={actualSrc}
                alt={alt}
                fill
                className="h-full w-full rounded-md object-cover"
                unoptimized
                placeholder={placeholder}
                onError={() => setError(true)}
            />
        </AspectRatio>
    </>
}
