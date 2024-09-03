import dynamicIconImports from "lucide-react/dynamicIconImports"

export type Integration = {
    title: string
    description: string
    href: string
    icon: keyof typeof dynamicIconImports // lucide dynamic loaded string
}