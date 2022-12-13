package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oxodao/photobooth/config"
	"github.com/oxodao/photobooth/models"
	"github.com/oxodao/photobooth/orm"
	"golang.org/x/exp/slices"
)

const SOCKET_TYPE_PHOTOBOOTH = "PHOTOBOOTH"
const SOCKET_TYPE_ADMIN = "ADMIN"

var SOCKET_TYPES = []string{
	SOCKET_TYPE_PHOTOBOOTH,
	SOCKET_TYPE_ADMIN,
}

type Sockets []*Socket

func (s Sockets) broadcastTo(to string, msgType string, data interface{}) {
	for _, socket := range s {
		if len(to) > 0 && socket.Type != to {
			continue
		}

		socket.Send(msgType, data)
	}
}

func (s Sockets) BroadcastPhotobooth(msgType string, data interface{}) {
	s.broadcastTo(SOCKET_TYPE_PHOTOBOOTH, msgType, data)
}

func (s Sockets) BroadcastTakePicture() {
	for _, socket := range s {
		if socket.Type != SOCKET_TYPE_PHOTOBOOTH {
			continue
		}

		socket.TakePicture()
	}
}

func (s Sockets) BroadcastAdmin(msgType string, data interface{}) {
	s.broadcastTo(SOCKET_TYPE_ADMIN, msgType, data)
}

func (s Sockets) BroadcastState() {
	for _, socket := range s {
		socket.sendState()
	}
}

type Socket struct {
	Type string
	Open bool
	Conn *websocket.Conn

	mtx *sync.Mutex
}

func (s *Socket) TakePicture() {
	if GET.Photobooth.IsTakingPicture || s.Type != SOCKET_TYPE_PHOTOBOOTH {
		return
	}

	GET.Photobooth.IsTakingPicture = true

	fmt.Println("Taking picture...")
	go func() {
		timeout := config.GET.Photobooth.DefaultTimer

		for timeout >= 0 {
			s.Send("TIMER", timeout)
			timeout--
			time.Sleep(1 * time.Second)
		}

		GET.Photobooth.IsTakingPicture = false
	}()
}

func (s *Socket) OnMessage(msg models.SocketMessage) {
	switch msg.MsgType {
	case "PONG":
		break
	case "GET_STATE":
		s.sendState()
	case "TAKE_PICTURE":
		s.TakePicture()
	case "SET_MODE":
		mode, ok := msg.Payload.(string)
		if !ok {
			s.Send("ERR_MODAL", "Bad request")
			break
		}

		if !slices.Contains(config.MODES, mode) {
			s.Send("ERR_MODAL", "Unknown mode")
			break
		}

		GET.Photobooth.CurrentMode = mode
		GET.Sockets.BroadcastState()
	case "SET_DATETIME":
		dt, ok := msg.Payload.(string)
		if !ok {
			s.Send("ERR_MODAL", "Bad request")
			break
		}

		time, err := time.Parse("2006-01-02 15:04:05", dt)
		if err != nil {
			s.Send("ERR_MODAL", "Failed to set date: "+err.Error())
			break
		}
		err = SetSystemDate(time)
		if err != nil {
			fmt.Println(err)
		}
	case "SET_EVENT":
		evtIdFloat, ok := msg.Payload.(float64)
		if !ok {
			s.Send("ERR_MODAL", "Failed to change event: Bad request")
			break
		}

		var evtId int64 = int64(evtIdFloat)

		evt, err := orm.GET.Events.GetEvent(evtId)
		if err != nil {
			s.Send("ERR_MODAL", "Failed to change event: "+err.Error())
			break
		}

		GET.Photobooth.CurrentState.CurrentEvent = &evtId
		GET.Photobooth.CurrentState.CurrentEventObj = evt

		err = orm.GET.AppState.SetState(GET.Photobooth.CurrentState)
		if err != nil {
			s.Send("ERR_MODAL", "Failed to save state: "+err.Error())
			break
		}

		GET.Sockets.BroadcastState()
	case "SHOW_DEBUG":
		(*GET.MqttClient).Publish("photobooth/button_press", 2, false, "DISPLAY_DEBUG")
	case "SHUTDOWN":
		(*GET.MqttClient).Publish("photobooth/button_press", 2, false, "SHUTDOWN")
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

func (s *Socket) sendState() error {
	settings := GET.GetFrontendSettings()
	if settings == nil {
		return errors.New("failed to send frontend_settings")
	}

	return s.Send("APP_STATE", settings)
}

func (p *Provider) Join(socketType string, socket *websocket.Conn) {
	sock := &Socket{
		Type: socketType,
		Conn: socket,
		Open: true,
		mtx:  &sync.Mutex{},
	}
	p.Sockets = append(p.Sockets, sock)

	go func() {
		for sock.Open {
			time.Sleep(1 * time.Second)
			sock.Send("PING", time.Now().Format("2006-01-02 15:04:05"))
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

	if socketType == SOCKET_TYPE_PHOTOBOOTH && config.GET.Photobooth.UnattendedInterval > 1 {
		go func() {
			for sock.Open {
				time.Sleep(time.Duration(config.GET.Photobooth.UnattendedInterval) * time.Minute)
				fmt.Println("Unattended picture")
				sock.Send("UNATTENDED_PICTURE", nil)
			}
		}()
	}

	sock.sendState()
}
