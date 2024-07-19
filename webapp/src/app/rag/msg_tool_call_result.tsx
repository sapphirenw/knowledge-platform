import { Button } from "@/components/ui/button";
import { ConversationMessage } from "@/types/conversation"
import { Document } from "@/types/document"
import { WebsitePage } from "@/types/websites";
import { File, FileText } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

export default function MessageToolCallResult({
    message,
    offset
}: {
    message: ConversationMessage
    offset: number
}) {
    const getItems = () => {
        const items: JSX.Element[] = []
        if (message.arguments == undefined) {
            return items
        }

        switch (message.name) {
            case "vector_query":
                for (let i = 0; i < message.arguments.docs.length; i++) {
                    items.push(<DocumentItem key={`doc-${i}`} doc={message.arguments.docs[i]} />)
                }
                for (let i = 0; i < message.arguments.pages.length; i++) {
                    items.push(<WebsitePageItem key={`page-${i}`} page={message.arguments.pages[i]} />)
                }
        }

        return items
    }

    return <div className="grid grid-cols-2 gap-2">
        {getItems()}
    </div>
}

function DocumentItem({ doc }: { doc: Document }) {
    return <Link href={`/datastore/documents/${doc.id}`} className="bg-secondary hover:opacity-75 transition-opacity rounded-lg grid place-items-center">
        <div className="flex items-center truncate w-full text-left">
            <div className="w-[20px] mr-2">
                <FileText size={20} />
            </div>
            <p className="truncate">{doc.filename}</p>
        </div>
    </Link>
}

function WebsitePageItem({ page }: { page: WebsitePage }) {
    const urlObj = new URL(page.url);

    return <a href={page.url} target="_blank" rel="" className="bg-secondary hover:opacity-75 transition-opacity rounded-lg grid place-items-center">
        <div className="w-full p-2">
            <div className="flex items-center space-x-2 px-2">
                <div className="min-w-[30px] min-h-[30px]">
                    <img className="w-[30px] h-[30px]" src={`${urlObj.origin}/favicon.ico`} />
                </div>
                <p className="text-wrap text-left">{page.url}</p>
            </div>
        </div>
    </a>
}