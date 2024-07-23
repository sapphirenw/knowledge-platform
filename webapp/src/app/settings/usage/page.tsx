import { getUsage, getUsageGrouped } from "@/actions/usage"
import { UsageGroupedRecord, UsageResponse } from "@/types/usage"
import CustomerUsageClient from "./client"

export default async function CustomerUsage({
    // searchParams,
}: {
        // searchParams: { [key: string]: string | undefined }
    }) {

    let usage: UsageResponse | undefined = undefined
    let usageGrouped: UsageGroupedRecord[] | undefined = undefined

    try {
        usage = await getUsage({ page: 1 })
        usageGrouped = await getUsageGrouped()
    } catch (e) {
        if (e instanceof Error) console.error(e)
    }

    return <CustomerUsageClient initialData={usage} usageGrouped={usageGrouped} />
}