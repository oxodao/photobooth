package models

type SocketMessage struct {
	MsgType string      `json:"type"`
	Payload interface{} `json:"payload"`
}
