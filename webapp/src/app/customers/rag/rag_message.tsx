import { ConversationMessage } from "@/types/conversation"

export default function RagMessage({
    message
}: {
    message: ConversationMessage
}) {
    // ignore system messages
    if (message.role == 0) {
        return <></>
    }

    const proseClass = "prose prose-slate prose-invert whitespace-pre-line max-w-none"

    const getBody = () => {
        switch (message.role) {
            case 1:
                // user
                return <div className="w-full flex justify-end">
                    <div className="bg-container p-4 rounded-2xl w-fit">
                        <p className={proseClass}>{message.message}</p>
                    </div>
                </div>
            case 2:
                // ai
                return <div className="flex content-start space-x-4">
                    <div className="w-12 h-12 bg-blue-900 rounded-full flex-shrink-0 font-bold text-white grid place-items-center">
                        <p>AI</p>
                    </div>
                    <p className={`${proseClass} prose-lg`}>{message.message}</p>
                </div>
            case 3:
                // tool call
                return <div className="">
                    [TOOL CALL] {message.name}
                </div>
            case 4:
                // tool result
                return <div className="">
                    [TOOL RESULT] {message.message}
                </div>
            default:
                return <div className="">{message.message}</div>
        }
    }
    return <div className="p-4">
        {getBody()}
    </div >
}