"use server";

import {Navbar} from "@/components/navbar";
import {WhitelistUsers} from "@/components/whitelist_users";
import {getAuthToken, redirectIfUnauthorized} from "@/lib/auth/server";
import {redirect} from "next/navigation";
import {ENDPOINT_PAGE_LOGIN} from "@/lib/endpoints";


export default async function Home() {
    await redirectIfUnauthorized();
    const auth_token = await getAuthToken();
    if (auth_token === null) {
        redirect(ENDPOINT_PAGE_LOGIN)
    }

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
