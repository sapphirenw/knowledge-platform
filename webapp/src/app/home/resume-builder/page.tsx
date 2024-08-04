import ResumeBuilderClient from "./client"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import { getResumeApplications, getResumeChecklist } from "@/actions/resume"
import CreateResumeApplicationButton from "./create"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { Separator } from "@/components/ui/separator"
import { Square, SquareCheckBig } from "lucide-react"
import ResumeChecklistView from "./checklist"

export default async function ResumeBuilder() {
    const queryClient = new QueryClient()

    // get the resume item
    // await queryClient.prefetchQuery({
    //     queryKey: ['resumeItem'],
    //     queryFn: getResumeItem,
    // })

    // get the checklist
    await queryClient.prefetchQuery({
        queryKey: ['resumeChecklist'],
        queryFn: getResumeChecklist,
    })

    // get the resume applications
    await queryClient.prefetchQuery({
        queryKey: ['resumeApplications'],
        queryFn: getResumeApplications,
    })

    return <HydrationBoundary state={dehydrate(queryClient)}>
        <div className="space-y-16">
            <ResumeChecklistView />
            <div className="space-y-2">
                <div className="flex items-center justify-between">
                    <Label htmlFor="user_documents">
                        <h3 className="text-lg font-bold">Applications</h3>
                    </Label>
                    <CreateResumeApplicationButton />
                </div>
                <ResumeBuilderClient />
            </div>
        </div>
    </HydrationBoundary>
}