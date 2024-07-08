"use server"

import { ModelRow } from "@/types/llm"
import { sendRequestV1 } from "./api"
import { getCID } from "./customer"

export async function getAvailableLLMs() {
    const cid = await getCID()
    let response = await sendRequestV1<ModelRow[]>({
        route: `customers/${cid}/llms`,
        // debugPrint: true,
    })
    return response
}