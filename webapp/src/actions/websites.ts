"use server"

import { HandleWebsiteRequest, HandleWebsiteResponse, Website, WebsitePage, WebsitePageContentResponse } from "@/types/websites";
import { getCID } from "./customer";
import { sendRequestV1 } from "./api";
import { redirect } from "next/navigation";

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

export async function getWebsite(siteId: string) {
    const cid = await getCID()
    return await sendRequestV1<Website>({
        route: `customers/${cid}/websites/${siteId}`,
    })
}

export async function getWebsitePages(siteId: string) {
    const cid = await getCID()
    return await sendRequestV1<WebsitePage[]>({
        route: `customers/${cid}/websites/${siteId}/pages`,
    })
}

export async function getWebsitePageContent(siteId: string, pageId: string) {
    const cid = await getCID()
    return await sendRequestV1<WebsitePageContentResponse>({
        route: `customers/${cid}/websites/${siteId}/pages/${pageId}/content?getCleaned=true&getChunked=true`,
    })
}

export async function deleteWebsite(siteId: string) {
    const cid = await getCID()
    await sendRequestV1<undefined>({
        route: `customers/${cid}/websites/${siteId}`,
        method: "DELETE",
    })

    // if here then the delete call worked, so redirect
    redirect("/settings/datastore")
}