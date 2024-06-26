"use client"

import { Conversation } from "@/types/conversation";
import { useState } from "react";
import Cookies from "js-cookie"
import { useQueryClient } from "@tanstack/react-query";

export default function RagSidebarClient({
    conversations,
    activeConvId,
}: {
    conversations: Conversation[],
    activeConvId: string,
}) {
    const queryClient = useQueryClient()
    const [selected, setSelected] = useState(activeConvId)

    const handleClick = async (c: Conversation) => {
        // set the cookie
        Cookies.set("conversationId", c.id, { sameSite: "Strict" })
        setSelected(c.id)

        // invalidate the conversation query
        await queryClient.invalidateQueries({ queryKey: ['conversation'] })
    }

    return <div className="">
        {conversations.map((c, index) => (
            <div key={`conv-${index}`}>
                <button className="w-full" onClick={() => handleClick(c)}>
                    <p className={`py-2 pl-4 text-left w-full rounded-xl hover:bg-secondary ${selected === c.id ? "bg-secondary" : ""}`}>{c.title}</p>
                </button>
            </div>
        ))}
    </div>
}