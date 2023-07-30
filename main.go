package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jasonsd19/chatroom-backend/internal/chatsession"
)

type WebsocketHandler struct {
	Chatroom *chatsession.Chatroom
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/ws", &WebsocketHandler{Chatroom: chatsession.CreateChatroom()})

	log.Printf("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server has shutdown")
	} else if err != nil {
		log.Println("Error - ", err)
		os.Exit(1)
	}
}

func (wh *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isValid(r, wh.Chatroom) {
		http.Error(w, "Invalid username", 400)
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	username := r.URL.Query().Get("username")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error - ", err)
		return
	}

	userClient := chatsession.CreateUserClient(username, wh.Chatroom, conn)
	wh.Chatroom.Client <- userClient
}

func isValid(r *http.Request, cr *chatsession.Chatroom) bool {
	if username := r.URL.Query().Get("username"); username != "" {
		if len(username) >= 3 && len(username) <= 15 {
			if _, exists := cr.Clients[username]; !exists {
				return true
			}
		}
	}
	return false
}
