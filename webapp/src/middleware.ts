import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// This function can be marked `async` if using `await` inside
export async function middleware(request: NextRequest) {
    // check the current customerId
    const customerId = request.cookies.get('cid')
    if (customerId == undefined || customerId.value == "") {
        // redirect to the login screen
        console.log("No customerId found")
        return NextResponse.redirect(new URL('/login', request.url))
    }

    // check the api key
    const apiKey = request.cookies.get('apiKey')
    if (apiKey == undefined || apiKey.value == "") {
        // redirect to the login screen
        console.log("No apiKey found")
        return NextResponse.redirect(new URL('/login', request.url))
    }

    return NextResponse.next()
}

// See "Matching Paths" below to learn more
export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api (API routes)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         */
        {
            source: '/((?!api|_next/static|_next/image|favicon.ico|login|beta-api-keys).*)',
            missing: [
                { type: 'header', key: 'next-router-prefetch' },
                { type: 'header', key: 'purpose', value: 'prefetch' },
            ],
        },
    ],
}