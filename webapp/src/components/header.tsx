import { cookies } from "next/headers"
import { Button } from "./ui/button"
import Link from "next/link"

export default function Header() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")

    return <header className="p-4 border-b border-b-border">
        <div className="safe-area">
            <div className="flex items-center justify-between">
                <div className=""></div>
                <div className="flex items-center space-x-2">
                    {customerId === undefined ? null : <div className="flex items-center space-x-2">
                        <Button variant="outline" asChild>
                            <a href="/rag">Chat</a>
                        </Button>
                        <Button variant="outline" asChild>
                            <a href="/models">Models</a>
                        </Button>
                        <Button variant="outline" asChild>
                            <a href="/datastore">Datastore</a>
                        </Button>
                        <Button variant="outline" asChild>
                            <a href="/usage">Usage</a>
                        </Button>
                    </div>}
                    <Button variant="outline" asChild>
                        <a href="/login">Login</a>
                    </Button>
                </div>
            </div>
        </div>
    </header>
}