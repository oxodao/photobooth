type appstate = {
    hwid: string;
};

type Event = {
    id: number;
    name: string;
    author: string;
    date: number;
    location?: string|null;
};

export type AppState = {
    app_state: appstate;
    current_event: Event|null;
    debug: boolean;
    use_hardware_flash: boolean;
};