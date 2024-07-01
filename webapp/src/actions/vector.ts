"use server"

import { VectorizeJob } from "@/types/vector"
import { cookies } from "next/headers"
import { getCID } from "./customer"
import { sendRequestV1 } from "./api"

export async function createVectorizeRequest(): Promise<VectorizeJob> {
    const cid = await getCID()
    return await sendRequestV1<VectorizeJob>({
        route: `customers/${cid}/vectorstore/vectorize`,
        method: "POST",
        body: JSON.stringify({
            documents: true,
            websites: true,
        })
    })
}

export async function getAllVectorizeRequests(): Promise<VectorizeJob[]> {
    const cid = await getCID()
    return await sendRequestV1<VectorizeJob[]>({
        route: `customers/${cid}/vectorstore/vectorize`,
    })
}