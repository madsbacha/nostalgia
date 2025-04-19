import {apiCurrentUserUrl, apiDiscordAuthInfoUrl, apiGetInfoUrl, apiUploadMediaUrl} from "@/lib/endpoints";
import {CurrentUser} from "@/lib/api/types";

export function ApplyAuthHeaders(auth_token: string, headers: HeadersInit | undefined = undefined) : HeadersInit | undefined {
    let new_headers = headers || {}
    new_headers = Object.assign({}, new_headers, {
        Authorization: `Bearer ${auth_token}`
    })
    return new_headers
}

export async function getCurrentUser(token: string) : Promise<CurrentUser> {
    const userResponse = await fetch(apiCurrentUserUrl(),
        {
            headers: ApplyAuthHeaders(token)
        })
    const user = await userResponse.json();

    return {
        id: user.id,
        username: user.username,
        avatar_url: user.avatar_url,
        permissions: user.permissions,
    }
}

export interface Info {
    title: string;
}

export async function getInfo() : Promise<Info> {
    const res = await fetch(apiGetInfoUrl())
    const info = await res.json();

    return {
        title: info.title,
    }
}

export interface DiscordAuthInfo {
    client_id: string;
    response_type: string;
    scopes: string;
}

export async function getDiscordAuthInfo() : Promise<DiscordAuthInfo> {
    const res = await fetch(apiDiscordAuthInfoUrl())
    const info = await res.json();

    return {
        client_id: info.client_id,
        response_type: info.response_type,
        scopes: info.scopes,
    }
}

export async function uploadMediaFile(file: File, title: string, description: string, tags: string) {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("title", title);
    formData.append("description", description);
    formData.append("tags", tags);

    const response = await fetch(apiUploadMediaUrl(), {
        method: "POST",
        credentials: "include",
        body: formData,
    });

    if (!response.ok) {
        throw new Error("Failed to upload media");
    }

    return response.json();
}
