"use client"

import { Document } from "@/types/document"

import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { humanFileSize } from "@/utils/humanFileSize"
import { useQuery } from "@tanstack/react-query"
import { getDocumentChunked, getDocumentCleaned } from "@/actions/document"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"
import { Separator } from "@/components/ui/separator"

export default function DocumentViewClient({
    document
}: {
    document: Document,
}) {

    const cleanedResp = useQuery({
        queryKey: ['documentCleaned', document.id],
        queryFn: () => getDocumentCleaned(document.id),
    })

    const chunkedResp = useQuery({
        queryKey: ['documentChunked', document.id],
        queryFn: () => getDocumentChunked(document.id),
    })

    const getCleanedView = () => {
        if (cleanedResp.status === "pending") {
            return <DefaultLoader />
        }
        if (cleanedResp.status === "error") {
            console.error(cleanedResp.error)
            return <ErrorPage msg="failed to get the cleaned content" />
        }

        return <p>{cleanedResp.data.cleaned}</p>
    }

    const getChunkedView = () => {
        if (chunkedResp.status === "pending") {
            return <DefaultLoader />
        }
        if (chunkedResp.status === "error") {
            console.error(chunkedResp.error)
            return <ErrorPage msg="failed to get the cleaned content" />
        }

        return <div className="space-y-2">
            {chunkedResp.data.chunks.map((item, i) => <p key={`chunk-${i}`}>{item}</p>)}
        </div>
    }

    return <div className="w-full space-y-8">
        <Table containerClassname="max-w-md mx-auto">
            <TableBody>
                <TableRow>
                    <TableCell className="font-medium">Filename</TableCell>
                    <TableCell>{document.filename}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">Parsed Filetype</TableCell>
                    <TableCell>{document.type}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">Size</TableCell>
                    <TableCell>{humanFileSize(document.sizeBytes)}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">Created</TableCell>
                    <TableCell>{new Date(document.createdAt).toLocaleString()}</TableCell>
                </TableRow>
                <TableRow>
                    <TableCell className="font-medium">Updated</TableCell>
                    <TableCell>{new Date(document.updatedAt).toLocaleString()}</TableCell>
                </TableRow>
            </TableBody>
        </Table>
        <div className="">
            <h4 className="font-medium">Cleaned Content:</h4>
            <Separator />
            <p className="prose break-all">{getCleanedView()}</p>
        </div>
        <div className="">
            <h4 className="font-medium">Chunks:</h4>
            <Separator />
            <p className="prose break-all">{getChunkedView()}</p>
        </div>
    </div>
}