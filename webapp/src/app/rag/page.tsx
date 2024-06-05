"use client"

import HandleRAG from '@/handlers/rag';
import { RAGRequest } from '@/types/rag';
import React, { KeyboardEvent, useState } from 'react';

export default function RAG() {
    const [input, setInput] = useState("")
    const [conversationId, setConversationId] = useState("")
    const [error, setError] = useState("")
    const [messages, setMessages] = useState<string[]>([])

    const enterKeyHandler = (event: KeyboardEvent<HTMLInputElement>) => {
        if (event.key === 'Enter') {
            handleSubmit()
        }
    };

    const handleSubmit = async () => {
        // clear the data
        console.log("clearing data ...")
        setMessages((prev) => [...prev, input])
        setInput("")

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
            setMessages((prev) => [...prev, response.data!.response])
            setConversationId(response.data!.conversationId)
        }
    }

    const getMessages = () => {
        const items = []
        for (let i = 0; i < messages.length; i++) {
            items.push(<div key={i} className='bg-gray-200 text-gray-800 p-2 rounded-xl w-fit justify-end'><p>{messages[i]}</p></div>)
        }
        return items
    }

    return <div className="h-screen w-screen bg-black">
        <div className="grid grid-rows-9 h-full">
            <div className="bg-green-400 row-span-8 overflow-scroll">
                {/* <div className="h-[2000px] w-[32px] bg-red-400">scroll</div> */}
                <div className="space-y-2 p-8 flex flex-col">
                    {getMessages()}
                </div>
            </div>
            <div className="bg-blue-400 row-span-1 flex justify-center items-center px-8 overflow-hidden">
                <div className="w-full p-3 pl-8 pr-3 rounded-full text-black max-w-[1000px] bg-white">
                    <div className="flex space-x-4">
                        <input
                            type="text"
                            onKeyDown={enterKeyHandler}
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder='Your query here ...'
                            className="bg-white w-full">
                        </input>
                        <button className="bg-gray-200 p-2 rounded-full" onClick={handleSubmit}>Submit</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
}