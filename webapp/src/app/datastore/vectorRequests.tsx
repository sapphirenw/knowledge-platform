"use client"

import { getAllVectorizeRequests } from "@/actions/vector"
import { useQuery } from "@tanstack/react-query"
import { Loader2 } from "lucide-react"
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

export default function VectorRequests() {
    // get the users files with react query
    const { status, data, error } = useQuery({
        queryKey: ['vectorRequests'],
        queryFn: ({ signal }) => getAllVectorizeRequests(),
        staleTime: 0,
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
        for (let i = 0; i < data.length; i++) {
            items.push(<TableRow id={`doc-${i}`}>
                <TableCell className="font-medium">{data[i].id}</TableCell>
                <TableCell>{data[i].status}</TableCell>
                <TableCell>{data[i].message}</TableCell>
                <TableCell className="text-right">{new Date(data[i].createdAt).toLocaleString()}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="border border-border rounded-md w-full overflow-hidden">
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead className="">Job ID</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Message</TableHead>
                    <TableHead className="text-right">Created</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {getTableRows()}
            </TableBody>
        </Table>
    </div>
}