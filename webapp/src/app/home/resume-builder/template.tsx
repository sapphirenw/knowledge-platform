import { Slash } from "lucide-react"

import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { headers } from "next/headers";

export default function Template({ children }: { children: React.ReactNode }) {
    const headersList = headers()
    const pathname = headersList.get('x-pathname') ?? "resume-builder";

    const items = pathname!.split("/").filter((val) => val.trim().length != 0)

    console.log(items)

    return <div className="overflow-scroll">
        <div className="p-8">
            <div className="safe-area">
                <Breadcrumb>
                    <BreadcrumbList>
                        {items.map((item, i) => <>
                            <BreadcrumbItem>
                                <BreadcrumbLink href={`/${items.slice(0, i + 1).join("/")}`}>{item.replaceAll("-", " ").split(" ").map((i) => (i.charAt(0).toUpperCase() + i.slice(1))).join(" ")}</BreadcrumbLink>
                            </BreadcrumbItem>
                            {i < items.length - 1 ? <BreadcrumbSeparator><Slash /></BreadcrumbSeparator> : null}
                        </>)}
                    </BreadcrumbList>
                </Breadcrumb>
            </div>
        </div>
        <div className="p-4">
            <div className="safe-area">
                {children}
            </div>
        </div>
    </div>
}