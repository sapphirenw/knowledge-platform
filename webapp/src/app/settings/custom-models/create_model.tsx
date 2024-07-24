"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"

import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Slider } from "@/components/ui/slider"

import { Input } from "@/components/ui/input"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { createLLM, getAvailableModels } from "@/actions/llm"
import DefaultLoader from "@/components/default_loader"
import AvailableModelView from "@/components/available_model"
import { Button } from "@/components/ui/button"
import { AvailableModel, CreateLLMRequest } from "@/types/llm"
import { Textarea } from "@/components/ui/textarea"
import { useState } from "react"
import { toast } from "@/components/ui/use-toast"


const FormSchema = z.object({
    title: z.string().min(3, {
        message: "The title must be at least 3 characters",
    }),
    availableModelName: z.string().min(1, {
        message: "The model name cannot be empty"
    }),
    temperature: z.number(),
    instructions: z.string().min(10, {
        message: "The message must be at least 10 characters"
    })
})

export default function CreateCustomerLLM() {
    const queryClient = useQueryClient()

    const [openDialog, setOpenDialog] = useState(false)
    const [isLoading, setIsLoading] = useState(false)

    // fetch the available models

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            availableModelName: "claude-3-5-sonnet-20240620",
            temperature: 1,
        },
    })

    async function onSubmit(data: z.infer<typeof FormSchema>) {
        // validate
        const req: CreateLLMRequest = {
            availableModelName: data.availableModelName,
            title: data.title,
            temperature: data.temperature,
            instructions: data.instructions
        }
        if (req.availableModelName.includes("claude")) {
            req.temperature = req.temperature / 2
        }

        // send the request
        setIsLoading(true)
        try {
            await createLLM(req)

            // invalidate query and close
            await queryClient.invalidateQueries({ queryKey: ['customerLLMs', false] })
            toast({
                title: "Success!",
                description: <p>Successfully created the LLM.</p>
            })
            setOpenDialog(false)
            form.reset()
        } catch (e) {
            if (e instanceof Error) console.log(e)
            toast({
                variant: "destructive",
                title: "On no!",
                description: <p>There was an issue creating your LLM.</p>
            })
        }
        setIsLoading(false)
    }

    return <Dialog open={openDialog} onOpenChange={setOpenDialog}>
        <DialogTrigger asChild>
            <Button onClick={() => setOpenDialog(true)}>Create New Model</Button>
        </DialogTrigger>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>Create A New LLM</DialogTitle>
                <DialogDescription>
                    Define a custom LLM that you can use to control various processes throughout this application. Whether that be chatting, summary generation, content generation, etc.
                </DialogDescription>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                    <FormField
                        control={form.control}
                        name="title"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Title</FormLabel>
                                <FormControl>
                                    <Input placeholder="My Awesome New Model" {...field} />
                                </FormControl>
                                <FormDescription>
                                    The title of your custom model. Make it creative and distinct!
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="availableModelName"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Language Model</FormLabel>
                                <Select onValueChange={field.onChange} defaultValue={field.value}>
                                    <FormControl>
                                        <SelectTrigger>
                                            <SelectValue placeholder="Claude-3.5 Sonnet" />
                                        </SelectTrigger>
                                    </FormControl>
                                    <AvailableModels />
                                </Select>
                                <FormDescription>
                                    Which lanuage model from a provider you want to use
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="temperature"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Temperature - {field.value}</FormLabel>
                                <Slider defaultValue={[field.value]} min={0} max={2} step={0.1} onValueChange={(val) => field.onChange(val[0])} />
                                <FormDescription>
                                    How creative you want the model to be. Higher values are more creative, but less predictable.
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="instructions"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Personality</FormLabel>
                                <FormControl>
                                    <Textarea
                                        placeholder="Give verbose instructions for the model to follow."
                                        className="resize-none"
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <DialogFooter>
                        <Button type="submit">
                            <div className="flex space-x-2 items-center">
                                {isLoading ? <DefaultLoader /> : <></>}
                                <p>Submit</p>
                            </div>
                        </Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
}

function AvailableModels() {
    const { status, data, error } = useQuery({
        queryKey: ['availableModels'],
        queryFn: () => getAvailableModels(""),
    })

    if (status === "pending") {
        return <DefaultLoader />
    }

    if (status === "error") {
        console.error(error)
        return <p>Failed to get the available models</p>
    }

    const grouped = data.reduce((acc: { [key: string]: AvailableModel[] }, obj) => {
        const key = obj.provider
        if (!acc[key]) {
            acc[key] = [];
        }
        acc[key].push(obj);
        return acc;
    }, {})

    const getItems = () => {
        const items = []

        for (const key in grouped) {
            items.push(
                <SelectGroup key={key}>
                    <SelectLabel>{key}</SelectLabel>
                    {grouped[key].map((item, i) => <SelectItem key={(item.id)} value={item.id}>{item.displayName}</SelectItem>)}
                </SelectGroup>
            )
        }

        return items
    }

    return <SelectContent>
        {getItems()}
    </SelectContent>

}