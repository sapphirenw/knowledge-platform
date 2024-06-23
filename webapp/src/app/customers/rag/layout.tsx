import Sidebar from "./rag_sidebar";

export default function RagLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return <div className="flex h-screen">
        <Sidebar />
        <div className="flex-grow h-full overflow-hidden">
            {children}
        </div>
    </div>
}