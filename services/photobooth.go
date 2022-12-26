package services

import (
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/photobooth/logs"
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
	logs.Info("Sync requested")
}

func (pb *Photobooth) OnExportEvent(client mqtt.Client, msg mqtt.Message) {
	eventIdStr := string(msg.Payload())
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		logs.Error("Failed to export event: bad eventid => ", eventIdStr)
		logs.Error(err)

		return
	}

	event, err := orm.GET.Events.GetEvent(eventId)
	if err != nil {
		logs.Error("Failed to export event:", err)
		return
	}

	pb.prv.Sockets.BroadcastAdmin("EXPORT_STARTED", event)

	exportedEvent, err := (NewEventExporter(event)).Export()
	if err != nil {
		logs.Error(err)
		return
	}

	pb.prv.Sockets.BroadcastAdmin("EXPORT_COMPLETED", exportedEvent)
}
