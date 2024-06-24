"use client"

import { Button } from "@/components/ui/button";
import FileUpload from "./fileUpload";
import UserFiles from "./userFiles";
import { useState } from "react";
import { toast } from "@/components/ui/use-toast";
import { Loader2 } from "lucide-react";
import { createVectorizeRequest } from "@/actions/vector";
import { ToastAction } from "@/components/ui/toast";

export default function Files() {
    const [isLoading, setIsLoading] = useState(false)

    const vectorize = async () => {
        setIsLoading(true)
        try {
            await createVectorizeRequest()
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

    return <div className="grid place-items-center p-12 gap-4">
        <Button onClick={() => vectorize()}>
            {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
            Vectorize Datastore
        </Button>
        <FileUpload />
        <UserFiles />
    </div>

}