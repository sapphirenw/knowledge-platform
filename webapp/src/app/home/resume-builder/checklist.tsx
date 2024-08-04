"use client"

import { getResumeChecklist } from "@/actions/resume"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { useQuery } from "@tanstack/react-query"
import { Square, SquareCheckBig } from "lucide-react"

export default function ResumeChecklistView() {
    const { data, status } = useQuery({
        queryKey: ['resumeChecklist'],
        queryFn: () => getResumeChecklist(),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        return <ErrorPage msg="" />
    }

    return <div className="space-y-16 max-w-2xl mx-auto">
        <StepCell
            checked={data[0].completed}
            step={1}
            desc="Upload your resume and some additional information about you that our models will parse to fill out your resume sections."
        >
            <div className=""></div>
        </StepCell>
        <StepCell
            checked={data[1].completed}
            step={2}
            desc="Verify the personal information that exists at the top of your resume. Make sure this information is accurate!"
        >
            <div className="grid place-items-center">
                <Button asChild>
                    <a href="/home/resume-builder/about">
                        <p>About Me</p>
                    </a>
                </Button>
            </div>
        </StepCell>
        <StepCell
            checked={data[2].completed}
            step={3}
            desc="Edit and tweak the various sections of your base resume. Make sure that our models have access to as much information as possible for these sections."
        >
            <div className=""></div>
        </StepCell>
    </div>
}

function StepCell({
    checked,
    step,
    desc,
    children,
}: {
    checked: boolean,
    step: number,
    desc: string,
    children: JSX.Element,
}) {
    return <div className="space-y-4">
        <div className="flex items-center space-x-4">
            <div className="my-auto">
                {checked ? <SquareCheckBig className="text-primary" /> : <Square />}
            </div>
            <div className="">
                <h3 className="font-medium opacity-50 text-sm">Step {step}</h3>
                <p className="text-lg">{desc}</p>
            </div>
        </div>
        <div className="h-[100px]">
            <div className="flex">
                <Separator orientation="vertical" />
                <div className="w-full">
                    <div className="mx-auto">
                        {children}
                    </div>
                </div>
            </div>
        </div>
    </div>
}