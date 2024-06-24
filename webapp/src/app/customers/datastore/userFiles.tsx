"use client"

import { listFolder } from "@/actions/document"
import { useQuery } from "@tanstack/react-query"
import { Loader2, Slash } from "lucide-react"

import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { humanFileSize } from "@/utils/humanFileSize"


export default function UserFiles() {
    // get the users files with react query
    const { status, data, error } = useQuery({
        queryKey: ['files'],
        queryFn: ({ signal }) => listFolder(),
        staleTime: 60 * 1000,
    })

    if (status === "pending") {
        return <Loader2 className="mr-2 h-4 w-4 animate-spin" />
    }

    if (status === "error") {
        console.log(error)
        return <div>There was an error</div>
    }

    const getTableRows = () => {
        const items = []
        for (let i = 0; i < data.documents.length; i++) {
            items.push(<TableRow id={`doc-${i}`}>
                <TableCell className="font-medium">{data.documents[i].filename}</TableCell>
                <TableCell>{humanFileSize(data.documents[i].sizeBytes)}</TableCell>
                <TableCell className="text-right">{new Date(data.documents[i].createdAt).toLocaleString()}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="w-full">
        <div className="p-4">
            <Breadcrumb>
                <BreadcrumbList>
                    <BreadcrumbSeparator>
                        <Slash />
                    </BreadcrumbSeparator>
                    <BreadcrumbItem>
                        <BreadcrumbLink href="/customers/folders">Root</BreadcrumbLink>
                    </BreadcrumbItem>
                </BreadcrumbList>
            </Breadcrumb>
        </div>
        <div className="border border-border rounded-md w-full overflow-hidden">
            <Table>
                <TableHeader>
                    <TableRow>
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