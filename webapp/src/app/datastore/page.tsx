import { Label } from "@radix-ui/react-label";
import VectorRequests from "./vector_requests";
import VectorizationRequest from "@/components/vectorization_request";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { HydrationBoundary, QueryClient, dehydrate } from "@tanstack/react-query";
import { getAllVectorizeRequests } from "@/actions/vector";

export default async function Datastore() {
    const queryClient = new QueryClient()

    await queryClient.prefetchQuery({
        queryKey: ['vectorRequests'],
        queryFn: getAllVectorizeRequests,
    })

    return <div className="grid place-items-center p-12 gap-4">
        <div className="flex items-center space-x-2">
            <VectorizationRequest />
            <Button variant="outline" asChild>
                <Link href="/datastore/files">Files</Link>
            </Button>
            <Button variant="outline" asChild>
                <Link href="/datastore/websites">Websites</Link>
            </Button>
        </div>
        <div className="w-full space-y-2">
            <Label>Vectorization Requests</Label>
            <HydrationBoundary state={dehydrate(queryClient)}>
                <VectorRequests />
            </HydrationBoundary>
        </div>
    </div>
}