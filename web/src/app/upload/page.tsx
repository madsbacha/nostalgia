"use server";

import {getAuthToken} from "@/lib/auth/server";
import {redirect} from "next/navigation";
import {ENDPOINT_PAGE_HOME} from "@/lib/endpoints";
import {Navbar} from "@/components/navbar";
import {UploadMediaForm} from "@/components/upload_media_form";
import {getCurrentUser} from "@/lib/api";

export default async function Home() {
    const auth_token = await getAuthToken();
    const user = auth_token != null ? await getCurrentUser(auth_token) : null;
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
