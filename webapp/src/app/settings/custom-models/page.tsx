import { getAvailableModels, getCustomerLLMs } from "@/actions/llm"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import CustomerModelsClient from "./client"
import CreateCustomerLLM from "./create_model"

export default async function CustomerModels() {
    // pre-fetch the models query
    const queryClient = new QueryClient()

    // get available models
    await queryClient.prefetchQuery({
        queryKey: ['availableModels'],
        queryFn: () => getAvailableModels(""),
    })

    // get customer created models
    await queryClient.prefetchQuery({
        queryKey: ['customerLLMs', false],
        queryFn: () => getCustomerLLMs(false)
    })

    return <div className="w-full space-y-4">
        <CreateCustomerLLM />
        <HydrationBoundary state={dehydrate(queryClient)}>
            <CustomerModelsClient />
        </HydrationBoundary>
    </div>

}