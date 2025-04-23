
export const ENDPOINT_PAGE_HOME = "/";
export const ENDPOINT_PAGE_LOGIN = "/login";
export const ENDPOINT_PAGE_LOGOUT = "/logout";

export const REDIRECT_UNAUTHORIZED = ENDPOINT_PAGE_LOGIN;
export const REDIRECT_AUTHORIZED = ENDPOINT_PAGE_HOME;
export const REDIRECT_AUTH_PERSIST_COOKIE = `${ENDPOINT_PAGE_LOGIN}?refresh`;

export const API_URL = () => {
    let url = process.env.NEXT_PUBLIC_API_URL;
    if (typeof window === 'undefined') {
        url = process.env.SSR_API_URL;
    }
    if (url === undefined) {
        throw new Error(`Missing API URL. ${typeof window !== 'undefined'}, SSR_API_URL: ${process.env.SSR_API_URL}, NEXT_PUBLIC_API_URL: ${process.env.NEXT_PUBLIC_API_URL}`);
    }
    return url;
}

export const apiAuthUrl = () => `${API_URL()}/api/auth/discord`;
export const apiGetInfoUrl = () => `${API_URL()}/api/info`;
export const apiDiscordAuthInfoUrl = () => `${API_URL()}/api/auth/discord`;
export const apiCurrentUserUrl = () => `${API_URL()}/api/users/@me`;
export const apiUploadMediaUrl = () => `${API_URL()}/api/media`;
export const apiMediaUrl = (mediaId: string) => `${API_URL()}/api/media/${mediaId}`;
export const apiMediaListUrl = () => `${API_URL()}/api/media`;
export const apiAllTagsUrl = () => `${API_URL()}/api/media/tags`;
export const apiUsersForWhitelistUrl = () => `${API_URL()}/api/users`;
export const apiUserPermissionsUrl = (userId: string) => `${API_URL()}/api/users/${userId}/permissions`;
