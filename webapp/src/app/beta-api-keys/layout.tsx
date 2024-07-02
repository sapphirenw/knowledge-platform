import DefaultLayout from "@/components/default_layout";

export default function Layout({ children }: { children: React.ReactNode }) {
    return <DefaultLayout>
        {children}
    </DefaultLayout>
}