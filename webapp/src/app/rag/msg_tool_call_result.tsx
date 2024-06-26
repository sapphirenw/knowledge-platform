import { Button } from "@/components/ui/button";
import { ConversationMessage } from "@/types/conversation"
import { Document } from "@/types/document"
import { File } from "lucide-react";

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
        }

        return items
    }

    return <div className="grid grid-cols-3 gap-4">
        {getItems()}
    </div>
}

function DocumentItem({ doc }: { doc: Document }) {
    return <Button variant="secondary">
        <div className="flex items-center space-x-2 truncate">
            <File />
            <p className="truncate">{doc.filename}</p>
        </div>
    </Button>
}