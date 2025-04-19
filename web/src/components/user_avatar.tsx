import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {getInitials} from "@/lib/utils";

interface UserAvatarProps {
    username: string;
    avatarUrl?: string;
}

export function UserAvatar({ username, avatarUrl }: UserAvatarProps) {
    return <Avatar className="h-9 w-9">
        <AvatarImage src={`${avatarUrl}`} alt={`@${username}`} />
        <AvatarFallback>{getInitials(username)}</AvatarFallback>
    </Avatar>
}
