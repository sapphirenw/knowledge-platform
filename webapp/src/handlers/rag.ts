"import server-only"

import { RAGRequest, RAGResponse } from "@/types/rag"

export default async function HandleRAG(req: RAGRequest): Promise<Resp<RAGResponse>> {
    try {
        let USER_ID = "018fe32b-e136-787f-98b8-c2017a64ee31"
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
        return {
            error: "Unknown error"
        }
    }
}