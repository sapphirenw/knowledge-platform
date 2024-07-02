"use client"

import { getWebsites } from "@/actions/websites"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { useQuery } from "@tanstack/react-query"

export default function UserWebsites() {
    // fetch the websites
    const siteResponse = useQuery({
        queryKey: ['websites'],
        queryFn: () => getWebsites(),
    })

    if (siteResponse.status === "error") {
        return <ErrorPage msg="" />
    }

    if (siteResponse.status === "pending") {
        return <DefaultLoader />
    }

    return <div className="w-full">
        <Table containerClassname="h-fit max-h-[500px] overflow-y-auto relative">
            <TableHeader className="sticky w-full top-0">
                <TableRow>
                    <TableHead className="">Domain</TableHead>
                    <TableHead className="">Page Count</TableHead>
                    <TableHead className="text-right">Created</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {siteResponse.data!.map((item, i) => <TableRow key={`site-${i}`}>
                    <TableCell className="font-medium">{item.domain}</TableCell>
                    <TableCell className="">{item.pageCount}</TableCell>
                    <TableCell className="text-right">{new Date(item.createdAt!).toLocaleString()}</TableCell>
                </TableRow>)}
            </TableBody>
        </Table>
    </div>
}