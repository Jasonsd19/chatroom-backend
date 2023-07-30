package chatsession

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	messageSizeLimit = 256
	messageInterval  = 300 * time.Second
)

type UserClient struct {
	Username   string
	Chatroom   *Chatroom
	Connection *websocket.Conn
	Message    chan *UserMessage
	ClientList chan *ClientList
}

type UserMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type ClientList struct {
	Clients []string `json:"user_list"`
}

func CreateUserClient(username string, cr *Chatroom, conn *websocket.Conn) *UserClient {
	// We could use a buffered channel for Message if we expect a large amount of traffic to the
	// chatroom however in my use case it's unlikely we will need that much throughput
	uc := UserClient{Username: username, Chatroom: cr, Connection: conn, Message: make(chan *UserMessage), ClientList: make(chan *ClientList)}

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
	uc.Connection.SetReadDeadline(time.Now().Add(messageInterval))

	for {
		var userMessage UserMessage
		err := uc.Connection.ReadJSON(&userMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error -", err)
			}
			return
		}
		uc.Connection.SetReadDeadline(time.Now().Add(messageInterval))
		uc.Chatroom.Relay <- &userMessage
	}
}

func (uc *UserClient) WriteMessage() {
	defer func() {
		uc.Chatroom.Remove <- uc
		uc.Connection.Close()
	}()

	for {
		select {
		case msg := <-uc.Message:
			err := uc.Connection.WriteJSON(msg)
			if err != nil {
				log.Println("Error -", err)
				return
			}
		case userList := <-uc.ClientList:
			err := uc.Connection.WriteJSON(userList)
			if err != nil {
				log.Println("Error -", err)
				return
			}
		}
	}
}
