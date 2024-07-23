"use client"

import { getUsage } from "@/actions/usage"
import { UsageGroupedRecord, UsageResponse } from "@/types/usage"
import { useState } from "react"

import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
    TableFooter,
} from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import DefaultLoader from "@/components/default_loader"

export default function CustomerUsageClient({
    initialData,
    usageGrouped,
}: {
    initialData: UsageResponse | undefined,
    usageGrouped: UsageGroupedRecord[] | undefined,
}) {
    const [data, setData] = useState<UsageResponse | undefined>(initialData)
    const [page, setPage] = useState(1)
    const [isLoading, setIsLoading] = useState(false)

    const canGetNext = () => {
        if (data === undefined) {
            return false
        }
        if (page >= data.metadata.pageCount) {
            return false
        }
        return true
    }

    const nextPage = () => {
        if (!canGetNext()) {
            return
        }

        refreshData(page + 1)
    }

    const canGetPrev = () => {
        if (data === undefined) {
            return false
        }
        if (page === 1) {
            return false
        }
        return true
    }

    const prevPage = () => {
        if (!canGetPrev()) {
            return
        }

        refreshData(page - 1)
    }

    const refreshData = async (newPage: number) => {
        setIsLoading(true)
        try {
            const newData = await getUsage({ page: page })
            setData(newData)
            setPage(newPage)
        } catch (e) {
            if (e instanceof Error) console.error(e)
            setData(undefined)
        }
        setIsLoading(false)
    }

    const getItems = () => {
        if (data === undefined) {
            return <div className="">Failed to get usage</div>
        }

        const items = []
        for (let i = 0; i < data.records.length; i++) {
            items.push(<TableRow key={`doc-${i}`}>
                <TableCell>{data.records[i].model}</TableCell>
                <TableCell>{new Date(data.records[i].createdAt).toLocaleString()}</TableCell>
                <TableCell>{data.records[i].inputTokens}</TableCell>
                <TableCell>{data.records[i].outputTokens}</TableCell>
                <TableCell className="text-right">{data.records[i].totalTokens}</TableCell>
            </TableRow>)
        }
        return items
    }

    const getUsageGrouped = () => {
        if (usageGrouped === undefined) {
            return <div className="">Failed to get Usage data</div>
        }

        return <Table containerClassname="h-fit overflow-y-auto relative">
            <TableHeader className="sticky w-full top-0">
                <TableRow>
                    <TableHead>Model</TableHead>
                    <TableHead>Total Input Tokens</TableHead>
                    <TableHead>Total Output Tokens</TableHead>
                    <TableHead>Input Cost</TableHead>
                    <TableHead>Output Cost</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {usageGrouped.map((item, i) => <TableRow key={`usageGrouped-${i}`}>
                    <TableCell>{item.model}</TableCell>
                    <TableCell>{item.inputTokensSum}</TableCell>
                    <TableCell>{item.outputTokensSum}</TableCell>
                    <TableCell>${item.inputCostCalculated}</TableCell>
                    <TableCell>${item.outputCostCalculated}</TableCell>
                </TableRow>)}
            </TableBody>
        </Table>
    }

    return <div className="w-full space-y-4">
        {getUsageGrouped()}
        <div className="w-full">
            <Table containerClassname="h-fit overflow-y-auto relative">
                <TableHeader className="sticky w-full top-0">
                    <TableRow>
                        <TableHead>Model</TableHead>
                        <TableHead>Date</TableHead>
                        <TableHead>Input Tokens</TableHead>
                        <TableHead>Output Tokens</TableHead>
                        <TableHead className="text-right">Total Tokens</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {getItems()}
                </TableBody>
            </Table>
            <div className="flex items-center justify-end space-x-2 py-4">
                <div className="flex-1 text-sm text-muted-foreground">
                </div>
                <div className="flex items-center space-x-2">
                    <div className="">
                        {isLoading ? <DefaultLoader /> : null}
                    </div>
                    <div className="">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => prevPage()}
                            disabled={!canGetPrev()}
                        >
                            Previous
                        </Button>
                    </div>
                    <div className="">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => nextPage()}
                            disabled={!canGetNext()}
                        >
                            Next
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    </div>
}