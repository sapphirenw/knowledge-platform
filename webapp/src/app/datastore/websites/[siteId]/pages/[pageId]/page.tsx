import { getWebsite, getWebsitePageContent, getWebsitePages } from "@/actions/websites";
import ErrorPage from "@/components/error_page";
import WebsitePageViewClient from "./client";

export default async function WebsitePageView({
    params,
}: {
    params: { siteId: string, pageId: string, };
}) {
    try {
        const pageContentResponse = await getWebsitePageContent(params.siteId, params.pageId)
        return <WebsitePageViewClient resp={pageContentResponse} />
    } catch {
        return <ErrorPage msg="There was an issue getting the website page content" />
    }
}