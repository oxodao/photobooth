package services

import (
	"fmt"
	"time"

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

func (pb *Photobooth) OnButtonPress(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Button press: ", string(msg.Payload()))

	switch string(msg.Payload()) {
	case "TAKE_PICTURE":
		pb.prv.Sockets.BroadcastTakePicture()
	case "DISPLAY_DEBUG":
		if pb.DisplayDebug {
			break
		}

		pb.DisplayDebug = true
		pb.prv.Sockets.BroadcastState()
		go func() {
			time.Sleep(30 * time.Second)
			pb.DisplayDebug = false
			pb.prv.Sockets.BroadcastState()
		}()
	}

	msg.Ack()
}

func (pb *Photobooth) OnSyncRequested(client mqtt.Client, mqt mqtt.Message) {
	fmt.Println("Sync requested")
}
