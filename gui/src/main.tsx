import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

import './assets/css/index.scss';
import WebsocketProvider from './hooks/ws';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <WebsocketProvider>
      <App />
    </WebsocketProvider>
  </React.StrictMode>
)
