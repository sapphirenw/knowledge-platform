"use client"

import { listFolder } from "@/actions/document"
import ErrorPage from "../error_page"
import { useQuery } from "@tanstack/react-query"
import { Square, SquareCheckBig } from "lucide-react"

import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

import { Button } from "../ui/button"
import { useState } from "react"
import DefaultLoader from "../default_loader"
import { Checkbox } from "../ui/checkbox"
import { Document } from "@/types/document"
import { toast } from "../ui/use-toast"
import SelectDocumentsTable from "./select_doc"

export default function SelectDocumentModel({
    title,
    limit,
    initalSelected,
    onSave,
}: {
    title?: string,
    limit?: number,
    initalSelected?: Document[],
    onSave: (docs: Document[]) => Promise<void>,
}) {
    const [isOpen, setIsOpen] = useState(false)
    const [selected, setSelected] = useState<Document[]>(initalSelected ?? [])
    const [isLoading, setIsLoading] = useState(false)

    const handleClick = async () => {
        if (selected.length === 0) {
            toast({
                variant: "destructive",
                title: "On no!",
                description: <p>You must select at least 1 document.</p>
            })
            return
        }

        setIsLoading(true)
        try {
            await onSave(selected)
            toast({
                title: "Success!",
                description: <p>{`Successfully attached (${selected.length}) documents.`}</p>
            })
            setIsOpen(false)
        } catch {
            toast({
                variant: "destructive",
                title: "On no!",
                description: <p>There was an issue selecting the documents.</p>
            })
        }
        setIsLoading(false)
    }

    const getContent = () => {
        if (isOpen) {
            return <SelectDocumentsTable
                limit={limit}
                docs={selected}
                setDocs={setSelected}
            />
        }

        return <DefaultLoader />
    }

    return <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogTrigger asChild>
            <Button onClick={() => setIsOpen(true)}>{title ?? "Select Documents"}</Button>
        </DialogTrigger>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>{title ?? "Select Documents"}</DialogTitle>
                <DialogDescription>
                </DialogDescription>
            </DialogHeader>
            {getContent()}
            <DialogFooter>
                <Button onClick={() => handleClick()}>
                    <div className="flex space-x-2 items-center">
                        {isLoading ? <DefaultLoader /> : <></>}
                        <p>Submit</p>
                    </div>
                </Button>
            </DialogFooter>
        </DialogContent>

    </Dialog>
}