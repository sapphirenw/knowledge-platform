import { Conversation } from "@/types/conversation"
import Link from "next/link"

export default function SidebarRow({ c }: { c: Conversation }) {

    return <Link href={`/rag?conversationId=${c.id}`}>
        <div className="py-2 px-4 rounded-xl hover:bg-container">{c.title}</div>
    </Link>
}