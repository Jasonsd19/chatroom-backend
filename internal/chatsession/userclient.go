package chatsession

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	messageSizeLimit = 256
	pongInterval     = 300 * time.Second
)

type UserClient struct {
	Chatroom   *Chatroom
	Connection *websocket.Conn
	Message    chan *UserMessage
}

type UserMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func CreateUserClient(cr *Chatroom, conn *websocket.Conn) *UserClient {
	// We could use a buffered channel for Message if we expect a large amount of traffic to the
	// chatroom however in my use case it's unlikely we will ned that much throughput
	uc := UserClient{Chatroom: cr, Connection: conn, Message: make(chan *UserMessage)}

	go uc.ReadMessage()
	go uc.WriteMessage()

	return &uc
}

func (uc *UserClient) ReadMessage() {
	defer func() {
		uc.Chatroom.Remove <- uc
		uc.Connection.Close()
	}()

	uc.Connection.SetReadLimit(messageSizeLimit)
	uc.Connection.SetReadDeadline(time.Now().Add(pongInterval))
	uc.Connection.SetPongHandler(func(string) error { uc.Connection.SetReadDeadline(time.Now().Add(pongInterval)); return nil })

	for {
		var userMessage UserMessage
		err := uc.Connection.ReadJSON(&userMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error -", err)
			}
			return
		}
		uc.Chatroom.Relay <- &userMessage
	}
}

func (uc *UserClient) WriteMessage() {
	for msg := range uc.Message {
		err := uc.Connection.WriteJSON(msg)
		if err != nil {
			log.Println("Error -", err)
			return
		}
	}
}
