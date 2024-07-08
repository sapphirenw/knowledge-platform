"use client"

import { handleWebsite } from "@/actions/websites";
import DefaultLoader from "@/components/default_loader";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { toast } from "@/components/ui/use-toast";
import { HandleWebsiteRequest, HandleWebsiteResponse } from "@/types/websites";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useQueryClient } from "@tanstack/react-query";

const FormSchema = z.object({
    domain: z.string(),
    whitelist: z.optional(z.string()),
    blacklist: z.optional(z.string()),
})

export default function WebsiteIngest() {
    const queryClient = useQueryClient()

    const [isLoading, setIsLoading] = useState(false)
    const [scrapeLoading, setScrapeLoading] = useState(false)
    const [data, setData] = useState<HandleWebsiteResponse | undefined>(undefined)

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
    })

    const onSubmit = async (data: z.infer<typeof FormSchema>) => {
        setIsLoading(true)

        // create the blacklist and whitelist
        let wlist: string[] = []
        if (data.whitelist !== undefined && data.whitelist !== "") {
            wlist = data.whitelist!.split(",")
        }

        let blist: string[] = []
        if (data.blacklist !== undefined && data.blacklist !== "") {
            blist = data.blacklist!.split(",")
        }

        const payload: HandleWebsiteRequest = {
            domain: data.domain,
            whitelist: wlist,
            blacklist: blist,
            insert: false,
        }

        // send the request
        try {
            const response = await handleWebsite(payload)
            setData(response)
        } catch (e) {
            toast({
                variant: "destructive",
                title: "Oh no!",
                description: <p>There was an internal issue handling the domain name.</p>
            })
        }

        setIsLoading(false)
    }

    const scrapePages = async () => {
        if (data === undefined) {
            toast({
                variant: "destructive",
                title: "Oh no!",
                description: <p>Your form state is not valid.</p>
            })
            return
        }

        setScrapeLoading(true)

        // send the request
        try {
            await handleWebsite({
                domain: data.site.domain,
                whitelist: data.site.whitelist,
                blacklist: data.site.blacklist,
                insert: true,
            })

            // clear the state
            form.reset()
            setData(undefined)

            // invalidate the site query
            await queryClient.invalidateQueries({ queryKey: ['websites'] })

            // notify user
            toast({
                title: "Success!",
                description: <p>Successfully ingested this website.</p>
            })
        } catch (e) {
            toast({
                variant: "destructive",
                title: "Oh no!",
                description: <p>There was an internal issue scraping the domain.</p>
            })
        }

        setScrapeLoading(false)
    }

    const getSite = () => {
        if (data === undefined) {
            return null
        }

        return <div className="space-y-2">
            <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                    <h4 className="text-lg font-bold">Search Results</h4>
                    <p className="text-sm text-muted-foreground font-medium">- {data.pages.length}</p>
                </div>
                <Button className="space-x-2" onClick={() => scrapePages()}>
                    {scrapeLoading ? <DefaultLoader /> : <></>}
                    <p>Scrape Pages</p>
                </Button>
            </div>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="">Url</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {data.pages.map((item, i) => <TableRow key={`page-${i}`}>
                        <TableCell className="font-medium">{item.url}</TableCell>
                    </TableRow>)}
                </TableBody>
            </Table>
        </div>
    }

    return <div className="w-full space-y-8">
        <div className="w-full">
            <Form {...form}>
                <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
                    <FormField
                        control={form.control}
                        name="domain"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Domain Name</FormLabel>
                                <FormControl>
                                    <Input type="text" placeholder="https://..." {...field} />
                                </FormControl>
                                <FormDescription>
                                    The address of the website you want to ingest.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="whitelist"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Whitelist</FormLabel>
                                <FormControl>
                                    <Input type="text" placeholder="route1/*,route2/*" {...field} />
                                </FormControl>
                                <FormDescription>
                                    A comma separated list of regex patterns that potential search results must match.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="blacklist"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Blacklist</FormLabel>
                                <FormControl>
                                    <Input type="text" placeholder="route1/*,route2/*" {...field} />
                                </FormControl>
                                <FormDescription>
                                    A comma separated list of regex patterns that potential search results must NOT match.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <div className="space-x-2">
                        <Button className="space-x-2" type="submit">
                            {isLoading ? <DefaultLoader /> : <></>}
                            <p>Search</p>
                        </Button>
                    </div>
                </form>
            </Form>
        </div>
        <div className="w-full">
            {getSite()}
        </div>
    </div>
}