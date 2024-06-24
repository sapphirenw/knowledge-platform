import { Label } from "@radix-ui/react-label";
import VectorRequests from "./vectorRequests";

export default function Datastore() {
    return <div className="grid place-items-center p-12 gap-4">
        <div className="w-full space-y-2">
            <Label>Vectorization Requests</Label>
            <VectorRequests />
        </div>
    </div>
}