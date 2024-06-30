import DefaultTemplate from "@/components/default_template";

export default function Template({ children }: { children: React.ReactNode }) {
    return <DefaultTemplate>
        <div className="grid place-items-center p-12 gap-4 safe-area">
            {children}
        </div>
    </DefaultTemplate>
}