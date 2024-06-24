import Sidebar from "./rag_sidebar";

export default function RagLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return <div className="flex flex-col flex-grow h-full">
        <div className="flex flex-row flex-grow h-full">
            <Sidebar />
            <div className="flex-grow flex flex-col overflow-hidden">
                {children}
            </div>
        </div>
    </div>
}