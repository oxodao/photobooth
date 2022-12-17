package services

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/photobooth/models"
)

type Photobooth struct {
	prv *Provider

	CurrentState    models.AppState
	IsTakingPicture bool
	CurrentMode     string
	DisplayDebug    bool
}

func (pb *Photobooth) OnSyncRequested(client mqtt.Client, mqt mqtt.Message) {
	fmt.Println("Sync requested")
}
