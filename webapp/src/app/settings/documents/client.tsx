"use client"

import { listFolder } from "@/actions/document"
import { useQuery } from "@tanstack/react-query"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

import { humanFileSize } from "@/utils/humanFileSize"
import DefaultLoader from "@/components/default_loader"
import Link from "next/link"


export default function DocumentsViewClient() {
    // get the users files with react query
    const { status, data, error } = useQuery({
        queryKey: ['documents', null],
        queryFn: () => listFolder(),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.log(error)
        return <div>There was an error</div>
    }

    const getTableRows = () => {
        const items = []
        for (let i = 0; i < data.documents.length; i++) {
            items.push(<TableRow key={`doc-${i}`}>
                <TableCell className="text-sm font-semibold opacity-50">
                    {data.documents[i].vectorSha256.trim() == "" ? "False" : "True"}
                </TableCell>
                <TableCell className="font-medium">
                    <Link className="text-primary hover:opacity-50 underline" href={`/settings/documents/${data.documents[i].id}`}>{data.documents[i].filename}</Link>
                </TableCell>
                <TableCell>{humanFileSize(data.documents[i].sizeBytes)}</TableCell>
                <TableCell className="text-right">{new Date(data.documents[i].createdAt).toLocaleString()}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="w-full">
        <div className="overflow-hidden">
            <Table containerClassname="h-fit max-h-[500px] overflow-y-auto relative">
                <TableHeader className="sticky w-full top-0">
                    <TableRow>
                        <TableHead className="">In Vector-Store</TableHead>
                        <TableHead className="">Filename</TableHead>
                        <TableHead>Size</TableHead>
                        <TableHead className="text-right">Created</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {getTableRows()}
                </TableBody>
            </Table>
        </div>
    </div>
}