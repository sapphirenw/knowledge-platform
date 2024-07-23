import SettingsSidebar from "@/components/settings_sidebar/sidebar";

export default function Template({ children }: { children: React.ReactNode }) {
    return <div className="flex flex-col flex-grow h-screen">
        {/* <Header /> */}
        <div className="flex flex-row flex-grow h-full overflow-hidden">
            <div className="">
                <SettingsSidebar />
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