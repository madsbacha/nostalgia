import {useCurrentUser, usePermissions, usePermissionsMutation, WhitelistUser} from "@/lib/api/client";
import {TableCell, TableRow} from "@/components/ui/table";
import {UserAvatar} from "@/components/user_avatar";
import {Switch} from "@/components/ui/switch";
import {SkeletonSwitch} from "@/components/skeleton_switch";
import {toast} from "sonner";
import {UserPermissions} from "@/lib/api/types";

interface PermissionsTableRowProps {
    user: WhitelistUser
}

const AddMsg = (username: string) => `Added permission to ${username}.`
const RemoveMsg = (username: string) => `Removed permission from ${username}.`
const FailedMsg = (username: string) => `Failed to update ${username}'s permissions.`

export function PermissionsTableRow({ user }: PermissionsTableRowProps) {
    const { user: currentUser } = useCurrentUser();
    const { trigger } = usePermissionsMutation(user.id)
    const { permissions, isLoading } = usePermissions(user.id)
    const showSkeleton = isLoading || permissions === undefined

    function toggleProperty(change:  (p: UserPermissions, value: boolean) => UserPermissions) {
        return async function handleToggle(value: boolean) {
            if (isLoading || permissions === undefined) return
            const p = change(permissions, value)
            try {
                await trigger(p)
                toast(value ? AddMsg(user.username) : RemoveMsg(user.username))
            } catch (e) {
                console.error(e)
                toast.error(FailedMsg(user.username))
            }
        }
    }
    const isRowForCurrentUser = currentUser?.id == user.id

    return (
        <TableRow key={user.id}>
            <TableCell className="font-medium">
                <div className="flex flex-row items-center">
                    <UserAvatar username={user.username} avatarUrl={user.avatar_url} />
                    <p className="ml-4">{user.username}</p>
                </div>
            </TableCell>
            <TableCell className="text-right">
                {showSkeleton ? <SkeletonSwitch /> : <Switch checked={permissions.can_manage_permissions} disabled={isRowForCurrentUser} onCheckedChange={toggleProperty((p, v) => ({...p, can_manage_permissions: v}))} />}
            </TableCell>
            <TableCell className="text-right">
                {showSkeleton ? <SkeletonSwitch /> : <Switch checked={permissions.can_upload_media} onCheckedChange={toggleProperty((p, v) => ({...p, can_upload_media: v}))} />}
            </TableCell>
            <TableCell className="text-right">
                {showSkeleton ? <SkeletonSwitch /> : <Switch checked={permissions.can_view_media} onCheckedChange={toggleProperty((p, v) => ({...p, can_view_media: v}))} />}
            </TableCell>
            <TableCell className="text-right">
                {showSkeleton ? <SkeletonSwitch /> : <Switch checked={permissions.whitelisted} disabled={isRowForCurrentUser} onCheckedChange={toggleProperty((p, v) => ({...p, whitelisted: v}))} />}
            </TableCell>
        </TableRow>
    )
}