import React from 'react'
import ReactDOM from 'react-dom/client'
import { ThemeProvider } from '@emotion/react'
import { createTheme, CssBaseline } from '@mui/material'
import App from './App'

import './assets/css/index.scss';
import WebsocketProvider from './hooks/ws';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <WebsocketProvider>
        <App />
      </WebsocketProvider>
    </ThemeProvider>
  </React.StrictMode>
)
