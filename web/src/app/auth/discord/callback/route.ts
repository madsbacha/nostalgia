import {getRedirectUri} from "@/lib/discord";
import {cookies} from "next/headers";
import {redirect} from "next/navigation";
import {AUTH_COOKIE_NAME} from "@/lib/auth/server";
import {apiAuthUrl, REDIRECT_AUTH_PERSIST_COOKIE} from "@/lib/endpoints";

export async function GET(request: Request) {
    const url = new URL(request.url);
    const code = url.searchParams.get('code');
    const error = url.searchParams.get('error');

    if (error) {
        // TODO: Handle error
        return new Response('Error: ' + error, { status: 400 });
    }

    if (!code) {
        // TODO: Handle error
        return new Response('Missing code', { status: 400 });
    }

    const redirect_uri = getRedirectUri();

    try {
        console.log(`ApiAuthUrl: ${apiAuthUrl()}`);
        const raw_response = await fetch(apiAuthUrl(), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ code, redirect_uri }),
        })
        const response = await raw_response.json();

        // TODO: Handle error when not authorized

        const isSecure = process.env.NODE_ENV === 'production' || process.env.FRONTEND_URL?.startsWith('https://')
        const cookieStore = await cookies()
        const maxAge = 60 * 60 * 24 * 7 // 1 week
        cookieStore.set(AUTH_COOKIE_NAME, response.token, { secure: isSecure, sameSite: "strict", maxAge: maxAge })
    } catch (e) {
        console.error(e);
        return new Response("Internal server error", { status: 500 });
    } finally {
        redirect(REDIRECT_AUTH_PERSIST_COOKIE)
    }
}