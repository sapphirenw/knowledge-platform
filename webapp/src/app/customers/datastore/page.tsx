"use client"

import { Button } from "@/components/ui/button";
import FileUpload from "./fileUpload";
import UserFiles from "./userFiles";
import { vectorizeDatastore } from "@/actions/vector";
import { useState } from "react";
import { toast } from "@/components/ui/use-toast";
import { Loader2 } from "lucide-react";

export default function Files() {
    const [isLoading, setIsLoading] = useState(false)

    const vectorize = async () => {
        setIsLoading(true)
        try {
            await vectorizeDatastore()
            toast({
                title: "Success!",
                description: "Successfully vectorized all documents.",
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
        <Button onClick={vectorize}>
            {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <></>}
            Vectorize Datastore
        </Button>
        <FileUpload />
        <UserFiles />
    </div>

}