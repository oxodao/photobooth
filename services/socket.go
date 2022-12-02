package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oxodao/photomaton/config"
	"github.com/oxodao/photomaton/models"
)

const SOCKET_TYPE_PHOTOBOOTH = "PHOTOBOOTH"

var SOCKET_TYPES = []string{
	SOCKET_TYPE_PHOTOBOOTH,
}

type Socket struct {
	Type string
	Open bool
	Conn *websocket.Conn

	isTakingPicture bool
	mtx             *sync.Mutex
}

func (s *Socket) OnMessage(msg models.SocketMessage) {
	switch msg.MsgType {
	case "PONG":
		break
	case "GET_STATE":
		s.SendState()
	case "TAKE_PICTURE":
		if s.isTakingPicture {
			return
		}
		s.isTakingPicture = true

		fmt.Println("Taking picture...")
		go func() {
			timeout := config.GET.Photobooth.DefaultTimer

			for timeout >= 0 {
				s.Send("TIMER", timeout)
				timeout--
				time.Sleep(1 * time.Second)
			}

			s.isTakingPicture = false
		}()
	case "":
		// Probably should be handled in another way
		return
	default:
		fmt.Printf("Unhandled socket message: %v\n", msg.MsgType)
		fmt.Println(msg)
	}
}

func (s *Socket) Send(msgType string, data interface{}) error {
	s.mtx.Lock()
	err := s.Conn.WriteJSON(models.SocketMessage{
		MsgType: msgType,
		Payload: data,
	})
	s.mtx.Unlock()

	return err
}

func (s *Socket) SendState() error {
	settings := GET.GetFrontendSettings()
	if settings == nil {
		return errors.New("failed to send frontend_settings")
	}

	return s.Send("APP_STATE", settings)
}
