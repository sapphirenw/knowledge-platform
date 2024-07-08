"use client"

import { getAvailableLLMs } from "@/actions/llm"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"
import { Button } from "@/components/ui/button"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { ModelRow } from "@/types/llm"
import { useQuery } from "@tanstack/react-query"
import { ChevronDown } from "lucide-react"

export default function RagLLMSelector({ currLLM, onSelect }: { currLLM?: ModelRow, onSelect: (model: ModelRow) => void }) {
    const { data, status, error } = useQuery({
        queryKey: ['availableLLMs'],
        queryFn: () => getAvailableLLMs(),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.error(error)
        return <ErrorPage msg="failed to get the llms" />
    }

    const getRows = () => {
        const items = []

        for (let i = 0; i < data.length; i++) {
            items.push(<button key={data[i].llm.id} className="w-full" onClick={() => onSelect(data[i])}>
                <RagLLMSelectorRow currLLM={currLLM} model={data[i]} />
            </button>)
        }
        return items
    }

    return <Popover>
        <PopoverTrigger asChild>
            <Button variant="outline">
                <div className="flex items-center space-x-2">
                    <p>{currLLM?.llm?.title ?? "Model"}</p>
                    <ChevronDown />
                </div>
            </Button>
        </PopoverTrigger>
        <PopoverContent className="w-80">
            <div className="">
                {getRows()}
            </div>
        </PopoverContent>
    </Popover>
}

function RagLLMSelectorRow({ currLLM, model }: { currLLM?: ModelRow, model: ModelRow }) {
    return <div className={`w-full hover:bg-secondary rounded-md transition-colors border ${(currLLM?.llm?.id ?? "") == model.llm.id ? "border-border" : "border-transparent"}`}>
        <div className="text-left px-4 py-2">
            <p className="">{model.llm.title}</p>
            <p className="text-sm opacity-50 font-medium">{model.availableModel.displayName}</p>
        </div>
    </div>
}