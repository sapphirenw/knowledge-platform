import { getAllConversations, getConversation } from '@/actions/conversation';
import RagClient from './rag_client';
import { ConversationMessage } from '@/types/conversation';
import ErrorPage from '@/components/error_page';
import { HydrationBoundary, QueryClient, dehydrate } from '@tanstack/react-query';
import { cookies } from 'next/headers';
import { getAvailableLLMs } from '@/actions/llm';

export default async function RAG() {
    const queryClient = new QueryClient()

    await queryClient.prefetchQuery({
        queryKey: ['conversation'],
        queryFn: getConversation,
    })

    await queryClient.prefetchQuery({
        queryKey: ['availableLLMs'],
        queryFn: getAvailableLLMs,
    })

    const cid = cookies().get("cid")?.value

    return <HydrationBoundary state={dehydrate(queryClient)}>
        <RagClient wsBaseUrl={`${process.env.DB_HOST}/v1/customers/${cid}/rag2`}></RagClient>
    </HydrationBoundary>
}