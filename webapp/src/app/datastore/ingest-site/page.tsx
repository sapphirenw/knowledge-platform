"use client"

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
import { insertWebsite, searchWebsite } from "@/actions/websites";
import { Checkbox } from "@/components/ui/checkbox";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";

const FormSchema = z.object({
    domain: z.string(),
    whitelist: z.optional(z.string()),
    blacklist: z.optional(z.string()),
    useSitemap: z.boolean(),
    allowOtherDomains: z.boolean(),
})

export default function WebsiteIngest() {
    const queryClient = useQueryClient()

    const [isLoading, setIsLoading] = useState(false)
    const [scrapeLoading, setScrapeLoading] = useState(false)
    const [data, setData] = useState<HandleWebsiteResponse | undefined>(undefined)
    const [selectedPages, setSelectedPages] = useState<string[]>([])

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            useSitemap: true,
            allowOtherDomains: false,
        }
    })

    const setDefaultBlacklist = (e: any) => {
        e.preventDefault()
        form.setValue("blacklist", "apple.com,google.com,wikipeida.com,x.com,twitter.com,github.com,gitlab.com")
    }

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
            useSitemap: data.useSitemap,
            allowOtherDomains: data.allowOtherDomains,
        }

        // send the request
        try {
            const response = await searchWebsite(payload)
            setData(response)
            setSelectedPages(response.pages.map((item, i) => item.url))
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
            await insertWebsite({
                domain: data.site.domain,
                whitelist: data.site.whitelist,
                blacklist: data.site.blacklist,
                pages: selectedPages,
            })

            // clear the state
            form.reset()
            setData(undefined)
            setSelectedPages([])

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
                    <p className="text-sm text-muted-foreground font-medium">- {selectedPages.length} Selected</p>
                </div>
                <Button className="space-x-2" onClick={() => scrapePages()}>
                    {scrapeLoading ? <DefaultLoader /> : <></>}
                    <p>Scrape Pages</p>
                </Button>
            </div>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="">Include</TableHead>
                        <TableHead className="">Url</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {data.pages.map((item, i) => <TableRow key={`page-${i}`}>
                        <TableCell className="w-fit">
                            <Switch

                                checked={selectedPages.indexOf(item.url) !== -1}
                                onCheckedChange={(e) => {
                                    if (e === true) {
                                        setSelectedPages((prev) => prev.concat(item.url))
                                    } else {
                                        setSelectedPages((prev) => prev.filter((val, i) => val !== item.url))
                                    }
                                }}
                            />
                        </TableCell>
                        <TableCell className="font-medium break-all">{item.url}</TableCell>
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
                                    {"The address of the website you want to ingest. The default scheme used will be \"https\"."}
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
                                <FormLabel>
                                    <div className="flex items-center space-x-4">
                                        <p>Blacklist</p>
                                        <button className="text-primary hover:opacity-50" onClick={(e) => setDefaultBlacklist(e)}>Set recommended defaults</button>
                                    </div>
                                </FormLabel>
                                <FormControl>
                                    <Input type="text" placeholder="route1/*,route2/*,domain3.com,..." {...field} />
                                </FormControl>
                                <FormDescription>
                                    A comma separated list of regex patterns that potential search results must NOT match.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="useSitemap"
                        render={({ field }) => (
                            <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                                <div className="space-y-0.5">
                                    <FormLabel className="text-base">
                                        Use the sitemap to parse the website pages.
                                    </FormLabel>
                                    <FormDescription>
                                        Using the sitemap to scrape a website is recommended and will return results faster. If required, disabling this option uses a more traditional webscraping approach.
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
                    {form.getValues().useSitemap ? null : <div>
                        <FormField
                            control={form.control}
                            name="allowOtherDomains"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                                    <div className="space-y-0.5">
                                        <FormLabel className="text-base">
                                            Allow other domains in scrape.
                                        </FormLabel>
                                        <FormDescription>
                                            {"Allow the scraper to visit websites that are not from the same domain as specified above. When setting this option, it is recommended to block some domains in the \"Blacklist\" section to clean the results."}
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
                    </div>}
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