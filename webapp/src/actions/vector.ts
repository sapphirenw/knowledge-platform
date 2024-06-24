"use server"

import { cookies } from "next/headers"

export async function vectorizeDatastore(): Promise<void> {
    const cid = cookies().get("cid")?.value
    if (cid == undefined) {
        throw new Error("no cid")
    }
    let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/vectorizeDocuments`, {
        method: "PUT",
        cache: 'no-store',
    })
    if (!response.ok) {
        const data = await response.text()
        console.log("there was an error:", data)
        throw new Error("There was an unknown error")
    }
}