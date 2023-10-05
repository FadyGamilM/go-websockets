package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientID   int64  `json:"client_id"`
	Username   string `json:"username"`
	RoomID     int64  `json:"room_id"`
	Connection *websocket.Conn
	Message    chan *Message
}

type Message struct {
	Content  string `json:"content"`
	Username string `json:"username"`
	RoomID   int64  `json:"room_id"`
}

// read message method

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		// close the connection for this client to the websocket
		c.Connection.Close()
	}()

	for {
		_, m, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ==> %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			Username: c.Username,
			RoomID:   c.RoomID,
		}

		hub.Broadcast <- msg
	}

}

// write message method
func (c *Client) WriteMessage() {
	defer func() {
		c.Connection.Close()
	}()
	for {
		msg, ok := <-c.Message
		if !ok {
			return
		}
		c.Connection.WriteJSON(msg)
	}
}
