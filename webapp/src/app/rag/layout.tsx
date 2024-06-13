import Sidebar from "./rag_sidebar";

export default function RagLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return <div className="flex">
        <Sidebar />
        <div className="w-full">{children}</div>
    </div>
}