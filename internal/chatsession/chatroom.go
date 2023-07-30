package chatsession

type Chatroom struct {
	Clients map[string]*UserClient
	Client  chan *UserClient
	Relay   chan *UserMessage
	Remove  chan *UserClient
}

func CreateChatroom() *Chatroom {
	cr := &Chatroom{
		Clients: map[string]*UserClient{},
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
		if _, exists := cr.Clients[client.Username]; !exists {
			cr.Clients[client.Username] = client
			cr.SendUserList()
		}
	}
}

func (cr *Chatroom) ReceiveAndNotify() {
	for msg := range cr.Relay {
		for _, client := range cr.Clients {
			client.Message <- msg
		}
	}
}

func (cr *Chatroom) RemoveClient() {
	for client := range cr.Remove {
		if _, exists := cr.Clients[client.Username]; exists {
			close(client.Message)
			delete(cr.Clients, client.Username)
			cr.SendUserList()
		}
	}
}

func (cr *Chatroom) SendUserList() {
	userList := ClientList{Clients: make([]string, len(cr.Clients))}

	index := 0
	for k := range cr.Clients {
		userList.Clients[index] = k
	}

	for _, client := range cr.Clients {
		client.ClientList <- &userList
	}
}
