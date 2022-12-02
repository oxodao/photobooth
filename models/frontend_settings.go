package models

import "github.com/oxodao/photomaton/config"

type FrontendSettings struct {
	AppState     AppState                `json:"app_state"`
	CurrentEvent *Event                  `json:"current_event"`
	Photobooth   config.PhotoboothConfig `json:"photobooth"`
	DebugMode    bool                    `json:"debug"`
}
