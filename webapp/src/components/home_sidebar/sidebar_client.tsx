"use client"

import { FileSpreadsheet, MessageSquareText, PanelLeftClose, PanelLeftOpen, Settings } from "lucide-react"
import { useState } from "react"
import Cookies from "js-cookie"
import { usePathname } from 'next/navigation'
import SidebarRowView from "../sidebar_row"

export default function HomeSidebarClient({
    customerId,
    initIsOpen
}: {
    customerId: string | undefined,
    initIsOpen: string,
}) {
    const pathname = usePathname()

    const [isOpen, setIsOpen] = useState(false)

    const toggleIsOpen = () => {
        setIsOpen(!isOpen)
    }

    const sidebar = () => {
        if (customerId === undefined) {
            return null
        }

        return <div className="space-y-2">
            <div className="">
                <SidebarRowView item={{
                    href: "/home/chat",
                    title: "Basic Chat",
                    icon: <MessageSquareText size={16} />
                }} path={pathname} isOpen={isOpen} variant="outline" />
            </div>
            <p className={`font-medium text-sm opacity-75 pt-4 pl-3`}>{isOpen ? "Integrations" : "Int"}</p>
            <div className="">
                <SidebarRowView item={{
                    href: "/home/resume-builder",
                    title: "Resume Builder",
                    icon: <FileSpreadsheet size={16} />
                }} path={pathname} isOpen={isOpen} variant="outline" />
            </div>
        </div>
    }

    // return <nav className={`border-r border-r-border h-full ${isOpen ? "w-[200px]" : "w-[60px]"} transition-all duration-300`}>
    return <nav
        onMouseEnter={() => setIsOpen(true)}
        onMouseLeave={() => setIsOpen(false)}
        className={`border-r border-r-border h-full hover:w-[200px] w-[60px] transition-all duration-300`}>
        <div className="w-full h-full p-2">
            <div className="flex flex-col justify-between h-full">
                <div className="space-y-2 pt-[50px]">
                    {/* <div className="opacity-50 px-2 text-right">
                        <button onClick={() => toggleIsOpen()}>
                            {isOpen ? <PanelLeftClose strokeWidth={1.5} /> : <PanelLeftOpen strokeWidth={1.5} />}
                        </button>
                    </div> */}
                    {sidebar()}
                </div>
                <SidebarRowView
                    item={{
                        href: "/settings",
                        title: "Settings",
                        icon: <Settings size={16} />
                    }}
                    path={pathname}
                    isOpen={isOpen}
                    variant="outline"
                />
            </div>
        </div>
    </nav>
}