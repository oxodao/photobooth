import { Alert, Snackbar } from "@mui/material";
import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { TextLoader } from "../components/loader";
import { AppState } from "../types/appstate";
import { WsMessage } from "../types/ws_message";
import { useApi } from "./auth";

/**
 * @TODO: Refacto lastError to have a state like in Prowty
 */

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
    const {showError} = useApi();
    const [ctx, setContext] = useState<WebsocketProps>(defaultState);

    const HOST = `ws://${window.location.host}/api/socket/admin`
    const { sendMessage, lastMessage, readyState } = useWebSocket(HOST);

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
        if (!lastMessage) {
            return;
        }


        const data = JSON.parse(lastMessage.data);
        let newCtx = { ...ctx, lastMessage: data };

        switch (data.type) {
            case "PING":
                sendMessage('{"type": "PONG"}')
                newCtx.currentTime = data.payload;
                break
            case "APP_STATE":
                newCtx.appState = data.payload;
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

        setContext(newCtx)
    }, [lastMessage]);

    return <WebsocketContext.Provider value={{
        ...ctx,
        sendMessage: (msgType: string, data?: any) => sendMessage(JSON.stringify({ type: msgType, payload: data })),
    }}>
        <>
            <TextLoader loading={readyState != ReadyState.OPEN} text={connectionStatus}>
                {children}
            </TextLoader>
        </>
    </WebsocketContext.Provider>
}

export function useWebsocket(): WebsocketContextProps {
    return useContext<WebsocketContextProps>(WebsocketContext);
}
