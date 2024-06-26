"use server"

import { Customer } from "@/types/customer"
import { cookies } from 'next/headers'

export async function getCustomer({
    signal,
    name,
}: {
    signal: AbortSignal | undefined,
    name: string,
}): Promise<Customer> {
    let response = await fetch(`${process.env.DB_HOST}/tests/customers/get?name=${name}`, {
        signal: signal,
        method: "GET",
        cache: 'no-store',
    })
    if (response.status != 200) {
        throw new Error(await response.text())
    }

    // set the cookie
    const c = await response.json() as Customer
    cookies().set('cid', c.id, { secure: true })

    return c
}

export async function getCID(): Promise<string> {
    const cid = cookies().get("cid")?.value
    if (cid == undefined) {
        throw new Error("failed to get the customer id")
    }
    return cid
}