"use client"

import { getCustomerLLMs } from "@/actions/llm"
import DefaultLoader from "@/components/default_loader"
import LLMView from "@/components/llm"
import { useQuery } from "@tanstack/react-query"

export default function CustomerModelsClient() {
    return <div className="w-full">
        <CustomerLLMsView />
    </div>
}

function CustomerLLMsView() {

    const { status, data, error } = useQuery({
        queryKey: ['customerLLMs', false],
        queryFn: () => getCustomerLLMs(false)
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.error(error)
        return <p>Failed to get your custom models</p>
    }

    return <div className="border border-border rounded-lg w-full">
        {data.map((item, i) => <LLMView key={item.llm.id} llm={item} />)}
    </div>
}