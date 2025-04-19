
import {
    Table,
    TableBody,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {WhitelistUser} from "@/lib/api/client";
import {PermissionsTableRow} from "@/components/permissions_table_row";

interface WhitelistTableProps {
    users: WhitelistUser[];
}

export function PermissionsTable({ users }: WhitelistTableProps) {
    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>User</TableHead>
                    <TableHead className="w-[200px] text-right">Can manage permissions</TableHead>
                    <TableHead className="w-[100px] text-right">Can upload</TableHead>
                    <TableHead className="w-[100px] text-right">Can view</TableHead>
                    <TableHead className="w-[100px] text-right">Whitelisted</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {users.map((user) => (
                    <PermissionsTableRow key={user.id} user={user} />
                ))}
            </TableBody>
        </Table>
    )
}
