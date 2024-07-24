import { getAvailableModels } from "@/actions/llm"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import AvailableModelsViewClient from "./client"

export default async function AvailableModelsView() {
    // pre-fetch the models query
    const queryClient = new QueryClient()

    // get available models
    await queryClient.prefetchQuery({
        queryKey: ['availableModels'],
        queryFn: () => getAvailableModels(""),
    })

    return <HydrationBoundary state={dehydrate(queryClient)}>
        <AvailableModelsViewClient />
    </HydrationBoundary>
}