"use server"

import { RAG2Request, RAGRequest, RAGResponse } from "@/types/rag"
import { cookies } from "next/headers"
import { getCID } from "./customer"
import { sendRequestV1 } from "./api"

// sends an invocation to create a new rag websocket
export async function createRagRequest(): Promise<RAG2Request> {
    const cid = await getCID()
    return await sendRequestV1<RAG2Request>({
        route: `customers/${cid}/rag2Init`,
        method: "GET"
    })
}

export async function handleRAG(req: RAGRequest): Promise<RAGResponse> {
    const cid = await getCID()
    // get the conversationId
    const convId = cookies().get("conversationId")?.value ?? ""
    return await sendRequestV1<RAGResponse>({
        route: `customers/${cid}/rag?conversationId=${convId}`,
        method: "POST",
        body: JSON.stringify(req),
    })
}