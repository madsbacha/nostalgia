export interface CurrentUser {
    id: string;
    username: string;
    avatar_url: string;
    permissions: UserPermissions;
}

export interface UserPermissions {
    can_view_media: boolean;
    can_upload_media: boolean;
    can_delete_media: boolean;
    can_update_media: boolean;
    whitelisted: boolean;
    can_manage_whitelist: boolean;
    can_manage_permissions: boolean;
}
