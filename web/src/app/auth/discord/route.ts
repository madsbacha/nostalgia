import {getDiscordUrl} from "@/lib/discord";
import {redirect} from "next/navigation";

export async function GET() {
    const discordLink = await getDiscordUrl();
    redirect(discordLink);
}