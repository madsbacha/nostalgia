"use server";

import {LoginButton} from "@/components/login_button";
import {isLoggedIn} from "@/lib/auth/server";
import {redirect} from "next/navigation";
import {ENDPOINT_PAGE_HOME, ENDPOINT_PAGE_LOGIN} from "@/lib/endpoints";
import {Redirect} from "@/components/redirect";

export default async function Page({
                                       searchParams,
                                   }: {
    searchParams: Promise<{ refresh: string | string[] | undefined }>
}) {
    const loggedIn = await isLoggedIn();
    if (loggedIn) {
        redirect(ENDPOINT_PAGE_HOME)
    }

    const { refresh } = await searchParams;

    return (
        <>
            <div className="h-dvh flex flex-col items-center justify-center bg-zinc-900">
                <div className="text-lg text-white font-medium mb-4">
                    klaphat
                </div>
                <div>
                    {refresh !== undefined ? <Redirect url={ENDPOINT_PAGE_LOGIN} /> : null}
                    <LoginButton/>
                </div>
            </div>
        </>
    )
}
