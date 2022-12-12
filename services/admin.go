package services

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/photobooth/config"
	"golang.org/x/exp/slices"
)

type Admin struct {
	prv *Provider
}

func (a *Admin) OnSetMode(c mqtt.Client, msg mqtt.Message) {
	mode := string(msg.Payload())
	if !slices.Contains(config.MODES, mode) {
		fmt.Println("given mode is not allowed")
		return
	}

	a.prv.Photobooth.CurrentMode = mode
	a.prv.Sockets.BroadcastState()
}
