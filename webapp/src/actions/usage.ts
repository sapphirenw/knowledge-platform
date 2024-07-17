import { UsageGroupedRecord, UsageResponse } from "@/types/usage";
import { getCID } from "./customer";
import { sendRequestV1 } from "./api";

export async function getUsage({
    page
}: {
    page: number,
}): Promise<UsageResponse> {
    const cid = await getCID()
    return await sendRequestV1<UsageResponse>({
        route: `customers/${cid}/usage?page=${page}&batchSize=10`,
        method: "GET"
    })
}

export async function getUsageGrouped(): Promise<UsageGroupedRecord[]> {
    const cid = await getCID()
    return await sendRequestV1<UsageGroupedRecord[]>({
        route: `customers/${cid}/usageGrouped`,
        method: "GET"
    })
}