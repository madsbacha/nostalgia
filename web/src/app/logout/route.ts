import {cookies} from "next/headers";
import {redirect} from "next/navigation";
import {REDIRECT_UNAUTHORIZED} from "@/lib/endpoints";
import {AUTH_COOKIE_NAME} from "@/lib/auth/server";

export async function GET() {
    const cookieStore = await cookies()
    cookieStore.delete(AUTH_COOKIE_NAME)
    redirect(REDIRECT_UNAUTHORIZED)
}