"use server";

import Link from "next/link";
import {ENDPOINT_PAGE_HOME} from "@/lib/endpoints";
import {getInfo} from "@/lib/api";

export async function NavbarLogo() {
    const info = await getInfo();
    return (
        <>
            <Link href={ENDPOINT_PAGE_HOME} className="content-center">
                <h1 className="text-lg font-medium no-underline cursor-pointer select-none">{info.title}</h1>
            </Link>
        </>
    )
}
