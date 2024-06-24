"use client"

import { Label } from "@radix-ui/react-label";
import VectorRequests from "./vectorRequests";
import VectorizationRequest from "@/components/vectorizationRequest";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Datastore() {
    return <div className="grid place-items-center p-12 gap-4">
        <div className="flex items-center space-x-2">
            <VectorizationRequest />
            <Button variant="outline" asChild>
                <Link href="/datastore/folders">Folders</Link>
            </Button>
            <Button variant="outline" asChild>
                <Link href="/datastore/websites">Websites</Link>
            </Button>
        </div>
        <div className="w-full space-y-2">
            <Label>Vectorization Requests</Label>
            <VectorRequests />
        </div>
    </div>
}