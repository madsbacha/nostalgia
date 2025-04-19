import {NavbarLogo} from "@/components/navbar_logo";
import {UserNavigation} from "@/components/user_navigation";
import Link from "next/link";
import {getAuthToken} from "@/lib/auth/server";
import {getCurrentUser} from "@/lib/api";

export async function Navbar() {
    const auth_token = await getAuthToken(); // TODO: Can this be made more clean/maintainable?
    const user = auth_token != null ? await getCurrentUser(auth_token) : null;

    return (
        <>
            <div className="p-2 flex flex-row justify-between">
                <NavbarLogo />
                <div className="flex flex-row justify-between align-middle">
                    {user?.permissions.can_upload_media &&
                        <Link href="/upload" className="font-medium mr-4 content-center">Upload</Link>}
                    <UserNavigation />
                </div>
            </div>
        </>
    )
}