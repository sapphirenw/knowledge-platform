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
import DefaultLoader from "@/components/default_loader"
import { Button } from "@/components/ui/button"
import Link from "next/link"


export default function UserFiles() {
    // get the users files with react query
    const { status, data, error } = useQuery({
        queryKey: ['files'],
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
                <TableCell className="font-medium">
                    <Link className="text-primary hover:opacity-50 underline" href={`/datastore/documents/${data.documents[i].id}`}>{data.documents[i].filename}</Link>
                </TableCell>
                <TableCell>{humanFileSize(data.documents[i].sizeBytes)}</TableCell>
                <TableCell className="text-right">{new Date(data.documents[i].createdAt).toLocaleString()}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="w-full">
        {/* <div className="p-4">
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
        </div> */}
        <div className="overflow-hidden">
            <Table containerClassname="h-fit max-h-[500px] overflow-y-auto relative">
                <TableHeader className="sticky w-full top-0">
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