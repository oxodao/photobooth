package models

import "github.com/oxodao/photobooth/config"

type FrontendSettings struct {
	AppState   AppState                `json:"app_state"`
	Photobooth config.PhotoboothConfig `json:"photobooth"`

	DebugDisplay bool                `json:"debug"`
	IPAddress    map[string][]string `json:"ip_addresses"`
	KnownEvents  []Event             `json:"known_events"`
	KnownModes   []string            `json:"known_modes"`

	CurrentMode string `json:"current_mode"`

	PhotoboothVersion string `json:"photobooth_version"`
	PhotoboothCommit  string `json:"photobooth_commit"`
}

type AdminSettings struct {
	AvailableModes []string `json:"available_modes"`
}
