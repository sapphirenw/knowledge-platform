import { ModelRow } from "@/types/llm";
import { Ellipsis } from "lucide-react";
import { Button } from "./ui/button";

export default function LLMView({ llm }: { llm: ModelRow }) {
    return <div className="w-full px-4 flex items-center justify-between">
        <div className="">
            <p className="text-lg font-semibold">{llm.llm.title}</p>
            <p className="text-sm opacity-50 font-medium">{llm.availableModel.displayName}</p>
        </div>
        <Button variant="secondary">
            <Ellipsis />
        </Button>
    </div>
}