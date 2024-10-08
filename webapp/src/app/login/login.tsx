"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { getCustomer } from '@/actions/customer'
import { Button } from '@/components/ui/button'
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
import { AuthError } from "@/types/errors"
import { navigate } from "@/actions/redirect"

const FormSchema = z.object({
    name: z.string().min(2, {
        message: "Username must be at least 2 characters.",
    }),
    authToken: z.string().min(1, {
        message: "Auth token cannot be empty",
    }),
})

export default function LoginClient({ cid }: { cid?: string }) {
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
            await getCustomer(data.name, data.authToken)
            toast({
                title: "Successfully autheticated!",
            })
            // redirect to settings
            navigate(`/settings`)
        } catch (e) {
            if (e instanceof AuthError) {
                toast({
                    variant: "destructive",
                    title: "Not Allowed.",
                    description: <p>Please check your auth token and try again</p>
                })
            } else {
                toast({
                    variant: "destructive",
                    title: "There was an error",
                    description: <p>{e instanceof Error ? e.message : "Unknown error"}</p>
                })
            }

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
                                Enter the details given to you to log into your session.
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
                                            <Input type="name" placeholder="(My name)" {...field} />
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
                                            <Input type="text" placeholder="* * * * *" {...field} />
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
            <p className="text-sm opacity-50">Current ID: {cid ?? "undefined"}</p>
            <Button variant="outline" onClick={() => removeData()}>Remove ID</Button>
        </div>
    </div>
}