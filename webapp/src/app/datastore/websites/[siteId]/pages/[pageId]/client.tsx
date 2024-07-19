"use client"

import { WebsitePageContentResponse } from "@/types/websites"
import { Separator } from "@/components/ui/separator"

export default function WebsitePageViewClient({
    resp,
}: {
    resp: WebsitePageContentResponse,
}) {

    const getCleaned = () => {
        if (resp.cleaned === undefined) {
            return null
        }
        return <div className="">
            <h4 className="font-medium">Cleaned Content:</h4>
            <Separator />
            <p className="break-all">{resp.cleaned}</p>
        </div>
    }

    const getChunked = () => {
        if (resp.chunks === undefined) {
            return null
        }
        return <div className="">
            <h4 className="font-medium">Chunks <span className="text-sm opacity-50">- {resp.chunks.length}</span></h4>
            <Separator />
            <div className="space-y-2">
                {resp.chunks.map((item, i) => <p key={`chunk-${i}`}>{item}</p>)}
            </div>
        </div>
    }

    return <div className="w-full space-y-8">
        <h3 className="text-lg font-bold">{resp.page.url}</h3>
        {getCleaned()}
        {getChunked()}
    </div>
}