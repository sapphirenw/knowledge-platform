import { getWebsite, getWebsitePages } from "@/actions/websites";
import WebsiteViewClient from "./client";
import ErrorPage from "@/components/error_page";

export default async function WebsiteView({
    params,
}: {
    params: { siteId: string };
}) {
    try {
        const site = await getWebsite(params.siteId)
        const pages = await getWebsitePages(params.siteId)
        return <WebsiteViewClient site={site} pages={pages} />
    } catch {
        return <ErrorPage msg="There was an issue getting the website" />
    }
}