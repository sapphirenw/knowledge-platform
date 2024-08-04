import { cn } from "@/lib/utils"
import { Button } from "./ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./ui/tooltip"
import Link from "next/link"

export type SidebarRow = {
    href: string
    title: string
    icon: JSX.Element
}

export default function SidebarRowView({
    item,
    path,
    isOpen,
    className,
    variant,
    iconClass,
}: {
    item: SidebarRow,
    path: string,
    isOpen: boolean,
    className?: string,
    variant?: "link" | "default" | "destructive" | "outline" | "secondary" | "ghost",
    iconClass?: string,
}) {
    return <TooltipProvider>
        <Tooltip>
            <TooltipTrigger asChild>
                <Button
                    variant={(path.match(item.href) ? "secondary" : variant ?? "ghost")}
                    className={`w-full ${path.match(item.href) ? "font-semibold" : ""}`}
                    size={isOpen ? "default" : "icon"}
                    asChild
                >
                    <Link href={item.href}>
                        <div className={cn(`flex items-center w-full ${isOpen ? "" : "grid place-items-center"}`, className)}>
                            <div className={`${iconClass}`}>{item.icon}</div>
                            {isOpen ? <p className="ml-2">{item.title}</p> : null}
                        </div>
                    </Link>
                </Button>
            </TooltipTrigger>
            <TooltipContent>
                <p>{item.title}</p>
            </TooltipContent>
        </Tooltip>
    </TooltipProvider>
}