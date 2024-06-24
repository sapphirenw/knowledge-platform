"use server"

import { ConversationResponse } from "@/types/conversation"
import { RAGRequest, RAGResponse } from "@/types/rag"
import { cookies } from "next/headers"

export async function HandleRAG(req: RAGRequest): Promise<Resp<RAGResponse>> {
    try {
        const cid = cookies().get("cid")?.value
        if (cid == undefined) {
            throw new Error("no cid")
        }

        let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/rag`, {
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