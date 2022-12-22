import { AlertColor } from "@mui/material";

export type SnackbarData = {
    open: boolean;
    message: string|null;
    type: AlertColor;
};