import { cookies } from "next/headers";
import LoginClient from "./login";

export default function Login() {
    const cookieStore = cookies()
    const customerId = cookieStore.get("cid")

    return <LoginClient cid={customerId?.value} />
}