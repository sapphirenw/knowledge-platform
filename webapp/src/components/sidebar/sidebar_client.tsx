"use client"

import { KeyRound, PanelLeftClose, PanelLeftOpen } from "lucide-react"
import { useState } from "react"
import Cookies from "js-cookie"

export type SidebarRow = {
    href: string
    title: string
    icon: JSX.Element
}

export default function SidebarClient({
    customerId,
    items,
    initIsOpen
}: {
    customerId: string | undefined,
    items: SidebarRow[],
    initIsOpen: string,
}) {

    const [isOpen, setIsOpen] = useState(initIsOpen == "true")

    const toggleIsOpen = () => {
        Cookies.set("isSideMenuOpen", `${!isOpen}`)
        setIsOpen(!isOpen)
    }

    return <nav className={`border-r border-r-border h-full ${isOpen ? "w-[200px]" : "w-[60px]"} transition-all duration-300`}>
        <div className="w-full h-full p-2">
            <div className="flex flex-col justify-between h-full">
                <div className="space-y-2">
                    <div className="opacity-50 px-2">
                        <button onClick={() => toggleIsOpen()}>
                            {isOpen ? <PanelLeftClose strokeWidth={1.5} /> : <PanelLeftOpen strokeWidth={1.5} />}
                        </button>
                    </div>
                    {customerId === undefined ? null : <div className="space-y-2">
                        {items.map((item, i) => <SidebarRowView key={`item-${i}`} item={item} isOpen={isOpen} />)}
                    </div>}
                </div>
                <SidebarRowView item={{
                    href: "/login",
                    title: "Login",
                    icon: <KeyRound strokeWidth={1.5} />
                }} isOpen={isOpen} />
            </div>
        </div>
    </nav>
}

function SidebarRowView({
    item,
    isOpen,
}: {
    item: SidebarRow,
    isOpen: boolean,
}) {
    return <div className="">
        <a href={item.href}>
            <div className="pl-2 pr-4 py-2 border border-border rounded-md hover:bg-border flex items-center space-x-2">
                <div className="w-[25px]">{item.icon}</div>
                {isOpen ? <p>{item.title}</p> : null}
            </div>
        </a>
    </div>
}