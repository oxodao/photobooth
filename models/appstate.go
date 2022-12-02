package models

type AppState struct {
	HardwareID   string  `json:"hwid" db:"hwid"`
	ApiToken     *string `json:"token" db:"token"`
	CurrentEvent *int64  `json:"-" db:"current_event"`
}
