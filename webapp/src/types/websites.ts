export type HandleWebsiteRequest = {
    domain: string
    whitelist?: string[]
    blacklist?: string[]
    useSitemap?: boolean
    allowOtherDomains?: boolean
    pages?: string[]
}

export type Website = {
    id: string
    customerId: string
    protocol: string
    domain: string
    path: string
    pageCount: number
    blacklist?: string[]
    whitelist?: string[]
    createdAt?: string
    updatedAt?: string
}

export type WebsitePage = {
    id: string
    customerId: string
    websiteId: string
    url: string
    sha256: string
    isValid: boolean
    metadata?: any
    summary: string
    summarySha256: string
    createdAt?: string
    updatedAt?: string
}

export type HandleWebsiteResponse = {
    site: Website
    pages: WebsitePage[]
}

export type WebsitePageContentResponse = {
    page: WebsitePage;
    cleaned?: string;
    chunks?: string[];
}