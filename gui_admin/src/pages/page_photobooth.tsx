import { Button, Card, CardActions, CardContent, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Paper, Table, TableBody, TableCell, TableContainer, TableHead, Typography } from "@mui/material";
import ExportListing from "../components/export_listing";
import { useWebsocket } from "../hooks/ws";
import App from "./App";

export default function PagePhotobooth() {
    const { appState, sendMessage } = useWebsocket();

    return <App>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>
                    Current event:  {!!appState?.app_state.current_event && <>{appState?.app_state.current_event?.name}</>}
                </Typography>
                <ul>
                    <li>Amount of picture handtaken: {appState?.app_state?.current_event?.amt_images_handtaken}</li>
                    <li>Amount of picture unattended: {appState?.app_state?.current_event?.amt_images_unattended}</li>
                </ul>
            </CardContent>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => sendMessage('REMOTE_TAKE_PICTURE', null)}>Remote take a picture</Button>
            </CardActions>
        </Card>
        {
            !!appState?.app_state?.current_event
            && <ExportListing />
        }
    </App>
}