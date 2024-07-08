import { ConversationMessage } from "@/types/conversation"
import MessageToolCallResult from "./msg_tool_call_result"
import MessageToolCall from "./msg_tool_call"


export default function RagMessage({
    message,
    offset
}: {
    message: ConversationMessage
    offset: number
}) {
    // ignore system messages
    if (message.role == 0) {
        return <></>
    }

    const proseClass = "prose prose-slate prose-invert whitespace-pre-line max-w-none"

    const getMessage = () => {
        switch (message.role) {
            case 1:
                // user
                return <div className="w-full flex justify-end">
                    <div className="bg-secondary p-4 rounded-2xl w-fit max-w-lg">
                        <p className={proseClass}>{message.message}</p>
                    </div>
                </div>
            case 2:
                // ai
                return <div className="flex content-start space-x-4">
                    <div className="w-12 h-12 bg-blue-900 rounded-full flex-shrink-0 font-bold text-white grid place-items-center">
                        <p>AI</p>
                    </div>
                    <div className={`${proseClass} prose-lg`}>{message.message}</div>
                </div>
            case 3:
                // tool call
                return <MessageToolCall message={message} offset={offset} />
            case 4:
                return <MessageToolCallResult message={message} offset={offset} />
            default:
                return <div className="">{message.message}</div>
        }
    }

    return <div className="px-4 py-2">{getMessage()}</div>
}

