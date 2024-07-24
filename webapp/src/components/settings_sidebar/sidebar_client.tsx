"use client"

import { AreaChart, BetweenHorizontalStart, Brain, Database, Earth, FileText, Home, KeyRound, Package, PanelLeftClose, PanelLeftOpen } from "lucide-react"
import { useState } from "react"
import Cookies from "js-cookie"
import { usePathname } from 'next/navigation'
import { Button } from "../ui/button"
import Link from "next/link"
import SidebarRowView from "../sidebar_row"
import { ThemeToggle } from "../theme_toggle"

export default function SettingsSidebarClient({
    customerId,
    initIsOpen
}: {
    customerId: string | undefined,
    initIsOpen: string,
}) {
    const pathname = usePathname()

    const [isOpen, setIsOpen] = useState(initIsOpen == "true")

    const toggleIsOpen = () => {
        Cookies.set("isSideMenuOpen", `${!isOpen}`)
        setIsOpen(!isOpen)
    }

    const sidebar = () => {
        if (customerId === undefined) {
            return null
        }

        return <div className="space-y-2">
            <p className={`font-medium text-sm opacity-75 pt-4 ${isOpen ? "pl-4" : "text-center"}`}>{isOpen ? "Datastore" : "Data"}</p>
            <div className="">
                <SidebarRowView item={{
                    href: "/settings/vector-requests",
                    title: "Vector Requests",
                    icon: <BetweenHorizontalStart size={16} />
                }} path={pathname} isOpen={isOpen} />
                <SidebarRowView item={{
                    href: "/settings/documents",
                    title: "Documents",
                    icon: <FileText size={16} />
                }} path={pathname} isOpen={isOpen} />
                <SidebarRowView item={{
                    href: "/settings/websites",
                    title: "Websites",
                    icon: <Earth size={16} />
                }} path={pathname} isOpen={isOpen} />
            </div>
            <p className={`font-medium text-sm opacity-75 pt-4 ${isOpen ? "pl-4" : "text-center"}`}>{isOpen ? "General" : "Gen"}</p>
            <div className="">
                <SidebarRowView item={{
                    href: "/settings/available-llms",
                    title: "Available LLMs",
                    icon: <Package size={16} />
                }} path={pathname} isOpen={isOpen} />
                <SidebarRowView item={{
                    href: "/settings/custom-models",
                    title: "Custom Models",
                    icon: <Brain size={16} />
                }} path={pathname} isOpen={isOpen} />
                <SidebarRowView item={{
                    href: "/settings/usage",
                    title: "Usage",
                    icon: <AreaChart size={16} />
                }} path={pathname} isOpen={isOpen} />
            </div>
        </div>
    }

    return <nav className={`border-r border-r-border h-full ${isOpen ? "w-[200px]" : "w-[60px]"} transition-all duration-300`}>
        <div className="w-full h-full p-2">
            <div className="flex flex-col justify-between h-full">
                <div className="space-y-2">
                    <div className="opacity-50 px-2 text-right">
                        <button onClick={() => toggleIsOpen()}>
                            {isOpen ? <PanelLeftClose strokeWidth={1.5} /> : <PanelLeftOpen strokeWidth={1.5} />}
                        </button>
                    </div>
                    <SidebarRowView
                        item={{
                            href: "/home",
                            title: "Home",
                            icon: <Home size={16} />
                        }}
                        path={pathname}
                        isOpen={isOpen}
                        className="text-center w-fit"
                        variant="default"
                    />
                    {sidebar()}
                </div>
                <div className="flex items-center space-x-2">
                    <SidebarRowView
                        item={{
                            href: "/login",
                            title: "Sign-out",
                            icon: <KeyRound size={16} />
                        }}
                        path={pathname}
                        isOpen={isOpen}
                        variant="outline"
                    />
                    <div className="aspect-square">
                        <ThemeToggle />
                    </div>
                </div>
            </div>
        </div>
    </nav>
}