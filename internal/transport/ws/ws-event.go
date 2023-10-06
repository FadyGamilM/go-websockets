package ws

import (
	"encoding/json"
	"log"
)

type Event struct {
	// we as a backend recieves this type and route the event to the appropriate event-handler
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// event handler is a function that performs a specific action based on the event type it receives
type EventHandler func(event Event, c *Client) error

const (
	Event_SendMessage = "send message"
)

// each event expects a different payload
// => payload of the send-message event
type SendMessageEventPayload struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

// => the send message event handler
func SendMessageEventHandler(event Event, c *Client) error {
	log.Printf("The received event is => %v \n", event)
	return nil
}
