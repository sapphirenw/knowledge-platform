"use server"

import { Conversation, ConversationResponse } from "@/types/conversation"

export async function GetAllConversations(): Promise<Resp<Conversation[]>> {
    try {
        let response = await fetch(`${process.env.DB_HOST}/customers/${process.env.TMP_USER_ID}/conversations`, {
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
        let response = await fetch(`${process.env.DB_HOST}/customers/${process.env.TMP_USER_ID}/conversations/${convId}`, {
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