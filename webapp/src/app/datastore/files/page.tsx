
import { HydrationBoundary, QueryClient, dehydrate } from "@tanstack/react-query";
import FileUpload from "./file_upload";
import UserFiles from "./user_files";
import VectorizationRequest from "@/components/vectorization_request";
import { listFolder } from "@/actions/document";

export default async function Files() {

    // pre-fetch the files
    const queryClient = new QueryClient()
    await queryClient.prefetchQuery({
        queryKey: ['files'],
        queryFn: listFolder,
    })

    return <div className="grid place-items-center p-12 gap-4 safe-area">
        <VectorizationRequest />
        <FileUpload />
        <HydrationBoundary state={dehydrate(queryClient)}>
            <UserFiles />
        </HydrationBoundary>
    </div>

}