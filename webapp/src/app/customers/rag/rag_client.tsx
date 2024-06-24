"use client"

import { HandleRAG } from '@/actions/rag';
import { ConversationMessage } from '@/types/conversation';
import React, { KeyboardEvent, useState } from 'react';
import RagMessage from './rag_message';

export default function RagClient({
    convId,
    msgs,
}: {
    convId: string | undefined
    msgs: ConversationMessage[]
}) {
    const [input, setInput] = useState("")
    const [conversationId, setConversationId] = useState(convId ?? "")
    const [error, setError] = useState("")
    const [messages, setMessages] = useState<ConversationMessage[]>(msgs)
    const [init, setInit] = useState(0)

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

        // send the request
        sendRequest(conversationId, input)
    }

    const sendRequest = async (conversationId: string, input: string,) => {
        // send the request
        console.log("sending the request ...")
        let req = {
            input: input,
            conversationId: conversationId,
        }
        let response = await HandleRAG(req)

        // parse the response
        if (response.error) {
            console.log("There was an error: ", response.error)
            setError(response.error)
        } else {
            console.log("Success!")
            console.log(response.data!)
            setMessages((prev) => [...prev, response.data!.message])
            setConversationId(response.data!.conversationId)

            // check whether to auto send the response
            if (response.data!.message.role == 3 || response.data!.message.role == 4) {
                sendRequest(response.data!.conversationId, "")
            }
        }
    }

    const getMessages = () => {
        const items = []
        for (let i = 0; i < messages.length; i++) {
            // ignore system messages
            if (messages[i].role != 0) {
                items.push(<div key={i}>
                    <RagMessage message={messages[i]} />
                </div>)
            }
        }
        return items
    }

    return <div className="flex flex-col h-full">
        <div className="bg-bg flex-grow overflow-scroll">
            <div className="flex h-full justify-center items-start">
                <div className="flex flex-col pb-16 max-w-[800px]">
                    {getMessages()}
                </div>
            </div>
        </div>
        <div className="bg-bg flex flex-col justify-center items-center px-8 pt-0 pb-4">
            <div className="w-full bg-container p-3 pl-8 pr-3 rounded-full max-w-[1000px]">
                <div className="flex space-x-4">
                    <input
                        type="text"
                        onKeyDown={enterKeyHandler}
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        placeholder="Your query here ..."
                        className="bg-container w-full"
                    />
                    <button
                        className="bg-slate-400 text-bg w-10 h-10 rounded-full font-bold flex-shrink-0"
                        onClick={handleSubmit}
                    >
                        &uarr;
                    </button>
                </div>
            </div>
            <div className="text-xs py-2 text-slate-500">AI can make mistakes, make sure you check important info.</div>
        </div>
    </div>
}