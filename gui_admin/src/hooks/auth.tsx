import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import Login from "../pages/Login";

type AuthProps = {
    connecting: boolean;
    password: string|null;
    lastAuthError: string|null;
};

type AuthContextProps = AuthProps & {
    connect: (password: string) => void;
    logout: () => void;
};

const defaultState: AuthProps = {
    connecting: false,
    password: null,
    lastAuthError: null,
};

const AuthContext = createContext<AuthContextProps>({
    ...defaultState,
    connect: () => {},
    logout: () => {},
});

export default function AuthProvider({ children }: { children: ReactNode }) {
    const [ctx, setContext] = useState<AuthProps>(defaultState);

    useEffect(() => {
        const password = localStorage.getItem('password');
        if (!!password){
            setContext({...ctx, password});
        }
    }, []);

    const connect = async (password: string) => {
        setContext({...ctx, connecting: true});

        const newCtx = {...ctx, connecting: false};

        const resp = await fetch(`/api/admin/login`, {
            'method': 'POST',
            'headers': {
                'Authorization': password,
            }
        })

        if (resp.status === 401) {
            newCtx.lastAuthError = 'Wrong password';
        } else {
            const data = await resp.text();

            if (data !== 'yes') {
                newCtx.lastAuthError = 'Wrong response from the photobooth';
            } else {
                localStorage.setItem('password', password);
                newCtx.lastAuthError = null;
                newCtx.password = password;
            }
        }

        setContext(newCtx);
    };

    const logout = () => {
        localStorage.removeItem('password');
        setContext({...ctx, password: null});
    };

    return <AuthContext.Provider value={{
        ...ctx,
        connect,
        logout,
    }}>
        { (!!ctx.password && !ctx.lastAuthError) && <>{children}</> }

        { (!ctx.password || !!ctx.lastAuthError) && <Login /> }
    </AuthContext.Provider>
}

export function useAuth(): AuthContextProps {
    return useContext<AuthContextProps>(AuthContext);
}