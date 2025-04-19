import {getDiscordAuthInfo} from "@/lib/api";

export function getRedirectUri(): string {
    const frontendUrl = process.env.FRONTEND_URL;
    if (frontendUrl === undefined) {
        throw new Error('Missing FRONTEND_URL environment variable');
    }
    return `${frontendUrl}/auth/discord/callback`;
}

export async function getDiscordUrl(): Promise<string> {
    const info = await getDiscordAuthInfo();
    const redirectUri = getRedirectUri();
    const url = createDiscordURL(
        info.client_id,
        redirectUri,
        info.response_type,
        info.scopes
    )
    return url.toString();
}

function createDiscordURL(clientId: string, redirectUri: string, responseType: string, scope: string) {
    const url = new URL('https://discord.com/api/oauth2/authorize');
    url.searchParams.set('client_id', clientId);
    url.searchParams.set('redirect_uri', redirectUri);
    url.searchParams.set('response_type', responseType);
    url.searchParams.set('scope', scope);
    return url;
}
