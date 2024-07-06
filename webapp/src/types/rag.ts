import { ConversationMessage } from "./conversation";

export type RAG2Request = {
    id: string
    path: string
}

export interface RAGRequest {
    input: string;
    checkQuality?: boolean;
    summaryModelId?: string;
    chatModelId?: string;
    folderIds?: string[];
    documentIds?: string[];
    websiteIds?: string[];
    websitePageIds?: string[];
}

export interface RAGResponse {
    conversationId: string
    documents?: any[]
    websitePages?: any[]
    message: ConversationMessage
}

// MessageType ragMessageType `json:"messageType"`

// // dependent on the message type

// ChatMessage    *gollm.Message `json:"chatMessage,omitempty"`
// ConversationId string         `json:"conversationId"`
// NewTitle       string         `json:"newTitle,omitempty"`
// Error          string         `json:"error,omitempty"`

export type RagMessagePayload = {
    messageType: string

    chatMessage?: ConversationMessage
    conversationId?: string
    newTitle?: string
    error?: string
}