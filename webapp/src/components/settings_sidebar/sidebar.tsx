import { cookies } from "next/headers";
import { AreaChart, Brain, Database, KeyRound, MessageCircleMore } from "lucide-react";
import SettingsSidebarClient, { SidebarRow } from "./sidebar_client";

export default function SettingsSidebar() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")?.value
    const isOpen = cookieStore.get("isSideMenuOpen")?.value ?? "true"

    const items: SidebarRow[] = [
        {
            href: "/rag",
            title: "Chat",
            icon: <MessageCircleMore strokeWidth={2} />
        },
        {
            href: "/settings/models",
            title: "Models",
            icon: <Brain strokeWidth={2} />
        },
        {
            href: "/settings/datastore",
            title: "Datastore",
            icon: <Database strokeWidth={2} />
        },
        {
            href: "/settings/usage",
            title: "Usage",
            icon: <AreaChart strokeWidth={2} />
        },
    ]

    return <SettingsSidebarClient customerId={customerId} items={items} initIsOpen={isOpen} />
}