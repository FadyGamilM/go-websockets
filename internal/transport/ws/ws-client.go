package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Manager *Manager
	Conn    *websocket.Conn
}

func NewClient(m *Manager, c *websocket.Conn) *Client {
	return &Client{
		Manager: m,
		Conn:    c,
	}
}

// so simply this method keep running in a for loop waitin for messages, and if we break from the for loop we reach the defer func which close the conn and remove this client
func (c *Client) ReadMessages() {
	// if we break from the below endless for loop, we want to clean the connection and resources because the connection is closed either normally or apnormally
	defer func() {
		// remove client will remove the client and close its connection
		c.Manager.RemoveClient(c)
	}()
	for {
		msgType, msgPayload, err := c.Conn.ReadMessage()
		if err != nil {
			// if the connection is closed abnormly, without being closed by the client or by the server, so there is something wrong
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ==> %v", err)
			}
			// if its closed for normal reasons, its okay lets break from listening to incoming messages to the currenct websocket connection
			break
		}

		log.Printf("The Message Type ==> %v", msgType)
		log.Printf("The Message Payload ==> %v", string(msgPayload))
	}
}
