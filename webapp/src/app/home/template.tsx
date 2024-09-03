import SidebarRowView from "@/components/sidebar/sidebar_row";
import { Settings } from "lucide-react";
import HomeBreadcrumb from "./breadcrumb";
import HomeLogo from "./logo";

export default function Template({ children }: { children: React.ReactNode }) {


    return <div className="flex flex-col flex-grow h-screen">
        <div className="border-b">
            <div className="flex items-center safe-area h-[50px] justify-between">
                <div className="flex items-center">
                    <div className="h-[40px] w-[175px]"><HomeLogo /></div>
                    <HomeBreadcrumb />
                </div>
                <div className="w-fit">
                    <SidebarRowView
                        item={{
                            href: "/settings",
                            title: "Settings",
                            icon: <Settings size={16} />
                        }}
                        isOpen={false} variant="outline"
                    />
                </div>
            </div>
        </div>
        {/* <div className="border-b border-border w-full">
            <div className="safe-area p-2 flex items-center space-x-2">
                <div className="pr-4">
                    <img
                        src="/aithing-light.svg"
                        height="22"
                        width="75"
                    />
                </div>
                <div className="flex items-center justify-between w-full">
                    <div className="flex items-center space-x-2">
                        <SidebarRowView
                            item={{
                                href: "/home/chat",
                                title: "Chat",
                                icon: <MessageSquareText size={16} />
                            }}
                            isOpen={false} variant="outline"
                        />
                        <SidebarRowView
                            item={{
                                href: "/home/resume-builder",
                                title: "Resume Builder",
                                icon: <FileSpreadsheet size={16} />
                            }}
                            isOpen={false} variant="outline"
                        />
                    </div>
                    <div className="w-fit">
                        <SidebarRowView
                            item={{
                                href: "/settings",
                                title: "Settings",
                                icon: <Settings size={16} />
                            }}
                            isOpen={false} variant="outline"
                        />
                    </div>
                </div>
            </div>
        </div> */}
        <div className="flex-grow flex flex-col overflow-hidden w-full">
            {children}
        </div>
    </div>
}