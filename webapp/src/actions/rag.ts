"use server"

import { RAGRequest, RAGResponse } from "@/types/rag"
import { cookies } from "next/headers"
import { getCID } from "./customer"

export async function handleRAG(req: RAGRequest): Promise<Resp<RAGResponse>> {
    try {
        const cid = await getCID()

        // get the conversationId
        const convId = cookies().get("conversationId")?.value ?? ""

        let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/rag?conversationId=${convId}`, {
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