import DefaultLayout from "@/components/default_layout";

export default function Template({ children }: { children: React.ReactNode }) {
    return <DefaultTemplate>
        {children}
    </DefaultTemplate>
}