import { ReactNode } from "react";
import { useWebsocket } from "../hooks/ws";

export default function Debug() {
  const {appState, currentTime} = useWebsocket();

  if (!appState) {
    return <div className="debug abstl">Something went wrong</div>
  }
  const datetime = currentTime ?? 'Datetime not available';
  const eventName = !!appState.app_state?.current_event ? appState.app_state.current_event.name : 'No event selected !';

  const D = (title: string, child: ReactNode) => <div><span style={{fontWeight: 'bold'}}>{title}</span>: {child}</div>

  return <>
    <div className="debug abstl">
      {<span>{eventName}</span>}
      {
        appState.debug && <>
          {D('Mode', <span>{appState.current_mode}</span>)}
          {D('Hardware flash', <span>{appState.use_hardware_flash ? 'true': 'false'}</span>)}
          {D('IPs', <ul>
              {
                appState.ip_addresses && Object.entries(appState.ip_addresses).filter(([_, x]) => x.length > 0).map(([key, inter]) => <li>
                    {key}: {inter.join(', ')}
                </li>)
              }
            </ul>)}
        </>
      }
    </div>
    <div className="debug abstr">
      <span>{datetime}</span>
      {
        appState.debug && <>
          {D('HWID', <span>{appState.app_state.hwid}</span>)}
          {D('Token', <span>{appState.app_state.token}</span>)}
        </>
      }
    </div>
  </>
}