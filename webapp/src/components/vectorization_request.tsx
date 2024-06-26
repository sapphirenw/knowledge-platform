"use client"

import { Loader2 } from "lucide-react"
import { Button } from "./ui/button"
import { ToastAction } from "@radix-ui/react-toast"
import { createVectorizeRequest } from "@/actions/vector"
import { toast } from "./ui/use-toast"
import { useState } from "react"
import { useQueryClient } from "@tanstack/react-query"
import DefaultLoader from "./default_loader"

export default function VectorizationRequest() {
    const [isLoading, setIsLoading] = useState(false)
    const queryClient = useQueryClient()

    const vectorize = async () => {
        setIsLoading(true)
        try {
            await createVectorizeRequest()
            // invalidate the query to get vector requests
            await queryClient.invalidateQueries({ queryKey: ['vectorRequests'] })
            toast({
                title: "Success!",
                description: "Successfully initiated a request to vectorize your datastore.",
                action: <ToastAction onClick={() => console.log("Hello")} altText="View Status">View Status</ToastAction>,
            })
        } catch (e) {
            toast({
                variant: "destructive",
                title: "Uh oh! Something went wrong.",
                description: "There was a problem with your request.",
            })
        }
        setIsLoading(false)
    }

    return <Button onClick={() => vectorize()}>
        {isLoading ? <DefaultLoader /> : <></>}
        Vectorize Datastore
    </Button>
}