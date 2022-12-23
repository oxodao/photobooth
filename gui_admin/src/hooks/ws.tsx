import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";
import { useApi } from "./auth";

const CONNECTION_MESSAGES = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Websocket closing",
    [ReadyState.CLOSED]: "Websocket closed",
    [ReadyState.UNINSTANTIATED]: "Websocket uninstantiated"
};

type WebsocketProps = {
    lastMessage: WsMessage | null;
    appState: AppState | null;
    currentTime: string | null;
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
    sendMessage: () => { },
});

export default function WebsocketProvider({ children }: { children: ReactNode }) {
    const [killSwitch, setKillswitch] = useState<number>(-1);
    const { showError } = useApi();
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);

    const HOST = `ws://${window.location.host}/api/socket/admin`
    const { sendMessage, lastMessage, readyState } = useWebSocket(HOST);

    useEffect(() => {
        if ([ReadyState.CLOSING, ReadyState.CLOSED, ReadyState.UNINSTANTIATED].includes(readyState)) {
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

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        const data = JSON.parse(lastMessage.data);

        switch (data.type) {
            case "PING":
                sendMessage('{"type": "PONG"}')
                setContext({...ctx, currentTime: data.payload})
                break
            case "APP_STATE":
                setContext({...ctx, appState: data.payload})
                break
            case "ERR_MODAL":
                showError(data.payload, 'error');
                break
            case "EXPORT_STARTED":
                showError('Export started', 'info');
                break
            case "EXPORT_COMPLETED":
                showError('Export completed', 'success');
                break
        }
    }, [lastMessage]);

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendMessage(JSON.stringify({ type: msgType, payload: data })),
    }}>
        <>
            <TextLoader loading={readyState != ReadyState.OPEN} text={CONNECTION_MESSAGES[readyState]}>
                {children}
            </TextLoader>
        </>
    </WebsocketContext.Provider>
}

export function useWebsocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}
