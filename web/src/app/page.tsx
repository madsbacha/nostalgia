"use server";

import {Navbar} from "@/components/navbar";
import {VideoGrid} from "@/components/video_grid";
import {redirectIfUnauthorized} from "@/lib/auth/server";


export default async function Home() {
    await redirectIfUnauthorized();

    return (
        <>
            <div className="flex flex-col items-center justify-center">
                <div className="container h-dvh ">
                    <Navbar />
                    <main>
                        <VideoGrid />
                    </main>
                </div>
            </div>
        </>
    )
}
