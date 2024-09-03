import Sidebar from "@/components/sidebar/sidebar";
import { AreaChart, BetweenHorizontalStart, Brain, Earth, FileText, Home, Package } from "lucide-react";

export default function Template({ children }: { children: React.ReactNode }) {
    return <div className="flex flex-col flex-grow h-screen">
        {/* <Header /> */}
        <div className="flex flex-row flex-grow h-full overflow-hidden">
            <div className="">
                <Sidebar props={{
                    header: {
                        href: "/home",
                        title: "Home",
                        icon: <Home size={16} />,
                    },
                    groups: [
                        {
                            title: "Datastore",
                            titleSm: "Data",
                            items: [
                                {
                                    href: "/settings/vector-requests",
                                    title: "Vector Requests",
                                    icon: <BetweenHorizontalStart size={16} />,
                                },
                                {
                                    href: "/settings/documents",
                                    title: "Documents",
                                    icon: <FileText size={16} />
                                },
                                {
                                    href: "/settings/websites",
                                    title: "Websites",
                                    icon: <Earth size={16} />
                                },
                            ],
                        },
                        {
                            title: "General",
                            titleSm: "Gen",
                            items: [
                                {
                                    href: "/settings/available-llms",
                                    title: "Available LLMs",
                                    icon: <Package size={16} />
                                },
                                {
                                    href: "/settings/custom-models",
                                    title: "Custom Models",
                                    icon: <Brain size={16} />
                                },
                                {
                                    href: "/settings/usage",
                                    title: "Usage",
                                    icon: <AreaChart size={16} />
                                },
                            ],
                        },
                    ],
                    footer: <p></p>,
                    allowsClose: false,
                }} />
            </div>
            <div className="flex-grow w-full h-full flex flex-col overflow-scroll">
                <div className="flex-grow flex flex-col w-full">
                    <div className="grid place-items-center p-12 gap-4 safe-area w-full">
                        {children}
                    </div>
                </div>
            </div>
        </div>
    </div>
}