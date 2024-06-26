import { getAllConversations, getConversation } from '@/actions/conversation';
import RagClient from './rag_client';
import { ConversationMessage } from '@/types/conversation';
import ErrorPage from '@/components/error_page';
import { HydrationBoundary, QueryClient, dehydrate } from '@tanstack/react-query';

export default async function RAG() {
    const queryClient = new QueryClient()

    await queryClient.prefetchQuery({
        queryKey: ['conversation'],
        queryFn: getConversation,
    })

    await queryClient.prefetchQuery({
        queryKey: ['allConversations'],
        queryFn: getAllConversations,
    })

    return <HydrationBoundary state={dehydrate(queryClient)}>
        <RagClient></RagClient>
    </HydrationBoundary>
}