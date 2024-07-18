import { getDocument } from "@/actions/document";
import ErrorPage from "@/components/error_page";
import DocumentViewClient from "./client";

export default async function DocumentView({
    params,
}: {
    params: { documentId: string };
}) {
    try {
        const doc = await getDocument(params.documentId)
        return <DocumentViewClient document={doc} />
    } catch {
        return <ErrorPage msg="There was an issue getting the document" />
    }
}