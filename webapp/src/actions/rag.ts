"use server"

import { ConversationResponse } from "@/types/conversation"
import { RAGRequest, RAGResponse } from "@/types/rag"

export async function HandleRAG(req: RAGRequest): Promise<Resp<RAGResponse>> {
    try {
        let response = await fetch(`${process.env.DB_HOST}/customers/${process.env.TMP_USER_ID}/rag`, {
            method: "POST",
            cache: 'no-store',
            body: JSON.stringify(req),
        })
        if (response.status != 200) {
            return {
                error: await response.text()
            }
        }
        return {
            data: await response.json() as RAGResponse
        }
    } catch (e) {
        console.log(e)
        return {
            error: "Unknown error"
        }
    }
}