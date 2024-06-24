"use client"

import { Loader2 } from "lucide-react"
import { Button } from "./ui/button"
import { ToastAction } from "@radix-ui/react-toast"
import { createVectorizeRequest } from "@/actions/vector"
import { toast } from "./ui/use-toast"
import { useState } from "react"
import { useQueryClient } from "@tanstack/react-query"

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
        {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
        Vectorize Datastore
    </Button>
}