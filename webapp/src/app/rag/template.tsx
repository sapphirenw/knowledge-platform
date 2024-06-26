import Header from "@/components/header";
import Sidebar from "./rag_sidebar";

export default function RagLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return <div className="flex flex-col flex-grow h-screen">
        <Header />
        <div className="flex flex-row flex-grow h-full overflow-hidden">
            <div className="w-[300px]">
                <Sidebar />
            </div>
            <div className="flex-grow flex flex-col overflow-hidden w-full">
                {children}
            </div>
        </div>
    </div>
}