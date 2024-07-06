"use client"

import { handleRAG } from '@/actions/rag';
import { ConversationMessage } from '@/types/conversation';
import { RagMessagePayload } from '@/types/rag';
import React, { KeyboardEvent, useEffect, useRef, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import DefaultLoader from '@/components/default_loader';
import RagMessage from './rag_message';
import Cookies from "js-cookie"
import { getConversation } from '@/actions/conversation';
import RagEmpty from './rag_empty';
import { toast } from '@/components/ui/use-toast';
import useWebSocket, { ReadyState } from 'react-use-websocket';

export default function RagClient({ wsBaseUrl }: { wsBaseUrl: string }) {
    const queryClient = useQueryClient()

    // for controlling the socketUrl
    const [socketUrl, setSocketUrl] = useState("");

    const [isLoading, setIsLoading] = useState(true)
    const [input, setInput] = useState("")
    const [messages, setMessages] = useState<ConversationMessage[]>([])
    const [isFirstMessage, setIsFirstMessage] = useState(true)

    // websocket
    const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

    const scrollableDivRef = useRef<HTMLDivElement | null>(null);

    // fetch the conversation
    const conv = useQuery({
        queryKey: ['conversation'],
        queryFn: () => getConversation(),
    })

    // handle the loading from the `useQuery` call. This will call mutliple times
    // based on the loading state of the call, and eventually results in the messages
    // being loaded into state and the screen scrolling to the bottom
    useEffect(() => {
        if (conv.status === "success") {
            // set the websocket state
            setSocketUrl(`${wsBaseUrl}?id=${conv.data.conversationId}`)

            // set the message state
            setMessages(conv.data!.messages)
            setIsFirstMessage(conv.data!.messages.length === 0)
            setTimeout(() => scrollToBottom(), 200)
        }

        if (conv.status === "error") {
            console.error("there was an error with the query")

        }

        if (conv.status === "error" || conv.status === "success") {
            setIsLoading(false)
        }
    }, [conv.data, conv.status])


    // handle when new messages are recieved in the websocket
    useEffect(() => {
        if (lastMessage !== null) {
            // parse the base64 payload into json
            const base64String = lastMessage.data.trim().replace(/^"|"$/g, '');
            const message = atob(base64String);
            const data = JSON.parse(message) as RagMessagePayload
            console.log(data)

            // process the message based on the type
            switch (data.messageType) {
                case "loading":
                    setIsLoading(true)
                    break
                case "newMessage":
                    // recieved a new chat message from the connection
                    setMessages((prev) => prev.concat(data.chatMessage!))
                    setTimeout(() => scrollToBottom(), 200)

                    // handle when to stop loading
                    if (data.chatMessage!.role === 2) {
                        setIsLoading(false)
                    }
                    break
                case "newConversationId":
                    // set the conversation id as a cookie
                    Cookies.set("conversationId", data.conversationId!, { secure: true, sameSite: "strict" })
                    break
                case "titleUpdate":
                    // invalidate the title query
                    queryClient.invalidateQueries({ queryKey: ['allConversations'] })
                    break
                case "error":
                    // the conversation is in an errored state
                    console.error("there was an unexpected issue:", data.error!)
                    toast({
                        variant: "destructive",
                        title: "Oh no!",
                        description: <p>{data.error ?? "There was an unknown error"}</p>
                    })
                    setIsLoading(false)
                    break
                default:
                    console.log("unexpected message type:", data.messageType)
                    setIsLoading(false)
            }
        }
    }, [lastMessage])



    const scrollToBottom = () => {
        if (scrollableDivRef.current) {
            scrollableDivRef.current.scrollTo({ top: scrollableDivRef.current.scrollHeight, behavior: 'smooth' });
        }
    };

    const enterKeyHandler = (event: KeyboardEvent<HTMLInputElement>) => {
        if (event.key === 'Enter') {
            handleSubmit()
        }
    };

    const handleSubmit = async () => {
        // clear the data
        console.log("clearing data ...")
        setMessages((prev) => [...prev, { role: 1, message: input, index: messages.length }])
        setInput("")
        setTimeout(() => scrollToBottom(), 200)

        // send the request on the websocket
        setIsLoading(true)
        sendMessage(input)

        // sendRequest(input)
    }

    const sendRequest = async (input: string,) => {
        // ensure the websocket is ready
        if (readyState !== ReadyState.OPEN) {
            console.log("the websocket is not open")
            return
        }

        // send the request
        console.log("sending the request ...")
        try {
            let req = {
                input: input,
            }
            const response = await handleRAG(req)
            console.log("Success!")
            console.log(response)
            setMessages((prev) => [...prev, response.message])
            Cookies.set("conversationId", response.conversationId)
            scrollToBottom()

            // invalidate the conversation list query for re-rendering the sidebar
            if (isFirstMessage) {
                setIsFirstMessage(false)
                await queryClient.invalidateQueries({ queryKey: ['allConversations'] })
            }

            // check whether to auto send the response
            if (response.message.role == 3 || response.message.role == 4) {
                await sendRequest("")
            }

            // invalidate the conversation list query for re-rendering the sidebar
            if (isFirstMessage) {
                setIsFirstMessage(false)
                await queryClient.invalidateQueries({ queryKey: ['allConversations'] })
            }
        } catch (e) {
            if (e instanceof Error) console.error(e)
            toast({
                variant: "destructive",
                title: "Oh no!",
                description: <p>{e instanceof Error ? e.message : "Unknown error"}</p>
            })
        }

        setIsLoading(false)
    }

    const getMessages = () => {
        if (conv.status === "pending") {
            // TODO -- add a default loader here
            return <DefaultLoader />
        }

        if (messages.length === 0) {
            return <RagEmpty />
        }

        const items = []
        for (let i = 0; i < messages.length; i++) {
            // ignore system messages
            if (messages[i].role != 0) {
                items.push(<div key={`rag_message-${i}`}>
                    <RagMessage message={messages[i]} offset={messages.length - i - 1} />
                </div>)
            }
        }
        return items
    }

    return <div className="flex flex-col flex-grow h-full overflow-hidden">
        <div ref={scrollableDivRef} className="bg-bg flex-grow overflow-scroll p-4">
            <div className="flex h-full justify-center items-start w-full">
                <div className="flex flex-col pb-16 max-w-[800px] w-full h-full">
                    {getMessages()}
                </div>
            </div>
        </div>
        <div className="bg-background flex flex-col justify-center items-center px-8 pt-0 pb-4">
            <div className="w-full bg-secondary p-3 pl-8 pr-3 rounded-full max-w-[1000px]">
                <div className="flex space-x-4">
                    <input
                        type="text"
                        onKeyDown={enterKeyHandler}
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        placeholder="Your query here ..."
                        className="bg-secondary w-full"
                    />
                    <button
                        className="bg-primary text-primary-foreground w-10 h-10 rounded-full font-bold flex-shrink-0"
                        onClick={handleSubmit}
                    >
                        <div className='grid place-items-center'>
                            {isLoading || conv.isLoading ? <DefaultLoader /> : <p>&uarr;</p>}
                        </div>
                    </button>
                </div>
            </div>
            <div className="text-xs py-2 text-slate-500">AI can make mistakes, make sure you check important info.</div>
        </div>
    </div>
}