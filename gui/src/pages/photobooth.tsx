import { useEffect, useRef, useState } from "react";
import Webcam from "react-webcam";
import '../assets/css/photobooth.scss';
import LockedModal from "../components/locked_modal";
import { useWebsocket } from "../hooks/ws";

const debugOpenImage = (img: string) => {
    const contentType = 'image/jpeg';

    const byteCharacters = atob(img.substr(`data:${contentType};base64,`.length));
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 1024) {
        const slice = byteCharacters.slice(offset, offset + 1024);

        const byteNumbers = new Array(slice.length);
        for (let i = 0; i < slice.length; i++) {
            byteNumbers[i] = slice.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);

        byteArrays.push(byteArray);
    }
    const blob = new Blob(byteArrays, { type: contentType });
    const blobUrl = URL.createObjectURL(blob);

    window.open(blobUrl, '_blank');
};

export default function Photobooth() {
    const webcamRef = useRef<Webcam>(null);

    const [timer, setTimer] = useState(-1);
    const { appState, lastMessage, sendMessage } = useWebsocket();

    const [flash, setFlash] = useState<boolean>(false);

    const [lastPicture, setLastPicture] = useState<string|null>(null);

    const takePicture = async (unattended: boolean) => {
        if (!webcamRef || !webcamRef.current) {
            return;
        }

        const imageSrc = webcamRef.current.getScreenshot();
        if (imageSrc) {
            let form = new FormData();

            form.append('image', imageSrc);
            form.append('unattended', unattended ? 'true' : 'false')
            form.append('event', ''+appState?.app_state?.current_event?.id)

            try {
                const resp = await fetch('/api/picture', {
                    method: 'POST',
                    body: form,
                });
         
                setTimer(-1);

                if (!unattended) {
                    const image = await resp.blob();
                    const url = URL.createObjectURL(image);

                    setLastPicture(url);

                    setTimeout(() => setLastPicture(null), 3500);
                }
            } catch {
                setTimer(-1);
            }
        }
    };

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        if (lastMessage.type == 'TIMER') {
            setTimer(lastMessage.payload)
            if (lastMessage.payload === 0) {
                // takePicture(false);
                setFlash(true);
                setTimeout(() => {
                    takePicture(false);
                    setFlash(false);
                }, 150);
            }

            return
        }

        if (lastMessage.type == 'UNATTENDED_PICTURE') {
            takePicture(true);
        }
    }, [lastMessage]);

    if (!appState) {
        return <div className="photobooth">NO STATE !</div>;
    }

    return <div className="photobooth">
        <Webcam
            forceScreenshotSourceSize
            ref={webcamRef}
            width={appState.photobooth.webcam_resolution.width}
            height={appState.photobooth.webcam_resolution.height}
            onClick={() => appState.current_mode !== 'DISABLED' && sendMessage('TAKE_PICTURE')}
            screenshotFormat="image/jpeg"
            videoConstraints={{ facingMode: 'user', ...appState.photobooth.webcam_resolution }}
        />
        {
            timer >= 0 &&
            <div className={`timer`}>{timer > 0 && timer}</div>
        }

        {
            flash &&
            <div className="timer flash"></div>
        }

        {
            appState.current_mode === 'DISABLED' && <LockedModal />
        }

        {
            lastPicture && <div className="picture_frame">
                <img src={lastPicture} alt="Last picture" />
            </div>
        }
    </div>
}