"use server";

import {getAuthToken, redirectIfUnauthorized} from "@/lib/auth/server";
import {redirect} from "next/navigation";
import {ENDPOINT_PAGE_HOME, ENDPOINT_PAGE_LOGIN} from "@/lib/endpoints";
import {Navbar} from "@/components/navbar";
import {UploadMediaForm} from "@/components/upload_media_form";
import {getCurrentUser} from "@/lib/api";

export default async function Home() {
    await redirectIfUnauthorized();
    const auth_token = await getAuthToken();
    if (auth_token === null) {
        redirect(ENDPOINT_PAGE_LOGIN)
    }

    const user = await getCurrentUser(auth_token);
    if (!user?.permissions.can_upload_media) {
        redirect(ENDPOINT_PAGE_HOME);
    }

    return (
        <>
            <div className="flex flex-col items-center justify-center">
                <div className="container h-dvh">
                    <Navbar />
                    <main>
                        <div className="flex min-h-svh w-full justify-center p-6 md:p-10">
                            <div className="w-full max-w-sm">
                                <UploadMediaForm />
                            </div>
                        </div>
                    </main>
                </div>
            </div>
        </>
    )
}
