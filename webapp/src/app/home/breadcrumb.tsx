"use client"

import { usePathname } from "next/navigation"
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Slash } from "lucide-react"

export default function HomeBreadcrumb() {
    const pathname = usePathname()
    const items = pathname!.split("/").filter((val) => val.trim().length != 0)

    return <Breadcrumb>
        <BreadcrumbList>
            <BreadcrumbSeparator><Slash /></BreadcrumbSeparator>
            {items.map((item, i) => <>
                <BreadcrumbItem>
                    <BreadcrumbLink href={`/${items.slice(0, i + 1).join("/")}`}>{item.replaceAll("-", " ").split(" ").map((i) => (i.charAt(0).toUpperCase() + i.slice(1))).join(" ")}</BreadcrumbLink>
                </BreadcrumbItem>
                {i < items.length - 1 ? <BreadcrumbSeparator><Slash /></BreadcrumbSeparator> : null}
            </>)}
        </BreadcrumbList>
    </Breadcrumb>
}