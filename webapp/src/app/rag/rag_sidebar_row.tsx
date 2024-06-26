"use client"

import { Conversation } from "@/types/conversation"
import Cookies from "js-cookie"

export default function SidebarRow({ c, activeConvId }: { c: Conversation, activeConvId: string }) {
    const handleClick = () => {
        // set the cookie
        Cookies.set("conversationId", c.id)

        // reload the page
        location.reload();
    }

    return <button onClick={handleClick} className="w-full">
        <p className={`py-2 pl-4 text-left w-full rounded-xl hover:bg-secondary ${activeConvId === c.id ? "bg-secondary" : ""}`}>{c.title}</p>
    </button>
}