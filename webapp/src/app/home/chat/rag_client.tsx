"use client"

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
import useWebSocket from 'react-use-websocket';
import RagLLMSelector from './rag_llm_selector';
import { ModelRow } from '@/types/llm';
import { House, Settings } from 'lucide-react';
import Link from 'next/link';

export default function RagClient({ wsBaseUrl }: { wsBaseUrl: string }) {
    const queryClient = useQueryClient()

    // for controlling the socketUrl
    const [socketUrl, setSocketUrl] = useState("");

    const [isLoading, setIsLoading] = useState(true)
    const [input, setInput] = useState("")
    const [messages, setMessages] = useState<ConversationMessage[]>([])
    const [isFirstMessage, setIsFirstMessage] = useState(true)
    const [currentChatLLM, setCurrentChatLLM] = useState<ModelRow | undefined>(undefined)

    const textareaRef = useRef<HTMLTextAreaElement>(null);

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
                    queryClient.invalidateQueries({ queryKey: ['allConversations'] })
                    break
                case "titleUpdate":
                    // invalidate the title query
                    queryClient.invalidateQueries({ queryKey: ['allConversations'] })
                    break
                case "changeChatLLM":
                    console.log("CHANGED CHAT LLM")
                    setCurrentChatLLM(data.chatLLM!)
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

    const enterKeyHandler = (event: KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            handleSubmit();
        }
    };

    const handleSubmit = async () => {
        if (input.trim() === "") {
            return
        }
        // clear the data
        console.log("clearing data ...")
        setMessages((prev) => [...prev, { role: 1, message: input, index: messages.length }])
        setInput("")
        setTimeout(() => scrollToBottom(), 200)

        // send the request on the websocket
        setIsLoading(true)
        sendMessage(JSON.stringify({ "messageType": "ragMessage", "message": input }))

        // sendRequest(input)
    }

    const changeChatLLM = (model: ModelRow) => {
        sendMessage(JSON.stringify({ "messageType": "changeChatLLM", "chatLLMId": model.llm.id }))
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

    useEffect(() => {
        if (textareaRef.current) {
            textareaRef.current.style.height = 'auto';
            textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
        }
    }, [input]);

    const handleInput = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
        const element = event.target;
        element.style.height = 'auto';
        element.style.height = `${element.scrollHeight}px`;
        setInput(event.target.value);
    };

    return <div className="flex flex-col flex-grow h-full overflow-hidden">
        <div ref={scrollableDivRef} className="bg-bg flex-grow overflow-scroll">
            <div className="sticky top-0 p-4">
                <RagLLMSelector currLLM={currentChatLLM} onSelect={changeChatLLM} />
            </div>
            <div className="flex h-full justify-center items-start w-full">
                <div className="flex flex-col max-w-[800px] w-full h-full">
                    {getMessages()}
                    <div className="min-h-16" />
                </div>
            </div>
        </div>
        <div className="bg-background flex flex-col justify-center items-center px-8 pt-0 pb-4">
            <div className="w-full bg-secondary p-2 pl-8 pr-3 rounded-[30px] max-w-[800px]">
                <div className="flex space-x-4">
                    <div className="w-full grid place-items-center">
                        <textarea
                            ref={textareaRef}
                            rows={1}
                            onKeyDown={enterKeyHandler}
                            value={input}
                            onChange={handleInput}
                            placeholder="Message our AI"
                            // className="bg-secondary w-full min-h-[1rem] resize-none"
                            className='m-0 resize-none border-0 bg-transparent px-0 block'
                            style={{
                                width: '100%',
                                maxHeight: `${10 * 1.2}em`, // Approximate height of maxRows lines
                                overflowY: 'auto',
                                resize: 'none',
                                boxSizing: 'border-box',
                                height: 'auto'
                            }}
                        />
                    </div>
                    <button
                        className="bg-primary hover:opacity-70 text-primary-foreground w-10 h-10 rounded-full font-bold flex-shrink-0"
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