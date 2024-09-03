import { Integration } from "@/types/integrations"

export async function getUserIntegrations(): Promise<Integration[]> {
    const integrations: Integration[] = []

    // basic chat that all users have access to
    integrations.push({
        title: "Chat",
        description: "A simple chat against your stored information with a wide-selection of models.",
        href: "/home/chat",
        icon: "message-square",
    })

    // TODO -- static integrations should be sourced from the user's preferences

    integrations.push({
        title: "Resume Builder",
        description: "Use Resume by AIThing to create custom resumes specific to job postings.",
        href: "/home/resume-builder",
        icon: "file-text",
    })

    return integrations
}