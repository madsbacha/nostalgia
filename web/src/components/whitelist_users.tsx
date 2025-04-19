"use client";

import {useUsersForWhitelist} from "@/lib/api/client";
import {PermissionsTable} from "@/components/permissions_table";

export function WhitelistUsers() {
    const { users, isLoading, isError } = useUsersForWhitelist();

    if (isLoading) return <></>
    if (isError || !users) return <div>Failed to load</div>

    return (
        <>
            <h1 className="text-2xl py-6">Permissions</h1>
            <div className="pb-6">
                <PermissionsTable users={users} />
            </div>
        </>
    )
}