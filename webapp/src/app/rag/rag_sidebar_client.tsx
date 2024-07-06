"use client"

import { Conversation } from "@/types/conversation";
import { useState } from "react";
import Cookies from "js-cookie"
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { SquarePlus } from "lucide-react";
import { getAllConversations } from "@/actions/conversation";
import DefaultLoader from "@/components/default_loader";
import ErrorPage from "@/components/error_page";

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

    return <div className="">
        <div className="pb-2">
            <button onClick={() => handleClick("")} className="w-full">
                <div className="py-2 pl-4 w-full rounded-xl hover:bg-secondary">
                    <div className="flex items-center space-x-4">
                        <SquarePlus />
                        <p>New</p>
                    </div>
                </div>
            </button>
        </div>
        {data.map((c, index) => (
            <div key={`conv-${index}`}>
                <button className="w-full" onClick={() => handleClick(c.id)}>
                    <p className={`py-2 pl-4 text-left w-full rounded-xl hover:bg-secondary ${selected === c.id ? "bg-secondary" : ""}`}>{c.title}</p>
                </button>
            </div>
        ))}
    </div>
}