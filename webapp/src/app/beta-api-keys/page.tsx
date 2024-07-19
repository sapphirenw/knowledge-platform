"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Button } from '@/components/ui/button'
import {
    AlertDialog,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
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
import { createApiKey } from "@/actions/beta_auth"
import { BetaApiKey } from "@/types/beta_auth"
import { Switch } from "@/components/ui/switch"

const FormSchema = z.object({
    name: z.string().min(2, {
        message: "Username must be at least 2 characters.",
    }),
    isAdmin: z.boolean(),
    authToken: z.string().min(1, {
        message: "Auth token cannot be empty",
    }),
})

export default function CreateBetaApiKey() {
    const [isLoading, setIsLoading] = useState(false)
    const [openDialog, setOpenDialog] = useState(false)
    const [data, setData] = useState<BetaApiKey | undefined>(undefined)

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            name: "",
            isAdmin: false,
        },
    })

    async function onSubmit(data: z.infer<typeof FormSchema>) {
        // send the request
        setIsLoading(true)
        try {
            const response = await createApiKey(data.name, data.isAdmin, data.authToken)
            setData(response)
            setOpenDialog(true)
            form.reset()
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

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text).then(() => {
            toast({
                title: "Success!",
                description: <p>Successfully copied api key to clipboard.</p>
            })
        }).catch(err => {
            toast({
                variant: "destructive",
                title: `Failed to copy: ${err}`,
            })
        });
    };

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
                            <CardTitle className="text-2xl">Create API Key</CardTitle>
                            <CardDescription>
                                If you are an admin, you can create a new API Key for a client here.
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
                                            This name will be tied to the API Key, so this is the name you will give them.
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="isAdmin"
                                render={({ field }) => (
                                    <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                                        <div className="space-y-0.5">
                                            <FormLabel className="text-base">Is Admin</FormLabel>
                                            <FormDescription>
                                                Grants the ability to view admin screens like debug chats.
                                            </FormDescription>
                                        </div>
                                        <FormControl>
                                            <Switch
                                                checked={field.value}
                                                onCheckedChange={field.onChange}
                                            />
                                        </FormControl>
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="authToken"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Auth Token</FormLabel>
                                        <FormControl>
                                            <Input type="text" placeholder="* * * * *" {...field} />
                                        </FormControl>
                                        <FormDescription>
                                            The auth token that was given to you specifically to create API Keys.
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
        <AlertDialog open={openDialog} onOpenChange={setOpenDialog}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Successfully Created Api Key</AlertDialogTitle>
                    <AlertDialogDescription>
                        <div className="space-y-4">
                            <p>CAUTION: You will only be able to view this information once.</p>
                            <div className="spave-y-2">
                                <p>Name: <code>{data?.name ?? "seesf"}</code></p>
                                <p>Appi Key: <code>{data?.id ?? "fsefe"}</code></p>
                            </div>
                        </div>
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Close</AlertDialogCancel>
                    <Button onClick={() => copyToClipboard(data?.id ?? "")}>Copy Key</Button>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    </div>
}