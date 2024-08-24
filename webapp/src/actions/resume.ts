"use server"

import { CreateResumeApplicationRequest, ResumeAbout, ResumeApplication, ResumeChecklistItem, ResumeItem } from "@/types/resume"
import { getCID } from "./customer"
import { sendRequestV1 } from "./api"
import { Document } from "@/types/document"

// export async function getResumeItems(): Promise<ResumeItem[]> {
//     const cid = await getCID()
//     return await sendRequestV1<ResumeItem[]>({
//         route: `customers/${cid}/resumes`,
//         method: "GET"
//     })
// }

// export async function createResumeItem(title: string): Promise<ResumeItem> {
// const cid = await getCID()
// return await sendRequestV1<ResumeItem>({
//     route: `customers/${cid}/resumes`,
//     method: "POST",
//     body: JSON.stringify({ "title": title }),
// })
// }

export async function getResumeItem(): Promise<ResumeItem> {
    const cid = await getCID()
    return await sendRequestV1<ResumeItem>({
        route: `customers/${cid}/resumes/${cid}`,
        method: "GET",
    })
}

export async function getResumeAbout(): Promise<ResumeAbout> {
    const cid = await getCID()
    return await sendRequestV1<ResumeAbout>({
        route: `customers/${cid}/resumes/${cid}/about`,
        method: "GET",
    })
}

export async function getResumeApplications(): Promise<ResumeApplication[]> {
    const cid = await getCID()
    return await sendRequestV1<ResumeApplication[]>({
        route: `customers/${cid}/resumes/${cid}/applications`,
        method: "GET",
    })
}

export async function createResumeApplication(req: CreateResumeApplicationRequest): Promise<ResumeApplication> {
    const cid = await getCID()
    return await sendRequestV1<ResumeApplication>({
        route: `customers/${cid}/resumes/${cid}/applications`,
        method: "POST",
        body: JSON.stringify(req),
    })
}

export async function getResumeChecklist() {
    const cid = await getCID()
    return await sendRequestV1<ResumeChecklistItem[]>({
        route: `customers/${cid}/resumes/${cid}/checklist`,
        method: "GET",
    })
}

export async function attachDocumentsToResume(docs: Document[]) {
    const cid = await getCID()
    return await sendRequestV1<undefined>({
        route: `customers/${cid}/resumes/${cid}/documents`,
        method: "POST",
        body: JSON.stringify({
            "documentIds": docs.map((item) => item.id),
        }),
    })
}

export async function getResumeDocuments() {
    const cid = await getCID()
    return await sendRequestV1<Document[]>({
        route: `customers/${cid}/resumes/${cid}/documents`,
        method: "GET",
    })
}
