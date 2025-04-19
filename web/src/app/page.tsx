"use server";

import {Navbar} from "@/components/navbar";
import {VideoGrid} from "@/components/video_grid";


export default async function Home() {
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
