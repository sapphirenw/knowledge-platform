import { GetConversation } from '@/api/conversation';
import RagClient from './rag_client';
import { ConversationMessage } from '@/types/conversation';

export default async function RAG({
    searchParams,
}: {
    searchParams?: { conversationId?: string; }
}) {
    const convId = searchParams?.conversationId

    const msgs: ConversationMessage[] = []
    if (convId == "new") { } else if (convId != undefined && convId != "") {
        // fetch the conversation associated with this
        console.log("fetching conversation ...")
        const response = await GetConversation(convId)
        if (response.error) {
            console.log("There was an error with the request: " + response.error)
        } else {
            msgs.push(...response.data!.messages)
        }
    }

    console.log("ConversationId = " + convId)


    return <RagClient convId={convId} msgs={msgs}></RagClient>
}