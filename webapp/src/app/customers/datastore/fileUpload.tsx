"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
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
import { toast } from "@/components/ui/use-toast"
import { useState } from "react"
import { uploadDocuments, validateDocuments } from "@/actions/document"
import { useQueryClient } from "@tanstack/react-query"

const FormSchema = z.object({
    files: z.any()
})

export default function FileUpload() {
    const queryClient = useQueryClient()

    const [isLoading, setIsLoading] = useState(false)

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
    })

    async function onSubmit() {
        setIsLoading(true)

        // parse the file array
        var inp = document.getElementById('files') as any;
        if (inp == null) {
            return
        }
        const filelist = inp.files as FileList
        if (filelist.length == 0) {
            setIsLoading(false)
            return
        }

        // create a form to hold the data
        const formData = new FormData();
        for (let i = 0; i < filelist.length; i++) {
            const file = filelist.item(i)
            formData.append(file!.name, file!)
        }

        // validate the files
        var errMessages: React.ReactNode[] = []
        const validationResponse = await validateDocuments(formData)
        for (let i = 0; i < validationResponse.length; i++) {
            if (validationResponse[i].error != undefined) {
                toast({
                    title: `Error with file: ${validationResponse[i].filename}`,
                    description: <p>{validationResponse[i].error}</p>
                })
                errMessages.push(<p>
                    <span className="font-bold">{validationResponse[i].filename}:</span>
                    {validationResponse[i].error}
                </p>)
            }
        }

        if (errMessages.length != 0) {
            toast({
                title: "File Validation Error",
                description: <div className="">
                    {errMessages.map((element, index) => (
                        <div key={index}>
                            {element}
                        </div>
                    ))}
                </div>
            })
            setIsLoading(false)
            return
        }

        // upload the files
        const successfullyUploaded = await uploadDocuments(formData)
        if (!successfullyUploaded) {
            toast({
                title: "Oh no!",
                description: <p>There was an internal issue uploading the files.</p>
            })
        } else {
            // invalidate the query that fetched the files
            await queryClient.invalidateQueries({ queryKey: ['files'] })
            toast({
                title: "Success!",
                description: <p>Successfully uploaded the file(s).</p>
            })
        }

        setIsLoading(false)
    }

    return <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="w-2/3 space-y-6">
            <FormField
                control={form.control}
                name="files"
                render={({ field }) => (
                    <FormItem>
                        <FormLabel>Select Files</FormLabel>
                        <FormControl>
                            <Input id="files" type="file" multiple placeholder="Select files" {...field} />
                        </FormControl>
                        <FormDescription>
                            Select files you want to upload
                        </FormDescription>
                        <FormMessage />
                    </FormItem>
                )}
            />
            <div className="space-x-2">
                <Button type="submit">
                    {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
                    Submit
                </Button>
            </div>
        </form>
    </Form>
}