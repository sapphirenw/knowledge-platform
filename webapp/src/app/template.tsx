import Footer from "@/components/footer";
import Header from "@/components/header";
import { cookies } from "next/headers";

export default function Template({ children }: { children: React.ReactNode }) {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")

    return <div className="flex flex-col min-h-screen">
        <Header />
        <div className="flex-grow w-full">{children}</div>
        <Footer />
    </div>
}