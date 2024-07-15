"use server"

import { HandleWebsiteRequest, HandleWebsiteResponse, Website } from "@/types/websites";
import { getCID } from "./customer";
import { sendRequestV1 } from "./api";

export async function searchWebsite(payload: HandleWebsiteRequest): Promise<HandleWebsiteResponse> {
    const cid = await getCID()
    return await sendRequestV1<HandleWebsiteResponse>({
        route: `customers/${cid}/websites`,
        method: "PUT",
        body: JSON.stringify(payload)
    })
}

export async function insertWebsite(payload: HandleWebsiteRequest): Promise<HandleWebsiteResponse> {
    const cid = await getCID()
    return await sendRequestV1<HandleWebsiteResponse>({
        route: `customers/${cid}/websites`,
        method: "POST",
        body: JSON.stringify(payload)
    })
}

export async function insertSingleWebsitePage(domain: string): Promise<boolean> {
    const cid = await getCID()
    await sendRequestV1({
        route: `customers/${cid}/insertSingleWebsitePage`,
        method: "POST",
        body: JSON.stringify({ "domain": domain })
    })
    return true
}

export async function getWebsites(): Promise<Website[]> {
    const cid = await getCID()
    return await sendRequestV1<Website[]>({
        route: `customers/${cid}/websites`,
    })
}