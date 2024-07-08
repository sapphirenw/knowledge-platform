export type LLM = {
    id: string
    customerId?: string
    title: string
    color?: string
    model: string
    temperature: number
    instructions: string
    isDefault: boolean
    public: boolean
    createdAt: string
    updatedAt: string
}

export type AvailableModel = {
    id: string
    provider: string
    displayName: string
    description: string
    inputTokenLimit: number
    outputTokenLimit: number
    currency: string
    inputCostPerMillionTokens: number
    outputCostPerMillionTokens: number
    depreciatedWarning: boolean
    isDepreciated: boolean
    createdAt: string
    updatedAt: string
}

export type ModelRow = {
    llm: LLM
    availableModel: AvailableModel
}