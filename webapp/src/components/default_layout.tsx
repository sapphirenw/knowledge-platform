import Footer from "./footer";
import Header from "./header";

export default function DefaultLayout({ children }: { children: React.ReactNode }) {
    return <div className="flex flex-col min-h-screen">
        <Header />
        <div className="flex-grow w-full h-full flex flex-col">
            <div className="flex-grow flex flex-col w-full">
                <div className="grid place-items-center p-12 gap-4 safe-area w-full">
                    {children}
                </div>
            </div>
        </div>
        <Footer />
    </div>
}