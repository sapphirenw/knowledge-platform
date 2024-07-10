import { getAvailableModels, getCustomerLLMs } from "@/actions/llm"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import CustomerModelsClient from "./models_client"
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

    return <HydrationBoundary state={dehydrate(queryClient)}>
        <CreateCustomerLLM />
        <CustomerModelsClient />
    </HydrationBoundary>

}