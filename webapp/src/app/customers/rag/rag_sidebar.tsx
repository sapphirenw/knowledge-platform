import { GetAllConversations } from "@/actions/conversation"
import { Conversation } from "@/types/conversation"
import Link from "next/link";
import { redirect } from 'next/navigation'
import SidebarRow from "./rag_sidebar_row";

export default async function Sidebar() {
    const getSidebar = async () => {
        const response = await GetAllConversations();
        if (response.error) {
            console.log("Server error: ", response.error);
            return <div className="">ERROR</div>;
        } else {
            return (
                <div className="w-max">
                    {response.data!.map((conversation, index) => (
                        <div key={`conv-${index}`}>
                            <SidebarRow c={conversation} />
                        </div>
                    ))}
                </div>
            );
        }
    };

    return (
        <nav className="bg-bg-dark border-r border-r-border p-4 overflow-y-scroll h-full">
            {await getSidebar()}
        </nav>
    );
}