"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { getCustomer } from '@/actions/customer'
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
import Cookies from "js-cookie"
import DefaultLoader from "@/components/default_loader"

const FormSchema = z.object({
    name: z.string().min(2, {
        message: "Username must be at least 2 characters.",
    }),
})

export default function Login() {
    const [isLoading, setIsLoading] = useState(false)

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            name: "",
        },
    })

    async function onSubmit(data: z.infer<typeof FormSchema>) {
        // send the request
        setIsLoading(true)
        try {
            const customer = await getCustomer({
                signal: undefined,
                name: data.name,
            })
            toast({
                title: "Successfully got your customerId",
                description: <p>id: {customer.id}</p>
            })
        } catch (e) {

            toast({
                variant: "destructive",
                title: "There was an error",
                description: <p>{e instanceof Error ? e.message : "Unknown error"}</p>
            })
        }
        setIsLoading(false)
    }

    const removeData = () => {
        Cookies.remove('cid')
        location.reload()
    }

    return <div className="safe-area p-16 grid place-items-center gap-4">
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Name</FormLabel>
                            <FormControl>
                                <Input type="name" placeholder="jake" {...field} />
                            </FormControl>
                            <FormDescription>
                                For now, this will map to your temporary ID, and can be used to fetch your customerId again.
                            </FormDescription>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <div className="space-x-2">
                    <Button type="submit">
                        {isLoading ? <DefaultLoader /> : <></>}
                        Submit
                    </Button>
                </div>
            </form>
        </Form>
        <div className="text-left w-full">
            <Button variant="outline" onClick={() => removeData()}>Remove ID</Button>
        </div>
    </div>
}