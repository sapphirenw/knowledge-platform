"use server"

import { Conversation, ConversationResponse } from "@/types/conversation"
import { cookies } from "next/headers"

export async function GetAllConversations(): Promise<Resp<Conversation[]>> {
    try {
        const cid = cookies().get("cid")?.value
        if (cid == undefined) {
            throw new Error("no cid")
        }
        let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/conversations`, {
            method: "GET",
            cache: 'no-store',
        })
        if (response.status != 200) {
            return {
                error: await response.text()
            }
        }
        return {
            data: await response.json() as Conversation[]
        }
    } catch (e) {
        console.log(e)
        return {
            error: "Unknown error"
        }
    }
}

export async function GetConversation(convId: string): Promise<Resp<ConversationResponse>> {
    try {
        const cid = cookies().get("cid")?.value
        if (cid == undefined) {
            throw new Error("no cid")
        }
        let response = await fetch(`${process.env.DB_HOST}/customers/${cid}/conversations/${convId}`, {
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