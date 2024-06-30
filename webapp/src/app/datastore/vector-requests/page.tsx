"use client"

import { getAllVectorizeRequests } from "@/actions/vector"
import { useQuery } from "@tanstack/react-query"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import DefaultLoader from "@/components/default_loader"
import VectorizationRequest from "@/components/vectorization_request"

export default function VectorRequests() {
    // get the users files with react query
    const { status, data, error } = useQuery({
        queryKey: ['vectorRequests'],
        queryFn: () => getAllVectorizeRequests(),
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
        for (let i = 0; i < data.length; i++) {
            items.push(<TableRow key={`doc-${i}`}>
                <TableCell className="font-medium">{data[i].id}</TableCell>
                <TableCell>{data[i].status}</TableCell>
                <TableCell>{data[i].message}</TableCell>
                <TableCell className="text-right">{new Date(data[i].createdAt).toLocaleString()}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="w-full space-y-4">
        <VectorizationRequest />
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