import { useEffect, useState } from "react";

import { Box, Button, Card, CardActions, CardContent, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, MenuItem, Modal, Select, SelectChangeEvent, Typography } from "@mui/material";
import { DateTime } from 'luxon';

import { useWebsocket } from "../hooks/ws";
import { Event } from '../types/appstate';
import App from "./App";

const style = {
    position: 'absolute' as 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 400,
    bgcolor: '#0e0e0e',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
};

export default function PageIndex() {
    const { sendMessage, appState, currentTime } = useWebsocket();

    const [knownEvents, setKnownEvents] = useState<Event[]>([]);
    const [currentEvent, setCurrentEvent] = useState<Event | null>(null);
    const [modes, setModes] = useState<string[]>([]);

    const [newEvent, setNewEvent] = useState<Event | null>(null);
    const [shutdown, setShutdown] = useState<boolean>(false);

    useEffect(() => {
        if (!appState) {
            setModes([]);
            setCurrentEvent(null);
            setKnownEvents([]);
            return;
        }

        setKnownEvents(appState.known_events);
        setModes(appState.known_modes);
        setCurrentEvent(appState.app_state.current_event);
    }, [appState]);

    const setMode = (evt: SelectChangeEvent) => sendMessage('SET_MODE', evt.target.value);

    const setEvent = (evt: SelectChangeEvent) => {
        const events = knownEvents.filter(x => x.id === (evt.target.value as unknown as number)); // wow such typescript
        if (events.length > 0) {
            if (!currentEvent) {
                sendMessage('SET_EVENT', events[0]);
            } else {
                setNewEvent(events[0]);
            }
        }
    };

    const submitDatetime = () => sendMessage('SET_DATETIME', DateTime.now().toFormat('yyyy-MM-dd HH:mm:ss'));

    return <App>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>Current event</Typography>
                {
                    <Select value={!!currentEvent ? "" + currentEvent.id : ""} label="Event" onChange={setEvent} style={{ marginTop: '1em' }}>
                        {
                            knownEvents.map(x => <MenuItem key={x.id} value={x.id}>{x.name}</MenuItem>)
                        }
                    </Select>
                }
            </CardContent>
        </Card>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>Mode</Typography>
                {
                    appState?.current_mode &&
                    <Select value={appState.current_mode} label="Mode" onChange={setMode} style={{ marginTop: '1em' }}>
                        {
                            modes.map(x => <MenuItem key={x} value={x}>{x}</MenuItem>)
                        }
                    </Select>
                }
            </CardContent>
        </Card>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>System time</Typography>
                <Typography variant="body1" style={{ textAlign: "center", marginTop: '2em' }}>{currentTime}</Typography>
            </CardContent>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={submitDatetime}>Set to my device's time</Button>
            </CardActions>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('SHOW_DEBUG', null)}>Show debug info (30 sec)</Button>
            </CardActions>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => setShutdown(true)}>Shutdown</Button>
            </CardActions>
        </Card>

        <Dialog open={!!newEvent} onClose={() => setNewEvent(null)}>
            <DialogTitle>Change event</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    You are updating the current event to "{newEvent?.name} (by {newEvent?.author})".
                    Doing so will make that all new pictures are sent to this event instead of the current one.
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setNewEvent(null)}>Cancel</Button>
                <Button onClick={() => {
                    sendMessage('SET_EVENT', newEvent?.id);
                    setNewEvent(null);
                }} autoFocus>Change event</Button>
            </DialogActions>
        </Dialog>

        <Dialog open={shutdown} onClose={() => setShutdown(false)}>
            <DialogTitle>Shutting down</DialogTitle>
            <DialogContent>
                <DialogContentText>You are trying to shutdown the photobooth. Are you sure ?</DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setShutdown(false)}>Cancel</Button>
                <Button onClick={() => {
                    sendMessage('SHUTDOWN', null);
                    setShutdown(false);
                }} autoFocus>Change event</Button>
            </DialogActions>
        </Dialog>
    </App>;
}