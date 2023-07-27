package chatsession

type Chatroom struct {
	Clients map[*UserClient]struct{}
	Client  chan *UserClient
	Relay   chan *UserMessage
	Remove  chan *UserClient
}

func CreateChatroom() *Chatroom {
	cr := &Chatroom{
		Clients: make(map[*UserClient]struct{}),
		Client:  make(chan *UserClient),
		Relay:   make(chan *UserMessage),
		Remove:  make(chan *UserClient),
	}

	go cr.RegisterClient()
	go cr.ReceiveAndNotify()
	go cr.RemoveClient()

	return cr
}

func (cr *Chatroom) RegisterClient() {
	for client := range cr.Client {
		if _, exists := cr.Clients[client]; !exists {
			cr.Clients[client] = struct{}{}
		}
	}
}

func (cr *Chatroom) ReceiveAndNotify() {
	for msg := range cr.Relay {
		for k := range cr.Clients {
			k.Message <- msg
		}
	}
}

func (cr *Chatroom) RemoveClient() {
	for client := range cr.Remove {
		if _, exists := cr.Clients[client]; exists {
			close(client.Message)
			delete(cr.Clients, client)
		}
	}
}
