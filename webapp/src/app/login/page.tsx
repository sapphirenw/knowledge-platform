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
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"

const FormSchema = z.object({
    name: z.string().min(2, {
        message: "Username must be at least 2 characters.",
    }),
    authToken: z.string().min(1, {
        message: "Auth token cannot be empty",
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

    return <div className="w-full space-y-4">
        <div className="w-full mx-auto">
            <Card className="w-full max-w-md mx-auto">
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <CardHeader>
                            <CardTitle className="text-2xl">Login</CardTitle>
                            <CardDescription>
                                Enter your email below to login to your account.
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="grid gap-4">
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
                            <FormField
                                control={form.control}
                                name="authToken"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Beta Auth Token</FormLabel>
                                        <FormControl>
                                            <Input type="password" placeholder="*****" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            Auth token that you have been supplied.
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </CardContent>
                        <CardFooter>
                            <div className="space-x-2">
                                <Button type="submit">
                                    {isLoading ? <DefaultLoader /> : <></>}
                                    Submit
                                </Button>
                            </div>
                        </CardFooter>
                    </form>
                </Form>
            </Card>
        </div>
        <div className="text-center w-full space-y-2">
            <p className="text-sm opacity-50">Current ID: {Cookies.get("cid") ?? "undefied"}</p>
            <Button variant="outline" onClick={() => removeData()}>Remove ID</Button>
        </div>
    </div>
}