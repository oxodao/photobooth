import Webcam from 'react-webcam'
import EventName from './components/event_name';
import { useWebsocket } from './hooks/ws';
import Photobooth from './pages/photobooth';
import Quiz from './pages/quiz';

function App() {
  const {currentMode, appState} = useWebsocket();

  return <>
    <EventName/>

    {
      (currentMode == 'PHOTOBOOTH' || currentMode == 'DISABLED')
      && <Photobooth disabled={currentMode == 'DISABLED'} />
    }

    {
      currentMode == 'QUIZ'
      && <Quiz />
    }
  </>
}

export default App
