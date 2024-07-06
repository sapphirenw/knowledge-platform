"use server"

import { Conversation, ConversationResponse } from "@/types/conversation"
import { getCID } from "./customer"
import { cookies } from "next/headers"
import { sendRequestV1 } from "./api"

export async function getAllConversations(): Promise<Conversation[]> {
    const cid = await getCID()
    let response = await sendRequestV1<Conversation[]>({
        route: `customers/${cid}/conversations`
    })
    return response
}

export async function getConversation(): Promise<ConversationResponse> {
    // read the cookie
    const convId = cookies().get("conversationId")?.value
    if (convId === "" || convId === undefined) {
        console.log("RETURNING EMPTY CONVERSATION")
        return {
            "conversationId": "",
            "messages": [],
        }
    }

    console.log("FETCHING CONVERSATION")

    // fetch with the conversationId
    try {
        const cid = await getCID()
        let response = await sendRequestV1<ConversationResponse>({
            route: `customers/${cid}/conversations/${convId}`
        })
        return response
    } catch (e) {
        if (e instanceof Error) console.error(e)
        // remove the convid
        cookies().delete("conversationId")
        throw e
    }
}