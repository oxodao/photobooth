package services

import (
	"fmt"
	"os/exec"
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
	case "SHUTDOWN":
		if err := GET.Shutdown(); err != nil {
			pb.prv.Sockets.broadcastTo("", "ERR_MODAL", "Failed to shutdown: "+err.Error())
			return
		}

		exec.Command("shutdown", "-h", "now").Run()
	}

	msg.Ack()
}

func (pb *Photobooth) OnSyncRequested(client mqtt.Client, mqt mqtt.Message) {
	fmt.Println("Sync requested")
}
