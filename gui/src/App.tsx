import Debug from './components/debug';
import Loader from './components/loader';
import { useWebsocket } from './hooks/ws';
import Photobooth from './pages/photobooth';
import Quiz from './pages/quiz';

function App() {
  const {appState} = useWebsocket();

    

  return <Loader loading={!appState}>
    <Debug />

    {
      (appState?.current_mode === 'PHOTOBOOTH' || appState?.current_mode === 'DISABLED')
      && <Photobooth />
    }

    {
      appState?.current_mode === 'QUIZ'
      && <Quiz />
    }
  </Loader>
}

export default App
