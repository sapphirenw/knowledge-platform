import { User, BriefcaseBusiness, Palette, GraduationCap, Lightbulb, FileText } from "lucide-react"

import Sidebar from "@/components/sidebar/sidebar";

export default function Template({ children }: { children: React.ReactNode }) {
    return <div className="flex flex-col flex-grow h-screen">
        {/* <Sidebar /> */}
        <div className="flex flex-row flex-grow h-full overflow-hidden">
            <div className="">
                <Sidebar props={{
                    groups: [
                        {
                            title: "Personal Information",
                            titleSm: "Info",
                            items: [
                                {
                                    href: "/home/resume-builder/information",
                                    title: "My Content",
                                    icon: <FileText size={16} />,
                                },
                                {
                                    href: "/home/resume-builder/about",
                                    title: "About Me",
                                    icon: <User size={16} />,
                                },
                                {
                                    href: "/home/resume-builder/experience",
                                    title: "Experience",
                                    icon: <BriefcaseBusiness size={16} />,
                                },
                                {
                                    href: "/home/resume-builder/projects",
                                    title: "Projects",
                                    icon: <Palette size={16} />,
                                },
                                {
                                    href: "/home/resume-builder/education",
                                    title: "Education",
                                    icon: <GraduationCap size={16} />,
                                },
                                {
                                    href: "/home/resume-builder/skills",
                                    title: "Skills",
                                    icon: <Lightbulb size={16} />,
                                },
                            ],
                        },
                    ],
                    footer: <p></p>,
                    allowsClose: true,
                }} />
            </div>
            <div className="flex-grow w-full h-full flex flex-col overflow-scroll">
                <div className="flex-grow flex flex-col w-full">
                    <div className="grid place-items-center p-12 gap-4 safe-area w-full">
                        <div className="overflow-scroll w-full">
                            <div className="p-4">
                                <div className="safe-area">
                                    {children}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
}