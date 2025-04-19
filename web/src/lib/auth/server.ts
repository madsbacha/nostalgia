import {cookies} from "next/headers";

export const AUTH_COOKIE_NAME = 'jwt'

export async function getAuthToken(): Promise<string | null> {
    const cookie_store = await cookies();
    const auth_session = cookie_store.get(AUTH_COOKIE_NAME)
    if (!auth_session) {
        return null
    }
    return auth_session.value
}
