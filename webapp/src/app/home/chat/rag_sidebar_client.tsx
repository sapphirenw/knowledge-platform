"use client"

import { Conversation } from "@/types/conversation";
import { useState } from "react";
import Cookies from "js-cookie"
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { SquarePlus } from "lucide-react";
import { getAllConversations } from "@/actions/conversation";
import DefaultLoader from "@/components/default_loader";
import ErrorPage from "@/components/error_page";
import { Button } from "@/components/ui/button";

export default function RagSidebarClient({
    activeConvId,
}: {
    activeConvId: string,
}) {
    const queryClient = useQueryClient()

    const [selected, setSelected] = useState(activeConvId)

    const { status, data, error } = useQuery({
        queryKey: ['allConversations'],
        queryFn: () => getAllConversations(),
    })

    const handleClick = async (convId: string) => {
        // set the cookie
        Cookies.set("conversationId", convId, { sameSite: "Strict" })
        setSelected(convId)

        // invalidate the conversation query
        await queryClient.invalidateQueries({ queryKey: ['conversation'] })
    }

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.error(error)
        return <ErrorPage msg="" />
    }

    return <div className="space-y-2">
        <div className="">
            <Button
                className="w-full text-left"
                variant="outline"
                onClick={() => handleClick("")}
            >
                <div className="flex items-center w-full text-left">
                    <SquarePlus size={16} />
                    <p className="ml-2">New</p>
                </div>
            </Button>
        </div>
        <p className={`font-medium text-sm opacity-75 pt-4`}>History</p>
        <div className="">
            {data.map((c, index) => (
                <div key={`conv-${index}`}>
                    <Button
                        variant={selected ? "secondary" : "ghost"}
                        className="w-full text-left"
                        onClick={() => handleClick(c.id)}
                    >
                        <p className="w-full overflow-clip overflow-ellipsis">{c.title}</p>
                    </Button>
                </div>
            ))}
        </div>
    </div>
}