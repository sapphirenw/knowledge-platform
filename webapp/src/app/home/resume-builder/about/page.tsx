import { getResumeAbout } from "@/actions/resume"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import ResumeAboutClient from "./client"

export default async function ResumeAbout() {
    const queryClient = new QueryClient()

    await queryClient.prefetchQuery({
        queryKey: ['resumeAbout'],
        queryFn: getResumeAbout,
    })

    return <HydrationBoundary state={dehydrate(queryClient)}>
        <ResumeAboutClient />
    </HydrationBoundary>
}