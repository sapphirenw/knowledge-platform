"use client"

import { attachDocumentsToResume, getResumeChecklist, getResumeDocuments } from "@/actions/resume"
import DefaultLoader from "@/components/default_loader"
import ErrorPage from "@/components/error_page"
import SelectDocumentModel from "@/components/select_document/select_doc_model"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { ResumeChecklistItem } from "@/types/resume"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { ChevronUp, Square, SquareCheckBig } from "lucide-react"
import { useState } from "react"

export default function ResumeChecklistView() {
    const queryClient = useQueryClient()

    const { data, status } = useQuery({
        queryKey: ['resumeChecklist'],
        queryFn: () => getResumeChecklist(),
    })

    const resumeDocsResponse = useQuery({
        queryKey: ['resumeDocuments'],
        queryFn: () => getResumeDocuments(),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        return <ErrorPage msg="" />
    }

    return <div className="space-y-8 max-w-2xl mx-auto">
        <StepCell
            item={data[0]}
            step={1}
            desc="Upload your resume and some additional information about you that our models will parse to fill out your resume sections."
        >
            <div className="grid place-items-center">
                <div className="flex space-x-2">
                    {resumeDocsResponse.status === "success" ? <SelectDocumentModel
                        title="Choose Resume"
                        limit={1}
                        initalSelected={resumeDocsResponse.data}
                        onSave={async (docs) => {
                            await attachDocumentsToResume(docs)
                            await queryClient.invalidateQueries({ queryKey: ['resumeChecklist'] })
                        }}
                    /> : <DefaultLoader />}
                    <Button variant="outline">
                        <p>Choose Other Information</p>
                    </Button>
                </div>
            </div>
        </StepCell>
        <StepCell
            item={data[1]}
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
            item={data[2]}
            step={3}
            desc="Edit and tweak the various sections of your base resume. Make sure that our models have access to as much information as possible for these sections."
        >
            <div className=""></div>
        </StepCell>
    </div>
}

function StepCell({
    item,
    step,
    desc,
    children,
}: {
    item: ResumeChecklistItem,
    step: number,
    desc: string,
    children: JSX.Element,
}) {

    const [isCollapsed, setIsCollapsed] = useState(item.completed)

    return <div className="space-y-4">
        <div className="flex items-center space-x-4">
            <div className="my-auto">
                {item.completed ? <SquareCheckBig className="text-primary" /> : <Square />}
            </div>
            <div className="">
                <h3 className="font-medium opacity-50 text-sm">Step {step}</h3>
                <p className="text-lg">{desc}</p>
            </div>
            <div className={`${isCollapsed ? "rotate-180" : ""} transition-all`}>
                <button onClick={() => setIsCollapsed(!isCollapsed)}>
                    <ChevronUp />
                </button>
            </div>
        </div>
        <div className={`${isCollapsed ? "h-0 opacity-0 invisible" : "h-[100px] opacity-100"} transition-all`}>
            <div className="flex">
                <Separator orientation="vertical" />
                <div className="w-full">
                    <div className="mx-auto text-center space-y-4">
                        <p className="text-red-400 font-medium">{item.message}</p>
                        <div>{children}</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
}