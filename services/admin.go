package services

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/photobooth/config"
	"github.com/oxodao/photobooth/logs"
	"golang.org/x/exp/slices"
)

type Admin struct {
	prv *Provider
}

func (a *Admin) OnSetMode(c mqtt.Client, msg mqtt.Message) {
	mode := string(msg.Payload())
	if !slices.Contains(config.MODES, mode) {
		logs.Error("given mode is not allowed")
		return
	}

	a.prv.Photobooth.CurrentMode = mode
	a.prv.Sockets.BroadcastState()
}
