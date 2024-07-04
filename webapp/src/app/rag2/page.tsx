import { createRagRequest } from "@/actions/rag";
import Rag2Client from "./client";
import ErrorPage from "@/components/error_page";

export default async function Rag2() {
    try {
        // send a request to start a rag ws
        const req = await createRagRequest()
        console.log(req)

        // create the url to create the ws with
        const url = `${process.env.DB_HOST}/${req.path}`

        return <Rag2Client wsUrl={url} />
    } catch (e) {
        if (e instanceof Error) console.log(e)
        return <ErrorPage msg="" />
    }
}