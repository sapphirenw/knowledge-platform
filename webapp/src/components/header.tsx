import { cookies } from "next/headers"

export default function Header() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")

    return <header className="p-4 border-b border-b-border">
        <div className="safe-area">
            <div className="flex items-center justify-between">
                <div className=""><p>customerId: {customerId?.value ?? "undefied"}</p></div>
            </div>
        </div>
    </header>
}