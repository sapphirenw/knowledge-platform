import { getAllConversations } from "@/actions/conversation"
import { SquarePlus } from "lucide-react";
import { cookies } from "next/headers";
import RagSidebarClient from "./rag_sidebar_client";
import { QueryClient } from "@tanstack/react-query";

export default async function Sidebar() {
    try {
        const queryClient = new QueryClient()
        await queryClient.prefetchQuery({
            queryKey: ['allConversations'],
            queryFn: getAllConversations,
        })

        return <nav className="border-r border-r-border p-4 overflow-y-scroll h-full w-full">
            <div className="w-full">
                <RagSidebarClient activeConvId={cookies().get("conversationId")?.value ?? ""} />
            </div>
        </nav>
    } catch (e) {
        if (e instanceof Error) console.log(e)
        return <div className="">ERROR</div>;
    }
}