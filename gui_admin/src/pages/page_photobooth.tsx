import { Button, Card, CardActions, CardContent, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Typography } from "@mui/material";
import { useState } from "react";
import { useWebsocket } from "../hooks/ws";
import App from "./App";

export default function PagePhotobooth() {
    const {appState} = useWebsocket();
    const [exportZipShown, setExportZipShown] = useState<boolean>(false);

    const remoteTakePicture = () => {
        // @TODO
    };

    const exportAsZip = () => {
        // @TODO: Through websocket until so that we can give the user a progress on the compression then give him a link that he can use to download as blob
    }

    return <App>
        <Card>
            <CardContent>
                <Typography variant="h2" fontSize={18}>Current event</Typography>
                <ul>
                    <li>Amount of picture handtaken: {appState?.app_state?.current_event?.amt_images_handtaken}</li>
                    <li>Amount of picture unattended: {appState?.app_state?.current_event?.amt_images_unattended}</li>
                </ul>
            </CardContent>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} onClick={() => remoteTakePicture}>Remote take a picture</Button>
            </CardActions>
        </Card>
        <Card>
            <CardActions>
                <Button style={{ width: '100%' }} color="error" onClick={() => setExportZipShown(true)}>Export as zip</Button>
            </CardActions>
        </Card>

        <Dialog open={exportZipShown} onClose={() => setExportZipShown(false)}>
            <DialogTitle>Export as zip</DialogTitle>
            <DialogContent>
                <DialogContentText>
                    You are trying to export the event {appState?.app_state?.current_event?.name}. <br />
                    This will create a zip with all the pictures and let you download, thus it could take a long time. <br />
                    Are you sure you want to continue ?
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setExportZipShown(false)}>Cancel</Button>
                <Button onClick={() => {
                    exportAsZip();
                    setExportZipShown(false);
                }} color="warning" autoFocus>Export zip</Button>
            </DialogActions>
        </Dialog>
    </App>
}