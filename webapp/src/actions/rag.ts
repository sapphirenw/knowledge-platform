"use server"

import { RAGRequest, RAGResponse } from "@/types/rag"
import { cookies } from "next/headers"
import { getCID } from "./customer"
import { sendRequestV1 } from "./api"

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