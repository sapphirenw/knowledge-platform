import { FetchConversation } from '@/api/rag';
import RagClient from './rag_client';
import { ConversationMessage } from '@/types/conversation';

export default async function RAG({
    searchParams,
}: {
    searchParams: { [key: string]: string | string[] | undefined }
}) {
    const convId = searchParams['conversationId'] as string | undefined
    console.log("ConversationId = " + convId)

    const msgs: ConversationMessage[] = []

    if (convId != "" && convId != undefined) {
        // fetch the conversation associated with this
        console.log("fetching conversation ...")
        const response = await FetchConversation(convId)
        if (response.error) {
            console.log("There was an error with the request: " + response.error)
        } else {
            msgs.push(...response.data!.messages)
            console.log(msgs)
        }
        console.log(response)
    }

    return <RagClient convId={convId} msgs={msgs}></RagClient>
}