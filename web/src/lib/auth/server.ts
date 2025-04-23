import {cookies} from "next/headers";
import {getCurrentUser} from "@/lib/api";
import {ENDPOINT_PAGE_LOGIN} from "@/lib/endpoints";
import {redirect} from "next/navigation";

export const AUTH_COOKIE_NAME = 'jwt'

export async function getAuthToken(): Promise<string | null> {
    const cookie_store = await cookies();
    const auth_session = cookie_store.get(AUTH_COOKIE_NAME)
    if (!auth_session) {
        return null
    }
    return auth_session.value
}

export async function isLoggedIn() {
    const auth_token = await getAuthToken();
    if (auth_token === null) {
        return false;
    }
    try {
        const user = await getCurrentUser(auth_token);
        return user !== null;
    } catch (e) {
        return e instanceof Error && e.message === "Unauthorized";
    }
}

export async function redirectIfUnauthorized() {
    const loggedIn = await isLoggedIn();
    if (!loggedIn) {
        redirect(ENDPOINT_PAGE_LOGIN)
    }
}
