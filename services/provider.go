package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oxodao/photomaton/config"
	"github.com/oxodao/photomaton/models"
	"github.com/oxodao/photomaton/orm"
	"github.com/oxodao/photomaton/utils"
)

var GET *Provider

type Provider struct {
	Sockets []*Socket
}

func (p *Provider) Join(socketType string, socket *websocket.Conn) {
	sock := &Socket{Type: socketType, Conn: socket, Open: true, mtx: &sync.Mutex{}}
	p.Sockets = append(p.Sockets, sock)

	go func() {
		for sock.Open {
			time.Sleep(5 * time.Second)
			sock.Send("PING", nil)
		}
	}()

	go func() {
		for {
			data := models.SocketMessage{}
			err := socket.ReadJSON(&data)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println("Unexpected close error: ", err)
					sock.Open = false
					return
				} else if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println("Websocket disconnected: ", err)
					sock.Open = false
					return
				}

				fmt.Println(err)
				continue
			}

			sock.OnMessage(data)
		}
	}()

	if config.GET.Photobooth.UnattendedInterval > 1 {
		go func() {
			for sock.Open {
				time.Sleep(time.Duration(config.GET.Photobooth.UnattendedInterval) * time.Minute)
				fmt.Println("Unattended picture")
				sock.Send("UNATTENDED_PICTURE", nil)
			}
		}()
	}

	sock.SendState()
}

func (p *Provider) GetFrontendSettings() *models.FrontendSettings {
	state, err := orm.GET.AppState.GetState()
	if err != nil {
		fmt.Println("Failed to get state: ", err)
		return nil
	}

	settings := models.FrontendSettings{
		AppState:     state,
		CurrentEvent: nil,
		Photobooth:   config.GET.Photobooth,
		DebugMode:    config.GET.DebugMode,
	}

	if state.CurrentEvent != nil {
		evt, err := orm.GET.Events.GetEvent(*state.CurrentEvent)
		if err != nil {
			fmt.Println("Failed to get current event: ", err)
			return nil
		}

		settings.CurrentEvent = evt
	}

	return &settings
}

func Load() error {
	for _, folder := range []string{"images"} {
		if err := utils.MakeOrCreateFolder(folder); err != nil {
			return err
		}
	}

	err := orm.Load()
	if err != nil {
		return err
	}

	prv := Provider{
		Sockets: []*Socket{},
	}

	GET = &prv

	return nil
}
