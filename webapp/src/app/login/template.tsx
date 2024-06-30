import DefaultTemplate from "@/components/default_template";

export default function Template({ children }: { children: React.ReactNode }) {
    return <DefaultTemplate>
        <div className="safe-area p-16 grid place-items-center gap-4 w-full">
            {children}
        </div>
    </DefaultTemplate>
}