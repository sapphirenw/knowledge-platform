export type UsageRecord = {
    id: string;
    customerId: string;
    conversationId?: string | null;
    model: string;
    inputTokens: number;
    outputTokens: number;
    totalTokens: number;
    createdAt: string;
};

export type UsageMetadata = {
    pageCount: number;
}

export type UsageResponse = {
    metadata: UsageMetadata;
    records: UsageRecord[];
}

export type UsageGroupedRecord = {
    model: string
    inputTokensSum: number
    outputTokensSum: number
    totalTokensSum: number
    inputCostPerMillionTokens: number
    outputCostPerMillionTokens: number
    inputCostCalculated: number
    outputCostCalculated: number
}