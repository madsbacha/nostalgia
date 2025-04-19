"use server";

import {Navbar} from "@/components/navbar";
import {WhitelistUsers} from "@/components/whitelist_users";


export default async function Home() {
    return (
        <>
            <div className="flex flex-col items-center justify-center">
                <div className="container h-dvh">
                    <Navbar />
                    <main>
                        <WhitelistUsers />
                    </main>
                </div>
            </div>
        </>
    )
}
