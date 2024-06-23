import { cookies } from 'next/headers'

export default function Login() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")


    return <div className="grid place-items-center">Login</div>
}