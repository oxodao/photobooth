import { useWebsocket } from "../hooks/ws";

export default function EventName() {
  const {currentMode, appState} = useWebsocket();

  let state = 'App state not available !';

  if (appState) {
    if (!appState.current_event) {
        state = 'No event selected !';
    } else {
        state = appState.current_event.name;
    }
  }

  return <span id="EVENT_NAME">{state}</span>
}