import Link from "next/link";
import { ThemeToggle } from "./theme_toggle";
import { Button } from "./ui/button";

export default function Footer() {
    return <footer className="p-4 border-t border-t-border">
        <div className="flex justify-between">
            <ThemeToggle />
            {/* <Button className="opacity-50" variant="link" asChild>
                <Link href="/beta-api-keys">Beta Api Keys</Link>
            </Button> */}
        </div>
    </footer>
}