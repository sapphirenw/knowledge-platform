import LoaderGrid from "@/components/loaders/loader_grid"
import { ConversationMessage } from "@/types/conversation"

export default function MessageToolCall({
    message,
    offset
}: {
    message: ConversationMessage
    offset: number
}) {

    const getLoadingText = () => {
        switch (message.name) {
            case "vector_query":
                return "Searching local information ..."
            default:
                return ""
        }
    }

    const getBody = () => {
        if (offset > 1) {
            return <p className="opacity-60 text-sm pl-4 pt-10">Message composed from the following information:</p>
        }

        return <div className="flex items-center space-x-4">
            <LoaderGrid />
            <p className="opacity-50 text-sm font-semibold">{!offset ? getLoadingText() : "Crafting response ..."}</p>
        </div>
    }

    // return a loading indicator
    return <div className="h-[50px] flex items-center">
        {getBody()}
    </div>
}