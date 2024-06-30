"use server"

import { HandleWebsiteRequest, HandleWebsiteResponse, Website } from "@/types/websites";
import { getCID } from "./customer";

export async function handleWebsite(payload: HandleWebsiteRequest): Promise<HandleWebsiteResponse> {
    try {
        const cid = await getCID()
        const response = await fetch(`${process.env.DB_HOST}/customers/${cid}/websites`, {
            method: "POST",
            cache: 'no-store',
            body: JSON.stringify(payload)
        })
        if (!response.ok) {
            console.log(await response.text())
            throw new Error("The request failed")
        }

        return await response.json() as HandleWebsiteResponse
    } catch (e) {
        if (e instanceof Error) console.log(e)
        throw e
    }
}

export async function getWebsites(): Promise<Website[]> {
    const cid = await getCID()
    const response = await fetch(`${process.env.DB_HOST}/customers/${cid}/websites`, {
        method: "GET",
        cache: 'no-store',
    })
    if (!response.ok) {
        throw new Error(await response.text())
    }

    const data = await response.json() as Website[]
    console.log(data)

    return data
}