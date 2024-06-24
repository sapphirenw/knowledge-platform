import Footer from "./footer";
import Header from "./header";

export default function DefaultTemplate({ children }: { children: React.ReactNode }) {
    return <div className="flex flex-col min-h-screen">
        <Header />
        <div className="flex-grow w-full h-full flex flex-col">
            <div className="flex-grow flex flex-col">{children}</div>
        </div>
        <Footer />
    </div>
}