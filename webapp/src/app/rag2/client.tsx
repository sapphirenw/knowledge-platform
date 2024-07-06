"use client"

import { useQueryClient } from "@tanstack/react-query"
import React, { useCallback, useEffect, useRef, useState } from "react"
import Cookies from "js-cookie"
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export default function Rag2Client({ wsUrl }: { wsUrl: string }) {
    //Public API that will echo messages sent to it back to the client

    const [socketUrl, setSocketUrl] = useState(wsUrl);
    const [messageHistory, setMessageHistory] = useState<MessageEvent<any>[]>([]);
    const [message, setMessage] = useState("")

    const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl);

    useEffect(() => {
        if (lastMessage !== null) {
            console.log(lastMessage)
            setMessageHistory((prev) => prev.concat(lastMessage));
        }
    }, [lastMessage]);

    const handleClickSendMessage = () => {
        console.log("sending:", message)
        sendMessage(JSON.stringify({ "message": message }))
        setMessage("")
    }

    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    return (
        <div className="safe-area w-full space-y-2">
            <div className="">
                <span>The WebSocket is currently {connectionStatus}</span>
            </div>
            <Input onChange={(e) => setMessage(e.target.value)} value={message} />
            <Button onClick={() => handleClickSendMessage()}>Send Message</Button>
            {lastMessage ? <p>Last message: {lastMessage.data}</p> : null}
            <div className="">
                <p>Messages:</p>
                <ul>
                    {messageHistory.map((message, idx) => (
                        <p key={idx}>{message ? message.data : null}</p>
                    ))}
                </ul>
            </div>
        </div>
    );
};