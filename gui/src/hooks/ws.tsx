import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";

type WebsocketProps = {
    lastMessage: WsMessage|null;
    appState: AppState|null;
    currentTime: string|null;
};

type WebsocketContextProps = WebsocketProps & {
    sendMessage: (msgType: string, data?: any | null) => void;
};

const defaultState: WebsocketProps = {
    lastMessage: null,
    appState: null,
    currentTime: null,
};

const WebsocketContext = createContext<WebsocketContextProps>({
    ...defaultState,
    sendMessage: (msgType: string, data?: any) => { },
});

export default function WebsocketProvider({ children }: { children: ReactNode }) {
    const HOST = `ws://${window.location.host}/api/socket/photobooth`
    const [killSwitch, setKillswitch] = useState<number>(-1);
    const { sendMessage, lastMessage, readyState } = useWebSocket(HOST);
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);

    useEffect(() => {
        if ([
            ReadyState.CLOSING,
            ReadyState.CLOSED,
            ReadyState.UNINSTANTIATED,
        ].includes(readyState)) {
            setKillswitch(setTimeout(() => {
                window.location.reload();
            }, 5000));
        } else {
            if (killSwitch >= 0) {
                clearTimeout(killSwitch);
                setKillswitch(-1);
            }
        }
    }, [readyState]);

    const connectionStatus = {
        [ReadyState.CONNECTING]: "Connecting",
        [ReadyState.OPEN]: "Open",
        [ReadyState.CLOSING]: "Websocket closing",
        [ReadyState.CLOSED]: "Websocket closed",
        [ReadyState.UNINSTANTIATED]: "Websocket uninstantiated"
    }[readyState];

    useEffect(() => {
        if (!lastMessage){
            return;
        }


        const data = JSON.parse(lastMessage.data);
        let newCtx = {...ctx, lastMessage: data};

        switch (data.type){
            case "PING":
                sendMessage('{"type": "PONG"}')
                newCtx.currentTime = data.payload;
                break
            case "APP_STATE":
                newCtx.appState = data.payload;
                break
        }

        setContext(newCtx)
    }, [lastMessage]);

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendMessage(JSON.stringify({type: msgType, payload: data})),
    }}>
        <TextLoader loading={readyState != ReadyState.OPEN} text={connectionStatus}>
            {children}
        </TextLoader>
    </WebsocketContext.Provider>
}

export function useWebsocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}