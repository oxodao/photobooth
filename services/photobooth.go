package services

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/orm"
)

type Photobooth struct {
	prv *Provider

	CurrentState    models.AppState
	IsTakingPicture bool
	CurrentMode     string
	DisplayDebug    bool
}

func (pb *Photobooth) OnSyncRequested(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Sync requested")
}

func (pb *Photobooth) OnExportEvent(client mqtt.Client, msg mqtt.Message) {
	eventIdStr := string(msg.Payload())
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		fmt.Println("Failed to export event: bad eventid => ", eventIdStr)
		fmt.Println(err)

		return
	}

	event, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		fmt.Println("Failed to export event:")
		fmt.Println(err)

		return
	}

	err = (NewEventExporter(event)).Export()
	if err != nil {
		fmt.Println(err)
		return
	}

	pb.prv.Sockets.BroadcastAdmin("EXPORT_COMPLETED", event)
}
