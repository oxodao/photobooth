import { ReactNode, useState } from 'react';

import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import LogoutIcon from '@mui/icons-material/Logout';


import { useWebsocket } from '../hooks/ws';
import { useAuth } from '../hooks/auth';
import CustomDrawer from '../components/drawer';
import { Stack } from '@mui/material';

function App({ children }: { children: ReactNode }) {
  const { logout } = useAuth();
  const [menuOpened, setMenuOpened] = useState<boolean>(false);
  const { appState } = useWebsocket();


  return (
    <div className="App">
      <AppBar position="static">
        <Toolbar>
          <IconButton size="large" edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }} onClick={() => setMenuOpened(true)}><MenuIcon /></IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>Mode: {appState?.current_mode}</Typography>
          <IconButton size="large" edge="start" color="inherit" aria-label="menu" onClick={logout}><LogoutIcon /></IconButton>
        </Toolbar>
      </AppBar>

      <CustomDrawer open={menuOpened} close={() => setMenuOpened(false)} />

      <Stack maxWidth="sm" spacing={2} margin="auto" marginTop={2}>
        {children}
      </Stack>
    </div>
  )
}

export default App
