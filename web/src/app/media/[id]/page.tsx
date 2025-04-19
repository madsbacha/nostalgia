"use server";

import {Navbar} from "@/components/navbar";
import {Video} from "@/app/media/[id]/video";

export default async function Home() {
    return (
        <>
            <div className="flex flex-col items-center justify-center">
                <div className="container h-dvh ">
                    <Navbar />
                    <main>
                        <div className="flex flex-col space-y-3">
                            <Video />
                        </div>
                    </main>
                </div>
            </div>
        </>
    )
}
