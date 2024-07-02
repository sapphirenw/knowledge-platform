export type VectorizeJob = {
    id: string;
    customerId: string;
    status: { vectorizeJobStatus: string, valid: boolean }
    message: string;
    error: string;
    documents: boolean;
    websites: boolean;
    createdAt: string;
    updatedAt: string;
};