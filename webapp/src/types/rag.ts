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