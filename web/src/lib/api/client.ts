import useSWR, {Fetcher} from "swr";
import {
    apiAllTagsUrl, apiCurrentUserUrl,
    apiMediaListUrl,
    apiMediaUrl,
    apiUserPermissionsUrl,
    apiUsersForWhitelistUrl
} from "@/lib/endpoints";
import useSWRMutation from "swr/mutation";
import {CurrentUser, UserPermissions} from "@/lib/api/types";

interface Media {
    id: string;
    title: string;
    description: string;
    source: string;
    mimeType: string;
}

export function useMedia(mediaId: string) {
    const fetcher: Fetcher<Media, string> = (...args) => {
        return fetch(...args, {credentials: "include"}).then(res => res.json())
    };

    const { data, isLoading, error } = useSWR(apiMediaUrl(mediaId), fetcher);

    return {
        media: data,
        isLoading,
        isError: error,
    };
}

interface MediaPreview {
    id: string;
    title: string;
    thumbnail_url: string;
    thumbnail_blurhash: string;
    tags: string[];
}

interface MediaList {
    media_list: MediaPreview[];
}


export function useMediaList() {
    const fetcher: Fetcher<MediaList, string> = (...args) => {
        return fetch(...args, {credentials:"include"}).then(res => res.json())
    };

    const { data, isLoading, error } = useSWR(apiMediaListUrl(), fetcher);

    return {
        mediaList: data,
        isLoading,
        isError: error,
    };
}

interface AllTagsResponse {
    tags: string[];
}

export function useTags() {
    const fetcher: Fetcher<AllTagsResponse, string> = (...args) => {
        return fetch(...args, {credentials:"include"}).then(res => res.json())
    };
    const { data, isLoading, error } = useSWR(apiAllTagsUrl(), fetcher);

    return {
        tags: data,
        isLoading,
        isError: error,
    };
}

type UsersForWhitelistResponse = WhitelistUser[];

export interface WhitelistUser {
    id: string;
    username: string;
    avatar_url: string;
    is_whitelisted: boolean;
}

export function useUsersForWhitelist() {
    const fetcher: Fetcher<UsersForWhitelistResponse, string> = (...args) => {
        return fetch(...args, {credentials: "include"}).then(res => res.json())
    };
    const { data, isLoading, error } = useSWR(apiUsersForWhitelistUrl(), fetcher)

    return {
        users: data,
        isLoading,
        isError: error
    };
}

export function usePermissions(userId: string) {
    const fetcher: Fetcher<UserPermissions, string> = (...args) => {
        return fetch(...args, {credentials: "include"}).then(res => res.json())
    };

    const { data, isLoading, error } = useSWR(apiUserPermissionsUrl(userId), fetcher)

    return {
        permissions: data,
        isLoading,
        isError: error
    };
}

async function updatePermissions(url: string, { arg }: { arg: UserPermissions }) {
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(arg),
        credentials: "include"
    });

    if (!response.ok) {
        throw new Error("Failed to update permissions");
    }

    return response.json();
}

export function usePermissionsMutation(userId: string) {
    const { trigger, isMutating } = useSWRMutation(apiUserPermissionsUrl(userId), updatePermissions);

    return {
        trigger,
        isMutating,
    }
}


export function useCurrentUser() {
    const fetcher: Fetcher<CurrentUser, string> = (...args) => {
        return fetch(...args, {credentials: "include"}).then(res => res.json())
    };

    const { data, isLoading, error } = useSWR(apiCurrentUserUrl(), fetcher)

    return {
        user: data,
        isLoading,
        isError: error
    };
}


