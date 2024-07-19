"use server"

import { BetaApiKey } from "@/types/beta_auth"
import { sendRequestV1 } from "./api"

export async function createApiKey(name: string, isAdmin: boolean, authToken: string): Promise<BetaApiKey> {
    return await sendRequestV1<BetaApiKey>({
        route: `beta/createBetaApiKey?name=${name}&isAdmin=${isAdmin}`,
        method: "POST",
        headers: new Headers({
            "x-master-auth-token": authToken,
        }),
        ignoreApiKey: true,
    })
}