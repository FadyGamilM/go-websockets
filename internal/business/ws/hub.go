package ws

import "fmt"

type Hub struct {
	Rooms map[string]*Room `json:"room"`
	// a channel to broadcast messages into room's clients
	Broadcast chan *Message
	// a channel to register a new client to a room
	Register chan *Client
	// a channel to remove a client from the room
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.Broadcast:
			// first check that the room exists (to safely broadcast a message to an existing room)
			room, ok := h.Rooms[fmt.Sprintf("%v", msg.RoomID)]
			if ok {
				// so the room we want to broadcast this msg into it is actually existing, broadcast the message to all its clients (if there is any)
				if len(room.Clients) > 0 {
					for _, client := range room.Clients {
						client.Message <- msg
					}
				}
			}
		case client := <-h.Register:
			// first check that the room exists (to safely broadcast a message to an existing room)
			room, ok := h.Rooms[fmt.Sprintf("%v", client.RoomID)]
			// if the room we want to broadcast this msg into it is actually existing
			if ok {
				// now check if the client we try to register is already registered before or not
				_, ok := room.Clients[fmt.Sprintf("%v", client.ClientID)]
				// if this client is not registered
				if !ok {
					// so register this user to the room
					room.Clients[fmt.Sprintf("%v", client.ClientID)] = client
				}
			}
		case client := <-h.Unregister:
			// first check that the room exists (to safely broadcast a message to an existing room)
			room, ok := h.Rooms[fmt.Sprintf("%v", client.RoomID)]
			// if the room we want to broadcast this msg into it is actually existing
			if ok {
				// now check if the client we try to register is already registered before or not
				_, ok := room.Clients[fmt.Sprintf("%v", client.ClientID)]
				// if this client is registered
				if ok {
					// so delete this user from the room
					delete(h.Rooms[fmt.Sprintf("%v", client.RoomID)].Clients, fmt.Sprintf("%v", client.ClientID))
					// broadcast to the other existing clients (if there is any) that this user has left the room
					h.Broadcast <- &Message{
						Content:  client.Username + " has left the channel",
						Username: client.Username,
						RoomID:   client.RoomID,
					}
					// close the messages channel of this removed 
					close(client.Message)
				}
			}
		}
	}
}
