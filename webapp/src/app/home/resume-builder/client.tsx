"use client"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import { useQuery } from "@tanstack/react-query"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"
import { getResumeApplications } from "@/actions/resume"
import Link from "next/link"

export default function ResumeBuilderClient() {

    // const resumeItemResponse = useQuery({
    //     queryKey: ['resumeItem'],
    //     queryFn: () => getResumeItem(),
    // })

    const resumeAppsResponse = useQuery({
        queryKey: ['resumeApplications'],
        queryFn: () => getResumeApplications(),
    })

    const getApplicationsTable = () => {

        if (resumeAppsResponse.status === "pending") {
            return <DefaultLoader />
        }

        if (resumeAppsResponse.status === "error") {
            return <ErrorPage msg="failed" />
        }

        const objs = []
        for (let i = 0; i < resumeAppsResponse.data.length; i++) {
            objs.push(<TableRow key={`app-${i}`}>
                <TableCell className="font-bold">
                    <Link className="w-fit text-primary hover:opacity-50 underline" href="">
                        <p>{resumeAppsResponse.data[i].title}</p>
                    </Link>
                </TableCell>
                <TableCell>
                    <div className="bg-secondary w-fit px-2 rounded-md">
                        <p className="text-primary text-sm font-bold">{resumeAppsResponse.data[i].status}</p>
                    </div>
                </TableCell>
                <TableCell className="text-right">{resumeAppsResponse.data[i].updatedAt.toLocaleString()}</TableCell>
            </TableRow>)
        }

        return <Table containerClassname="h-fit max-h-[500px] overflow-y-auto relative">
            <TableHeader className="sticky w-full top-0">
                <TableRow>
                    <TableHead className="">Title</TableHead>
                    <TableHead className="">Status</TableHead>
                    <TableHead className="text-right">Updated</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {objs.length === 0 ? <TableRow><TableCell><p className="">No applications found.</p></TableCell></TableRow> : objs}
            </TableBody>
        </Table>
    }

    return <div className="">
        {getApplicationsTable()}
    </div>
}