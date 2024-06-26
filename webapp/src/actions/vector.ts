"use server"

import { VectorizeJob } from "@/types/vector"
import { cookies } from "next/headers"
import { getCID } from "./customer"

export async function createVectorizeRequest(): Promise<VectorizeJob> {
    const cid = await getCID()
    let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/vectorstore/vectorize`, {
        method: "POST",
        cache: 'no-store',
        body: JSON.stringify({
            documents: true,
            websites: true,
        })
    })
    if (!response.ok) {
        const data = await response.text()
        console.log("there was an error:", data)
        throw new Error("There was an unknown error")
    }

    const data = await response.json() as VectorizeJob
    console.log(data)
    return data
}

export async function getAllVectorizeRequests(): Promise<VectorizeJob[]> {
    const cid = await getCID()
    let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/vectorstore/vectorize`, {
        method: "GET",
        next: { revalidate: 3 },
    })
    if (!response.ok) {
        const data = await response.text()
        console.log("there was an error:", data)
        throw new Error("There was an unknown error")
    }
    const data = await response.json() as VectorizeJob[]
    console.log(data)
    return data
}