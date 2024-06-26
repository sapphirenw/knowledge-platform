"use server"

import { Conversation, ConversationResponse } from "@/types/conversation"
import { getCID } from "./customer"
import { cookies } from "next/headers"

export async function getAllConversations(): Promise<Conversation[]> {
    const cid = await getCID()
    let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/conversations`, {
        method: "GET",
        cache: 'no-store',
    })
    if (!response.ok) {
        console.log(await response.text())
        throw new Error("failed to fetch the data")
    }
    return await response.json() as Conversation[]
}

export async function getConversation(): Promise<ConversationResponse> {
    // read the cookie
    const convId = cookies().get("conversationId")?.value
    if (convId === "" || convId == undefined) {
        return {
            "conversationId": "",
            "messages": [],
        }
    }

    // fetch with the conversationId
    const cid = await getCID()
    let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/conversations/${convId}`, {
        method: "GET",
        cache: 'no-store',
    })
    if (!response.ok) {
        console.log(await response.text())
        throw new Error("failed to fetch the data")
    }

    return await response.json() as ConversationResponse
}