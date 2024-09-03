"use client"

import { usePathname } from "next/navigation"

export default function HomeLogo() {
    const pathname = usePathname()

    const getIcon = () => {
        if (pathname.includes("/home/resume-builder")) {
            return "/resume-light.svg"
        } else {
            return "/aithing-light.svg"
        }
    }

    return <img
        src={getIcon()}
        className="max-w-full max-h-full object-scale-down"
    />
}