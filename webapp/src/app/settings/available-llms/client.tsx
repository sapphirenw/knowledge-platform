"use client"

import { getAvailableModels } from "@/actions/llm"
import DefaultLoader from "@/components/default_loader"
import { useQuery } from "@tanstack/react-query"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

export default function AvailableModelsViewClient() {
    const { status, data, error } = useQuery({
        queryKey: ['availableModels'],
        queryFn: () => getAvailableModels(""),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.error(error)
        return <div>There was an error</div>
    }

    const getTableRows = () => {
        const items = []

        const d = data.sort((a, b) => a.provider.localeCompare(b.provider))
        for (let i = 0; i < data.length; i++) {
            items.push(<TableRow key={`model-${i}`}>
                <TableCell className="text-sm font-semibold opacity-50">{data[i].provider}</TableCell>
                <TableCell>{data[i].displayName}</TableCell>
                <TableCell>{data[i].description}</TableCell>
                <TableCell>${data[i].inputCostPerMillionTokens}</TableCell>
                <TableCell>${data[i].outputCostPerMillionTokens}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="w-full">
        <div className="overflow-hidden">
            <Table containerClassname="">
                <TableHeader className="sticky w-full top-0">
                    <TableRow>
                        <TableHead>Provider</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Description</TableHead>
                        <TableHead>{"Input Cost (Million Tokens)"}</TableHead>
                        <TableHead>{"Output Cost (Million Tokens)"}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {getTableRows()}
                </TableBody>
            </Table>
        </div>
    </div>
}