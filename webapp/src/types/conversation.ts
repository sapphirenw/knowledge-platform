export interface ConversationMessage {
    role: number
    message: string
    index: number
    id?: string
    name?: string
    toolArguments?: any
}

export interface ConversationResponse {
    conversationId: string
    title: string
    messages: ConversationMessage[]
}