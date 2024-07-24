import { cookies } from "next/headers";
import SettingsSidebarClient from "./sidebar_client";

export default function SettingsSidebar() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")?.value
    const isOpen = cookieStore.get("isSideMenuOpen")?.value ?? "true"

    return <SettingsSidebarClient customerId={customerId} initIsOpen={isOpen} />
}