import { useCallback, useEffect, useLayoutEffect, useRef } from 'react';

export default function useKeyPress (keys: string[], callback: (event: any) => void, node = null) {
    const callbackRef = useRef(callback);

    useLayoutEffect(() => {
        callbackRef.current = callback;
    });

    const handleKeyPress = useCallback(
        (event: any) => {
            if (keys.some((key) => event.key === key)) {
                callbackRef.current(event);
            }
        },
        [keys]
    );

    useEffect(() => {
        const targetNode = node ?? document;
        targetNode && targetNode.addEventListener("keydown", handleKeyPress);

        return () => targetNode && targetNode.removeEventListener("keydown", handleKeyPress);
    }, [handleKeyPress, node]);
};
