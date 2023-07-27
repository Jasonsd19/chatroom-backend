package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jasonsd19/chatroom-backend/internal/chatsession"
)

func main() {
	setupHandlers()

	log.Print("Starting server...")
	err := http.ListenAndServe(":3001", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server has shutdown")
	} else if err != nil {
		log.Println("Error - ", err)
		os.Exit(1)
	}
}

func setupHandlers() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	cr := chatsession.CreateChatroom()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error - ", err)
			return
		}

		userClient := chatsession.CreateUserClient(cr, conn)
		cr.Client <- userClient
	})
}
