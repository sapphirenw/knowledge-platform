"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"

import { Input } from "@/components/ui/input"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { useState } from "react"
import { toast } from "@/components/ui/use-toast"
import { Textarea } from "@/components/ui/textarea"
import { getResumeAbout } from "@/actions/resume"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"


const FormSchema = z.object({
    name: z.string().min(1, {
        message: "Cannot be empty",
    }),
    email: z.string().min(1, {
        message: "Cannot be empty",
    }).email({ message: "Invalid email address" }),
    phone: z.string(),
    title: z.string().min(1, { message: "Cannot be empty" }),
    location: z.string().min(1, { message: "Cannot be empty" }),

    github: z.string().url().startsWith("https://", { message: "Must be a valid url" }).includes("git"),
    linkedin: z.string().url().includes("https://github.com", {
        message: "Must be a valid LinkedIn link"
    }),
})

export default function ResumeAboutClient() {
    const queryClient = useQueryClient()

    const { data, status } = useQuery({
        queryKey: ["resumeAbout"],
        queryFn: () => getResumeAbout(),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        return <ErrorPage msg="" />
    }

    const [isLoading, setIsLoading] = useState(false)

    // fetch the available models

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            name: data.name,
            email: data.email,
            phone: data.phone,
            title: data.title,
            location: data.location,
            github: data.github,
            linkedin: data.linkedin,
        }
    })

    async function onSubmit(data: z.infer<typeof FormSchema>) {
        // send the request
        setIsLoading(true)
        try {

            toast({
                title: "Success!",
                description: <p>Successfully updated your preferenes.</p>
            })
            form.reset()
        } catch (e) {
            if (e instanceof Error) console.log(e)
            toast({
                variant: "destructive",
                title: "On no!",
                description: <p>There was an issue updating your preferences.</p>
            })
        }
        setIsLoading(false)
    }

    return <div className="max-w-lg mx-auto">
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <div className="flex items-center space-x-2 w-full">
                    <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                            <FormItem className="w-full">
                                <FormLabel>
                                    <p>Name <span className="text-xs text-red-500"> *</span></p>
                                </FormLabel>
                                <FormControl>
                                    <Input placeholder="Jake" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="phone"
                        render={({ field }) => (
                            <FormItem className="w-full">
                                <FormLabel>
                                    <p>Phone <span className="opacity-50 text-xs">(Optional)</span></p>
                                </FormLabel>
                                <FormControl>
                                    <Input placeholder="Phone number" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                </div>
                <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                <p>Email <span className="text-xs text-red-500"> *</span></p>
                            </FormLabel>
                            <FormControl>
                                <Input type="email" placeholder="me@email.com" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="title"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                <p>Job Title <span className="text-xs text-red-500"> *</span></p>
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Software Engineer" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="location"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>
                                <p>Location <span className="text-xs text-red-500"> *</span></p>
                            </FormLabel>
                            <FormControl>
                                <Input placeholder="Paris, FR" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <div className="space-y-2">
                    <FormField
                        control={form.control}
                        name="github"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>
                                    <p>Links <span className="opacity-50 text-xs">(Optional)</span></p>
                                </FormLabel>
                                <FormControl>
                                    <Input placeholder="https://github.com/..." {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="linkedin"
                        render={({ field }) => (
                            <FormItem>
                                <FormControl>
                                    <Input placeholder="https://linkedin.com/..." {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                </div>
            </form>
        </Form>
    </div>
}