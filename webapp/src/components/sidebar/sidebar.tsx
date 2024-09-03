import { cookies } from "next/headers"
import SidebarClient from "./sidebar_client"
import { SidebarRow } from "./sidebar_row"

export type SidebarGroup = {
    title: string
    titleSm: string
    items: SidebarRow[]
}

export type SidebarProps = {
    header?: SidebarRow
    groups: SidebarGroup[]
    footer?: JSX.Element
    allowsClose: boolean
}

export default function Sidebar({
    props,
}: {
    props: SidebarProps,
}) {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")?.value
    const isOpen = cookieStore.get("isSideMenuOpen")?.value ?? "true"

    return <SidebarClient props={props} customerId={customerId} initIsOpen={isOpen} />
}