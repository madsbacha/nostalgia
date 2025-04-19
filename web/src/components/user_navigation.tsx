import {
    DropdownMenu,
    DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import {Button} from "@/components/ui/button";
import {getCurrentUser} from "@/lib/api";
import Link from "next/link";
import {getAuthToken} from "@/lib/auth/server";
import {UserAvatar} from "@/components/user_avatar";


export async function UserNavigation() {
    const auth_token = await getAuthToken();
    const user = auth_token != null ? await getCurrentUser(auth_token) : null;

    const canManagePermissions = user?.permissions.can_manage_permissions ?? false;
    const shouldViewAdminMenu = canManagePermissions;
    const adminMenu = <>
        <DropdownMenuGroup>
            {canManagePermissions && <DropdownMenuItem>
                <Link href={"/permissions"} className="w-full">
                    Permissions
                </Link>
            </DropdownMenuItem>}
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
    </>

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                    <UserAvatar username={user?.username ?? "unknown"} avatarUrl={user?.avatar_url} />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56" align="end" forceMount>
                <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                        <p className="text-sm font-medium leading-none">{user?.username ?? "unknown"}</p>
                        {/*<p className="text-xs leading-none text-muted-foreground">*/}
                        {/*    m@example.com*/}
                        {/*</p>*/}
                    </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                {/*<DropdownMenuGroup>*/}
                {/*    <DropdownMenuItem>*/}
                {/*        <Link href={"/profile"} className="w-full">*/}
                {/*            Profile*/}
                {/*        </Link>*/}
                {/*    </DropdownMenuItem>*/}
                {/*    <DropdownMenuItem>*/}
                {/*        <Link href={"/profile/uploads"} className="w-full">*/}
                {/*            Uploads*/}
                {/*        </Link>*/}
                {/*    </DropdownMenuItem>*/}
                {/*    <DropdownMenuItem>*/}
                {/*        <Link href={"/settings"} className="w-full">*/}
                {/*            Settings*/}
                {/*        </Link>*/}
                {/*    </DropdownMenuItem>*/}
                {/*</DropdownMenuGroup>*/}
                {/*<DropdownMenuSeparator />*/}
                {shouldViewAdminMenu && adminMenu}
                <DropdownMenuItem>
                    <Link href={"/logout"} className="w-full" prefetch={false}>
                        Log out
                    </Link>
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}