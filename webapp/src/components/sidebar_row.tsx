import { cn } from "@/lib/utils"
import { Button } from "./ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./ui/tooltip"

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
                    variant={variant ?? (path.match(item.href) ? "secondary" : "ghost")}
                    className={`w-full ${path.match(item.href) ? "font-semibold" : ""}`}
                    // size={isOpen ? "default" : "icon"}
                    asChild
                >
                    <a href={item.href}>
                        <div className={cn(`flex items-center ${isOpen ? "w-full text-left" : ""}`, className)}>
                            <div className={`${iconClass}`}>{item.icon}</div> {isOpen ? <p className="ml-2">{item.title}</p> : null}
                        </div>
                    </a>
                </Button>
            </TooltipTrigger>
            <TooltipContent>
                <p>{item.title}</p>
            </TooltipContent>
        </Tooltip>
    </TooltipProvider>
}