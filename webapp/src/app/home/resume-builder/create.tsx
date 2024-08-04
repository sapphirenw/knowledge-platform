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

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"


import { Input } from "@/components/ui/input"
import { useQueryClient } from "@tanstack/react-query"
import DefaultLoader from "@/components/default_loader"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { toast } from "@/components/ui/use-toast"
import { createResumeApplication } from "@/actions/resume"
import { Textarea } from "@/components/ui/textarea"


const FormSchema = z.object({
    title: z.string().min(3, {
        message: "Must be at least 3 characters",
    }),
    link: z.string().url(),
    companySite: z.string(),
    rawText: z.string().min(1, {
        message: "The raw text cannot be empty"
    })
})

export default function CreateResumeApplicationButton() {
    const queryClient = useQueryClient()

    const [openDialog, setOpenDialog] = useState(false)
    const [isLoading, setIsLoading] = useState(false)

    // fetch the available models

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
    })

    async function onSubmit(data: z.infer<typeof FormSchema>) {
        // send the request
        setIsLoading(true)
        try {
            await createResumeApplication({
                title: data.title,
                link: data.link,
                companySite: data.companySite,
                rawText: data.rawText,
            })

            // invalidate query and close
            await queryClient.invalidateQueries({ queryKey: ['resumeApplications'] })
            toast({
                title: "Success!",
                description: <p>Successfully created the application.</p>
            })
            setOpenDialog(false)
            form.reset()
        } catch (e) {
            if (e instanceof Error) console.log(e)
            toast({
                variant: "destructive",
                title: "On no!",
                description: <p>There was an issue creating the application.</p>
            })
        }
        setIsLoading(false)
    }

    return <Dialog open={openDialog} onOpenChange={setOpenDialog}>
        <DialogTrigger asChild>
            <Button onClick={() => setOpenDialog(true)}>Create New</Button>
        </DialogTrigger>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>Create Application</DialogTitle>
                <DialogDescription>
                    An application allows you to create tailored resumes for this job posting using AI against your own data.
                </DialogDescription>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                    <FormField
                        control={form.control}
                        name="title"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Title</FormLabel>
                                <FormControl>
                                    <Input placeholder="Default" {...field} />
                                </FormControl>
                                <FormDescription>
                                    The name of this application / job posting.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="link"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>
                                    <p>Link <span className="text-sm opacity-50">(optional)</span></p>
                                </FormLabel>
                                <FormControl>
                                    <Input placeholder="https://..." {...field} />
                                </FormControl>
                                <FormDescription>
                                    A link to the job posting.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="companySite"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>
                                    <p>Company Website <span className="text-sm opacity-50">(optional)</span></p>
                                </FormLabel>
                                <FormControl>
                                    <Input placeholder="https://..." {...field} />
                                </FormControl>
                                <FormDescription>
                                    An optional reference to the company website.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="rawText"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Job Posting Text</FormLabel>
                                <FormControl>
                                    <Textarea
                                        placeholder="We want you to have 5 years experience in ..."
                                        className="resize-none"
                                        {...field}
                                    />
                                </FormControl>
                                <FormDescription>
                                    The raw text of the job posting.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <DialogFooter>
                        <Button type="submit">
                            <div className="flex space-x-2 items-center">
                                {isLoading ? <DefaultLoader /> : <></>}
                                <p>Submit</p>
                            </div>
                        </Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
}