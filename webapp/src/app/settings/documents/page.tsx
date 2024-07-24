import { Label } from "@radix-ui/react-label";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { HydrationBoundary, QueryClient, dehydrate } from "@tanstack/react-query";
import { getWebsites } from "@/actions/websites";
import { listFolder } from "@/actions/document";
import { Separator } from "@/components/ui/separator";
import DocumentsViewClient from "./client";

export default async function DocumentsView() {
    const queryClient = new QueryClient()

    await queryClient.prefetchQuery({
        queryKey: ['documents', null],
        queryFn: listFolder,
    })

    return <div className="space-y-4 w-full">
        {/* <div className="grid place-items-center">
            <p className="text-center text-sm text-muted-foreground max-w-lg">When you add documents or websites, you must manually ensure your vectorized data is in-sync. You can manage and queue requests <Button className="p-0 h-fit" variant="link" asChild><Link href="/settings/datastore/vector-requests">here</Link></Button>.</p>
        </div>
        <Separator className="my-4" /> */}
        <div className="w-full space-y-2">
            <div className="space-y-16">
                <div className="space-y-2 w-full h-max-[500px]">
                    <div className="flex items-center justify-between">
                        <Label htmlFor="user_documents">
                            <h3 className="text-lg font-bold">Documents</h3>
                        </Label>
                        <Button asChild>
                            <Link href="/settings/documents/upload-file">Upload Files</Link>
                        </Button>
                    </div>
                    <div id="user_documents">
                        <HydrationBoundary state={dehydrate(queryClient)}>
                            <DocumentsViewClient />
                        </HydrationBoundary>
                    </div>
                </div>
            </div>
        </div>
    </div>
}