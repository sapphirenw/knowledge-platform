import { cookies } from "next/headers";
import HomeSidebarClient from "./sidebar_client";

export default function HomeSidebar() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")?.value
    const isOpen = cookieStore.get("isSideMenuOpen")?.value ?? "true"

    return <HomeSidebarClient customerId={customerId} initIsOpen={isOpen} />
}