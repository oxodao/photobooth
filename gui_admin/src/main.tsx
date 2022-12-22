import React from 'react'
import ReactDOM from 'react-dom/client'
import { ThemeProvider } from '@emotion/react'
import { createTheme, CssBaseline } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterLuxon } from '@mui/x-date-pickers/AdapterLuxon';
import { createHashRouter, RouterProvider } from "react-router-dom";

import ApiProvider from './hooks/auth'
import WebsocketProvider from './hooks/ws'

import './assets/css/index.scss'
import PageIndex from './pages/page_index';
import PagePhotobooth from './pages/page_photobooth';
import PageQuiz from './pages/page_quiz';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

const router = createHashRouter([
  {
    path: '/',
    element: <PageIndex />,
  },
  {
    path: '/photobooth',
    element: <PagePhotobooth />
  },
  {
    path: '/quiz',
    element: <PageQuiz />
  },
]);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <LocalizationProvider dateAdapter={AdapterLuxon}>
        <ApiProvider>
          <WebsocketProvider>
            <RouterProvider router={router} />
          </WebsocketProvider>
        </ApiProvider>
      </LocalizationProvider>
    </ThemeProvider>
  </React.StrictMode>,
)

