"use server"

import { AvailableModel, CreateLLMRequest, ModelRow } from "@/types/llm"
import { sendRequestV1 } from "./api"
import { getCID } from "./customer"

export async function getCustomerLLMs(includeAll: boolean) {
    const cid = await getCID()
    let response = await sendRequestV1<ModelRow[]>({
        route: `customers/${cid}/llms?includeAll=${includeAll}`,
        // debugPrint: true,
    })
    return response
}

export async function getAvailableModels(provider: string) {
    let response = await sendRequestV1<AvailableModel[]>({
        route: `llms/availableModels?provider=${provider}`,
    })
    return response
}

export async function createLLM(request: CreateLLMRequest) {
    const cid = await getCID()
    await sendRequestV1({
        route: `customers/${cid}/llms`,
        method: "POST",
        body: JSON.stringify(request),
    })
}

export async function updateLLM(id: string, request: CreateLLMRequest) {
    const cid = await getCID()
    await sendRequestV1({
        route: `customers/${cid}/llms/${id}`,
        method: "PUT",
        body: JSON.stringify(request),
    })
}