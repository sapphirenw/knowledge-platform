import { cookies } from "next/headers";
import { AreaChart, Brain, Database, KeyRound, MessageCircleMore } from "lucide-react";
import SidebarClient, { SidebarRow } from "./sidebar_client";

export default function Sidebar() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")?.value
    const isOpen = cookieStore.get("isSideMenuOpen")?.value ?? "true"

    const items: SidebarRow[] = [
        {
            href: "/rag",
            title: "Chat",
            icon: <MessageCircleMore strokeWidth={1.5} />
        },
        {
            href: "/models",
            title: "Models",
            icon: <Brain strokeWidth={1.5} />
        },
        {
            href: "/datastore",
            title: "Datastore",
            icon: <Database strokeWidth={1.5} />
        },
        {
            href: "/usage",
            title: "Usage",
            icon: <AreaChart strokeWidth={1.5} />
        },
    ]

    return <SidebarClient customerId={customerId} items={items} initIsOpen={isOpen} />
}