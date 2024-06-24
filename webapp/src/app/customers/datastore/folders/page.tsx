
import FileUpload from "./fileUpload";
import UserFiles from "./userFiles";
import VectorizationRequest from "@/components/vectorizationRequest";

export default function Files() {
    return <div className="grid place-items-center p-12 gap-4">
        <VectorizationRequest />
        <FileUpload />
        <UserFiles />
    </div>

}