import {buttonVariants} from "@/components/ui/button";
import Link from "next/link";

export function LoginButton() {

    return <Link
        href={"/auth/discord"}
        className={buttonVariants({variant:"secondary"})}
        prefetch={false}>
        Login
    </Link>
}