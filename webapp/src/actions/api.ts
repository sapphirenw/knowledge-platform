"use server"

import { AuthError } from "@/types/errors"
import { cookies } from "next/headers"
import { RedirectType, redirect } from "next/navigation"

// wrapper function to handle things such as auth
// Overload for when skipJsonParse is true
export async function sendRequestV1({
    route,
    method,
    cache,
    headers,
    body,
    ignoreApiKey,
    debugPrint,
    skipJsonParse,
}: {
    route: string,
    method?: string,
    cache?: RequestCache,
    headers?: Headers,
    body?: BodyInit,
    ignoreApiKey?: boolean,
    debugPrint?: boolean,
    skipJsonParse: true,
}): Promise<undefined>;

// Overload for when skipJsonParse is false or not provided
export async function sendRequestV1<T>({
    route,
    method,
    cache,
    headers,
    body,
    ignoreApiKey,
    debugPrint,
    skipJsonParse,
}: {
    route: string,
    method?: string,
    cache?: RequestCache,
    headers?: Headers,
    body?: BodyInit,
    ignoreApiKey?: boolean,
    debugPrint?: boolean,
    skipJsonParse?: false,
}): Promise<T>;

export async function sendRequestV1<T>({
    route,
    method,
    cache,
    headers,
    body,
    ignoreApiKey,
    debugPrint,
    skipJsonParse,
}: {
    route: string,
    method?: string,
    cache?: RequestCache,
    headers?: Headers,
    body?: BodyInit,
    ignoreApiKey?: boolean,
    debugPrint?: boolean,
    skipJsonParse?: boolean,
}): Promise<T | undefined> {
    try {
        if (ignoreApiKey !== true) {
            const apiKey = cookies().get("apiKey")?.value;
            if (apiKey == undefined) {
                throw new AuthError("invalid api key");
            }

            if (headers === undefined) {
                headers = new Headers({
                    "x-api-key": apiKey,
                });
            } else {
                headers.append("x-api-key", apiKey);
            }
        }

        const response = await fetch(`${process.env.INTERNAL_API_HOST}/v1/${route}`, {
            method: method ?? "GET",
            cache: cache ?? 'no-store',
            headers: headers,
            body: body,
        });

        if (!response.ok) {
            if (response.status == 403) {
                throw new AuthError("Unauthenticated response from the server");
            }
            throw new Error(await response.text());
        }

        if (response.status === 204 || skipJsonParse) {
            return undefined;
        }

        const rawData = await response.json()
        if (debugPrint === true) {
            console.log(rawData);
        }

        return rawData as T;
    } catch (e) {
        if (e instanceof AuthError) {
            console.log("Authentication error, removing all cookies:", e);
            cookies().delete("cid");
            cookies().delete("apiKey");
            redirect("/login")
            throw new AuthError("Not Allowed.");
        } else if (e instanceof Error) {
            console.error(e);
            throw e;
        } else {
            throw new Error("an unknown error occurred");
        }
    }
}
