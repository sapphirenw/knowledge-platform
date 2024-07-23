"use client"

import { Website, WebsitePage } from "@/types/websites"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import Link from "next/link"
import { useState } from "react"
import { deleteWebsite } from "@/actions/websites"
import { toast } from "@/components/ui/use-toast"
import {
    AlertDialog,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import DefaultLoader from "@/components/default_loader"

export default function WebsiteViewClient({
    site,
    pages,
}: {
    site: Website,
    pages: WebsitePage[],
}) {
    const [openDialog, setOpenDialog] = useState(false)
    const [isLoading, setIsLoading] = useState(false)

    const deleteSite = async () => {
        setIsLoading(true)
        try {
            await deleteWebsite(site.id)
            setOpenDialog(false)
            toast({
                title: "Success!",
                description: <p>Successfully deleted the website</p>
            })
        } catch {
            toast({
                variant: "destructive",
                title: `Failed to delete the site`,
            })
        }
        setIsLoading(false)
    }

    const getPagesTable = () => {
        return <Table containerClassname="">
            <TableHeader className="sticky w-full top-0">
                <TableRow>
                    <TableHead className="">URL</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {pages.map((item, i) => <TableRow key={`page-${i}`}>
                    <TableCell className="font-medium">
                        <Link className="text-primary hover:opacity-50 underline" href={`/datastore/websites/${site.id}/pages/${item.id}`}>{item.url}</Link>
                    </TableCell>
                </TableRow>)}
            </TableBody>
        </Table>
    }

    return <div className="w-full space-y-8">
        <h3 className="text-lg font-bold">{`${site.domain}${site.path}`}</h3>
        <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
            <AlertDialogTrigger asChild>
                <Button variant="destructive">Delete Website</Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <AlertDialogDescription>
                        <div className="space-y-4">
                            <p>This action is permament, and will require you to re-scrape the website in order to access any of the content.</p>
                        </div>
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Close</AlertDialogCancel>
                    <Button variant="destructive" onClick={() => deleteSite()}>
                        <div className="flex items-center space-x-2">
                            {isLoading ? <DefaultLoader /> : null}
                            <p>Delete</p>
                        </div>
                    </Button>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
        {getPagesTable()}
    </div>
}