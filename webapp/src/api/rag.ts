"use server"

import { ConversationResponse } from "@/types/conversation"
import { RAGRequest, RAGResponse } from "@/types/rag"

export async function HandleRAG(req: RAGRequest): Promise<Resp<RAGResponse>> {
    try {
        let USER_ID = "01900764-490c-7a3d-a513-4111f579a7b1"
        let response = await fetch(`http://localhost:8000/customers/${USER_ID}/rag`, {
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

export async function FetchConversation(convId: string): Promise<Resp<ConversationResponse>> {
    try {
        let response = await fetch(`http://localhost:8000/conversations/${convId}`, {
            method: "GET",
            cache: 'no-store',
        })
        if (response.status != 200) {
            return {
                error: await response.text()
            }
        }
        return {
            data: await response.json() as ConversationResponse
        }
    } catch (e) {
        console.log(e)
        return {
            error: "Unknown error"
        }
    }
}