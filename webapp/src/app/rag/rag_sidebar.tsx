import { getAllConversations } from "@/actions/conversation"
import SidebarRow from "./rag_sidebar_row";
import { SquarePlus } from "lucide-react";
import { cookies } from "next/headers";

export default async function Sidebar() {
    try {
        const data = await getAllConversations()

        return <nav className="border-r border-r-border p-4 overflow-y-scroll h-full w-full">
            <div className="w-full">
                <div className="pb-2">
                    <a className="w-full" href={`/rag`}>
                        <div className="py-2 pl-4 w-full rounded-xl hover:bg-secondary">
                            <div className="flex items-center space-x-4">
                                <SquarePlus />
                                <p>New</p>
                            </div>
                        </div>
                    </a>
                </div>
                {data.map((conversation, index) => (
                    <div key={`conv-${index}`}>
                        <SidebarRow c={conversation} activeConvId={cookies().get("conversationId")?.value ?? ""} />
                    </div>
                ))}
            </div>
        </nav>
    } catch (e) {
        if (e instanceof Error) console.log(e)
        return <div className="">ERROR</div>;
    }
}