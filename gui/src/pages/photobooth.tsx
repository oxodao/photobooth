import { useEffect, useRef, useState } from "react";
import Webcam from "react-webcam";
import '../assets/css/photobooth.scss';
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

export default function Photobooth({disabled}: {disabled: boolean}) {
    const webcamRef = useRef<Webcam>(null);
    const [timer, setTimer] = useState(-1);
    const { appState, lastMessage, sendMessage } = useWebsocket();

    const takePicture = (unattended: boolean) => {
        if (!webcamRef || !webcamRef.current) {
            return;
        }

        const imageSrc = webcamRef.current.getScreenshot();
        if (imageSrc) {
            // debugOpenImage(imageSrc);

            let form = new FormData();
            form.append('image', imageSrc);
            form.append('unattended', unattended ? 'true' : 'false')
            form.append('event', ''+appState?.current_event?.id)

            fetch('/api/picture', {
                method: 'POST',
                body: form,
            }).then(() => {
                setTimer(-1);
            }).catch(() => {
                setTimer(-1);
            })
        }
    };

    useEffect(() => {
        if (!lastMessage) {
            return;
        }

        if (lastMessage.type == 'TIMER') {
            setTimer(lastMessage.payload)
            if (lastMessage.payload === 0) {
                takePicture(false);
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
            width={1280}
            height={720}
            onClick={() => !disabled && sendMessage('TAKE_PICTURE')}
            screenshotFormat="image/jpeg"
            videoConstraints={{ facingMode: 'user' }}
        />
        {
            timer >= 0 &&
            <div className={`timer ${appState.use_hardware_flash ? '' : 'flash'}`}>{timer > 0 && timer}</div>
        }
    </div>
}