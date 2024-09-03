import { getUserIntegrations } from "@/actions/integrations"
import Icon from "@/components/icon"
import { Integration } from "@/types/integrations"
import Link from "next/link"

export default async function HomeView() {

    const getContent = async () => {
        try {
            const integrations = await getUserIntegrations()

            return <div className="grid grid-cols-3 gap-4">
                {integrations.map((item, i) => <IntegrationCell key={`integration-${i}`} integration={item} />)}
            </div>
        } catch (e) {
            return <div className="">Error getting integrations</div>
        }
    }

    return <div className="flex flex-col flex-grow h-screen">
        <div className="p-8">
            <div className="safe-area">
                <div className="space-y-2">
                    <h2 className="text-lg font-bold">Integrations</h2>
                    {await getContent()}
                </div>
            </div>
        </div>
    </div>
}

function IntegrationCell({
    integration,
}: {
    integration: Integration,
}) {
    return <Link href={integration.href}>
        <div className="group hover:border-primary hover:bg-primary hover:cursor-pointer transition-all border rounded-md p-4">
            <div className="flex space-x-2">
                <div className="resize-non pt-[3px]">
                    <Icon size={18} name={integration.icon} />
                </div>
                <div className="text-left">
                    <h3 className="font-medium">{integration.title}</h3>
                    <p className="text-sm text-left font-light opacity-70">{integration.description}</p>
                </div>
            </div>
        </div>
    </Link>
}