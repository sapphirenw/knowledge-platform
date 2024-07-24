import HomeSidebar from "@/components/home_sidebar/sidebar";

export default function Template({ children }: { children: React.ReactNode }) {
    return <div className="flex flex-col flex-grow h-screen">
        <div className="flex flex-row flex-grow h-full overflow-hidden">
            <div className="">
                <HomeSidebar />
            </div>
            <div className="flex-grow flex flex-col overflow-hidden w-full">
                {children}
            </div>
        </div>
    </div>
}