import { AvailableModel } from "@/types/llm";

export default function AvailableModelView({ model }: { model: AvailableModel }) {
    return <div className="flex items-center justify-between w-full">
        <div className="text-left w-full">
            <p className="font-semibold">{model.displayName}</p>
            <p className="max-w-md">{model.description}</p>
        </div>
        <p className="text-sm bg-primary text-white rounded-full px-2">{model.provider}</p>
    </div>
}