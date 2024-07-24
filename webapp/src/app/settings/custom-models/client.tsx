"use client"

import { getCustomerLLMs } from "@/actions/llm"
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
import { Settings } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function CustomerModelsClient() {
    return <div className="w-full">
        <CustomerLLMsView />
    </div>
}

function CustomerLLMsView() {

    const { status, data, error } = useQuery({
        queryKey: ['customerLLMs', false],
        queryFn: () => getCustomerLLMs(false)
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.error(error)
        return <p>Failed to get your custom models</p>
    }

    if (data.length === 0) {
        return <div className="">No custom models found.</div>
    }

    const getTableRows = () => {
        const items = []

        for (let i = 0; i < data.length; i++) {
            items.push(<TableRow key={`model-${i}`}>
                <TableCell className="w-[50px]">
                    <Button size="icon" variant="outline"><Settings size={16} /></Button>
                </TableCell>
                <TableCell className="text-sm font-semibold opacity-50">{data[i].availableModel.id}</TableCell>
                <TableCell>{data[i].llm.title}</TableCell>
                <TableCell>{data[i].llm.temperature.toPrecision(2)}</TableCell>
                <TableCell>{data[i].llm.instructions}</TableCell>
            </TableRow>)
        }
        return items
    }

    return <div className="w-full">
        <div className="overflow-hidden">
            <Table containerClassname="">
                <TableHeader className="sticky w-full top-0">
                    <TableRow>
                        <TableHead className="w-[50px]"></TableHead>
                        <TableHead>LLM</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Temperature</TableHead>
                        <TableHead>Instructions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {getTableRows()}
                </TableBody>
            </Table>
        </div>
    </div>
}