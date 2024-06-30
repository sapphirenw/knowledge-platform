import { Label } from "@radix-ui/react-label";
import VectorRequests from "./vector-requests/page";
import VectorizationRequest from "@/components/vectorization_request";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { HydrationBoundary, QueryClient, dehydrate } from "@tanstack/react-query";
import { getAllVectorizeRequests } from "@/actions/vector";
import { getWebsites } from "@/actions/websites";
import { listFolder } from "@/actions/document";
import UserFiles from "./user_files";
import UserWebsites from "./user_websites";
import { Separator } from "@/components/ui/separator";

export default async function Datastore() {
    const queryClient = new QueryClient()

    // await queryClient.prefetchQuery({
    //     queryKey: ['vectorRequests'],
    //     queryFn: getAllVectorizeRequests,
    // })

    await queryClient.prefetchQuery({
        queryKey: ['files'],
        queryFn: listFolder,
    })
    await queryClient.prefetchQuery({
        queryKey: ['websites'],
        queryFn: getWebsites,
    })

    return <div className="space-y-4 w-full">
        <div className="grid place-items-center">
            <Button asChild>
                <Link href="/datastore/vector-requests">My Vector Requests</Link>
            </Button>
        </div>
        <Separator className="my-4" />
        <div className="w-full space-y-2">
            <HydrationBoundary state={dehydrate(queryClient)}>
                <div className="space-y-16">
                    <div className="space-y-2 w-full h-max-[500px]">
                        <div className="flex items-center justify-between">
                            <Label htmlFor="user_files">
                                <h3 className="text-lg font-bold">Files</h3>
                            </Label>
                            <Button asChild>
                                <Link href="/datastore/upload-file">Upload Files</Link>
                            </Button>
                        </div>
                        <div id="user_files">
                            <UserFiles />
                        </div>
                    </div>
                    <div className="space-y-2 w-full h-max-[500px]">
                        <div className="flex items-center justify-between">
                            <Label htmlFor="user_websites">
                                <h3 className="text-lg font-bold">Websites</h3>
                            </Label>
                            <Button asChild>
                                <Link href="/datastore/ingest-site">Ingest New Site</Link>
                            </Button>
                        </div>
                        <div id="user_websites">
                            <UserWebsites />
                        </div>
                    </div>
                </div>
            </HydrationBoundary>
        </div>
    </div>
}