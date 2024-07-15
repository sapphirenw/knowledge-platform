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
import { insertSingleWebsitePage } from "@/actions/websites"


const FormSchema = z.object({
    domain: z.string().min(3, {
        message: "Must be a valid url",
    }).url(),
})

export default function InsertSingleWebsitePageButton() {
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
            await insertSingleWebsitePage(data.domain)

            // invalidate query and close
            await queryClient.invalidateQueries({ queryKey: ['websites'] })
            toast({
                title: "Success!",
                description: <p>Successfully inserted the page.</p>
            })
            setOpenDialog(false)
            form.reset()
        } catch (e) {
            if (e instanceof Error) console.log(e)
            toast({
                variant: "destructive",
                title: "On no!",
                description: <p>There was an issue inserting the page.</p>
            })
        }
        setIsLoading(false)
    }

    return <Dialog open={openDialog} onOpenChange={setOpenDialog}>
        <DialogTrigger asChild>
            <Button variant="secondary" onClick={() => setOpenDialog(true)}>Ingest Single Page</Button>
        </DialogTrigger>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>Ingest Single Website Page</DialogTitle>
                <DialogDescription>
                    Opt to insert a single website page into your datastore instead of an entire scraped website. Useful for public profiles, etc.
                </DialogDescription>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                    <FormField
                        control={form.control}
                        name="domain"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Domain</FormLabel>
                                <FormControl>
                                    <Input placeholder="https://..." {...field} />
                                </FormControl>
                                <FormDescription>
                                    The domain name you want to ingest. Must be a valid url.
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