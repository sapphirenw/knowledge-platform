"use client"

import { listFolder } from "@/actions/document"
import { useQuery } from "@tanstack/react-query"
import DefaultLoader from "../default_loader"
import ErrorPage from "../error_page"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { SetStateAction, useState } from "react"
import { Document } from "@/types/document"
import { Checkbox } from "../ui/checkbox"

export default function SelectDocumentsTable({
    limit,
    docs,
    setDocs,
}: {
    limit?: number,
    docs: Document[],
    setDocs: (value: SetStateAction<Document[]>) => void,
}) {
    const { data, status } = useQuery({
        queryKey: ['documents', null],
        queryFn: () => listFolder(),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        return <ErrorPage msg="" />
    }

    return <Table containerClassname="h-fit max-h-[500px] overflow-y-auto relative">
        <TableHeader className="sticky w-full top-0">
            <TableRow>
                <TableHead className="w-[50px]"></TableHead>
                <TableHead className="">Filename</TableHead>
            </TableRow>
        </TableHeader>
        <TableBody>
            {data.documents.map((item, i) => <TableRow
                key={`doc-${i}`}
                onClick={() => {
                    if (docs.map((val) => val.id).indexOf(item.id) === -1) {
                        if (limit !== undefined) {
                            if (docs.length < limit) {
                                setDocs((prev) => prev.concat([item]))
                            }
                        } else {
                            setDocs((prev) => prev.concat([item]))
                        }
                    } else {
                        setDocs((prev) => prev.filter((val) => val.id !== item.id))
                    }
                }}
                className=" hover:cursor-pointer"
            >
                <TableCell className="w-[50px]">
                    <div className="flex items-center">
                        <Checkbox checked={docs.map((val) => val.id).indexOf(item.id) !== -1} />
                    </div>
                </TableCell>
                <TableCell className="font-medium">{item.filename}</TableCell>
            </TableRow>)}
        </TableBody>
    </Table>
}