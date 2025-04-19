"use server";

import {LoginButton} from "@/components/login_button";

export default async function Home() {

  return (
      <>
        <div className="h-dvh flex flex-col items-center justify-center bg-zinc-900">
            <div className="text-lg text-white font-medium mb-4">
              klaphat
            </div>
            <div>
                <LoginButton />
            </div>
        </div>
      </>
  )
}
