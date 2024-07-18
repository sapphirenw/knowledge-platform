export type FileValidationResponse = {
    filename: string
    error?: string
}

export type PresignedUrlResponse = {
    uploadUrl: string
    method: string
    documentId: string
}

export type Document = {
    createdAt: string;
    customerId: string;
    datastoreId: string;
    datastoreType: string;
    filename: string;
    id: string;
    parentId: string | null;
    sha256: string;
    sizeBytes: number;
    summary: string;
    summarySha256: string;
    vectorSha256: string;
    type: string;
    updatedAt: string;
    validated: boolean;
};

export type ListFolderResponse = {
    // self: null,
    folders: [],
    documents: Document[]
}

export type DocumentCleanedResponse = {
    document: Document
    cleaned: string
}

export type DocumentChunkedResponse = {
    document: Document
    chunks: string[]
}