export interface RAGRequest {
    input: string;
    conversationId: string;
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
    documents: any[]
    websitePages: any[]
    response: string
}