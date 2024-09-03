"use client"

import { PanelLeftClose, PanelLeftOpen } from "lucide-react"
import { useState } from "react"
import { usePathname } from 'next/navigation'
import SidebarRowView from "./sidebar_row"
import { SidebarProps } from "./sidebar"
import Cookies from "js-cookie"

export default function SidebarClient({
    props,
    customerId,
    initIsOpen,
}: {
    props: SidebarProps,
    customerId: string | undefined,
    initIsOpen: string,
}) {
    const [isOpen, setIsOpen] = useState(props.allowsClose ? initIsOpen == "true" : true)

    const toggleIsOpen = () => {
        Cookies.set("isSideMenuOpen", `${!isOpen}`)
        setIsOpen(!isOpen)
    }

    const sidebar = () => {
        if (customerId === undefined) {
            return null
        }

        return <div className="space-y-2">
            {props.groups.map((group, i) => <div key={`group-${i}`} className="space-y-2">
                <p className={`font-medium text-sm opacity-75 pt-4 line-clamp-1 ${isOpen ? "pl-4" : "text-center"}`}>{isOpen ? group.title : group.titleSm}</p>
                <div className="">
                    {group.items.map((item, j) => <SidebarRowView
                        key={`item-${j}`}
                        item={item}
                        isOpen={isOpen}
                    />)}
                </div>
            </div>)}
        </div>
    }

    return <nav className={`border-r border-r-border h-full ${isOpen ? "w-[200px]" : "w-[60px]"} transition-all duration-300`}>
        <div className="w-full h-full p-2">
            <div className="flex flex-col justify-between h-full">
                <div className="space-y-2">
                    {props.allowsClose ? <div className="opacity-50 px-2 text-right">
                        <button onClick={() => toggleIsOpen()}>
                            {isOpen ? <PanelLeftClose strokeWidth={1.5} /> : <PanelLeftOpen strokeWidth={1.5} />}
                        </button>
                    </div> : null}
                    {props.header !== undefined ? <SidebarRowView
                        item={props.header}
                        isOpen={isOpen}
                    /> : null}
                    {sidebar()}
                </div>
                {props.footer}
            </div>
        </div>
    </nav>
}