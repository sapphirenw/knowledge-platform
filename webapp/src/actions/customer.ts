"use server"

import { Customer } from "@/types/customer"
import { cookies } from 'next/headers'
import { sendRequestV1 } from "./api"

export async function getCID(): Promise<string> {
    const cid = cookies().get("cid")?.value
    if (cid == undefined) {
        throw new Error("failed to get the customer id")
    }
    return cid
}

export async function getCustomer(name: string, authToken: string): Promise<Customer> {
    const customer = await sendRequestV1<Customer>({
        route: `beta/customers/get?name=${name}&authToken=${authToken}`,
        ignoreApiKey: true,
    })

    // set cookies to expire after a month
    cookies().set('cid', customer.id, { secure: true, sameSite: "strict", expires: new Date(Date.now() + (30 * 86400 * 1000)) })
    cookies().set('apiKey', authToken, { secure: true, sameSite: "strict", expires: new Date(Date.now() + (30 * 86400 * 1000)) })
    return customer
}