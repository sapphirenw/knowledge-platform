"use server"

import { Document, DocumentChunkedResponse, DocumentCleanedResponse, FileValidationResponse, ListFolderResponse, PresignedUrlResponse } from "@/types/document"
import { humanFileSize } from "@/utils/humanFileSize"
import { generateSHA256 } from "@/utils/sha"
import { cookies } from "next/headers"
import { getCID } from "./customer"
import { sendRequestV1 } from "./api"

export async function listFolder() {
    const cid = await getCID()
    const response = await sendRequestV1<ListFolderResponse>({
        route: `customers/${cid}/folders`
    })

    return response
}

/**
 * Validates that all of the files passed to this function are valid to be uploaded to the docstore.
 * @param form form data to parse. Must contain a list of files in the `form.values()` method
 * @returns a list of validation responses
 */
export async function validateDocuments(form: FormData): Promise<FileValidationResponse[]> {
    const response: FileValidationResponse[] = []
    const maxFileSize = parseInt(process.env.MAX_FILE_SIZE ?? "1000000")

    const entries = Array.from(form.values())
    for (let i = 0; i < entries.length; i++) {
        const file = entries[i] as File
        const rec: FileValidationResponse = { filename: file.name }
        if (file.size > maxFileSize) {
            rec.error = `The file is too big. Max size = ${humanFileSize(maxFileSize)}. File size = ${humanFileSize(file.size)}`
        }
        response.push(rec)
    }

    return response
}

/**
 * Upload documents to the fileserver
 * @param form form data from an html form. Must contain a list of documents
 * @returns boolean for whether the upload was successful or not
 */
export async function uploadDocuments(form: FormData): Promise<boolean> {
    try {
        const cid = await getCID()

        // parse through all files
        const entries = Array.from(form.values())
        for (let i = 0; i < entries.length; i++) {
            // handle the file
            const file = entries[i] as File
            const buffer = Buffer.from(await file.arrayBuffer());
            const sig = await generateSHA256(buffer)

            // generate the presigned url
            const payload = {
                filename: file.name,
                mime: file.type,
                signature: sig,
                size: file.size
            }
            console.log(payload)
            const presignedData = await sendRequestV1<PresignedUrlResponse>({
                route: `customers/${cid}/generatePresignedUrl`,
                method: "POST",
                body: JSON.stringify(payload),
            })
            console.log("created pre-signed url")

            // parse the response
            const url = Buffer.from(presignedData.uploadUrl, 'base64').toString('utf-8');

            // upload the file
            const uploadResp = await fetch(url, {
                method: presignedData.method,
                headers: {
                    'Content-Type': file.type
                },
                body: buffer
            });
            if (!uploadResp.ok) {
                console.log(await uploadResp.text())
                throw new Error("failed to upload the file")
            }

            console.log("successfully uploaded file")

            // notify of a successful upload
            await sendRequestV1<undefined>({
                route: `customers/${cid}/documents/${presignedData.documentId}/validate`,
                method: "PUT",
            })
            console.log("Successfully notified of successful upload")
        }

        return true
    } catch (e) {
        if (e instanceof Error) console.log(e)
        return false
    }
}

export async function getDocument(documentId: string) {
    const cid = await getCID()
    return await sendRequestV1<Document>({
        route: `customers/${cid}/documents/${documentId}`
    })
}

export async function getDocumentCleaned(documentId: string) {
    const cid = await getCID()
    return await sendRequestV1<DocumentCleanedResponse>({
        route: `customers/${cid}/documents/${documentId}/cleaned`
    })
}

export async function getDocumentChunked(documentId: string) {
    const cid = await getCID()
    return await sendRequestV1<DocumentChunkedResponse>({
        route: `customers/${cid}/documents/${documentId}/chunked`
    })
}