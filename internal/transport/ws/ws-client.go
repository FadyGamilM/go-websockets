package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	Manager *Manager
	Conn    *websocket.Conn
	// egress is a sync mechanism, since that gorilla mux websocket conn allows for us only one concurrent write at a time, we have to read the messages one by one, and put them in an unbuffered channel, then read these messages (in the writeMessage method) one by one, so when a clinet_A want to send a message and client_B want to send a message at the same time, the one who comes first will put its message on the egress channel first, so the other client message will be wait untill the message of client_A get read from the unbuffered channel ..
	egress chan Event
}

func NewClient(m *Manager, c *websocket.Conn) *Client {
	return &Client{
		Manager: m,
		Conn:    c,
		egress:  make(chan Event),
	}
}

// so simply this method keep running in a for loop waitin for messages, and if we break from the for loop we reach the defer func which close the conn and remove this client
func (c *Client) ReadMessages() {
	// if we break from the below endless for loop, we want to clean the connection and resources because the connection is closed either normally or apnormally
	defer func() {
		// remove client will remove the client and close its connection
		c.Manager.RemoveClient(c)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses when we recieve a pong from the client
	c.Conn.SetPongHandler(c.pongHandler)

	for {
		_, msgPayload, err := c.Conn.ReadMessage()
		if err != nil {
			// if the connection is closed abnormly, without being closed by the client or by the server, so there is something wrong
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ==> %v", err)
			}
			// if its closed for normal reasons, its okay lets break from listening to incoming messages to the currenct websocket connection
			break
		}

		var e Event
		err = json.Unmarshal(msgPayload, &e)
		if err != nil {
			log.Printf("error marshling the msgPayload into the evnet type : %v \n", err)
			// break to clean up the resources
			break
		}
		// if we marshel it and the msgPayload is on the form of event (contains payload and type) -> then route it via the manager
		err = c.Manager.RouteEventsToHandlers(e, c)
		if err != nil {
			log.Printf("error routing the event to the appropriate event handler : %v \n", err)
			break
		}

		// TODO => this part is for just printing the messages to all clients (no essential and not supposed to be here i am doing this to ensure that my backend working correctly ..)
		// put the message from the client (react app) into egress channel which some user send it (broadcast it)
		for wsClient, _ := range c.Manager.Clients {
			if wsClient != c {

				wsClient.egress <- e
			}
		}
		// TODO => End of the broadcasting to clients part ..
	}
}

func (c *Client) WriteMessages() {
	// i defined this ticker and the hearbeat logic in the writeMessage because we need to ping (send) to the client and this is the role of the write message method
	// define a ticker to ping on a ping interval rate defined previously which is 90% of the pong-wait time
	heartbeatTicker := time.NewTicker(pingInterval)
	defer func() {
		// close the ticker to free the reousrce
		heartbeatTicker.Stop()
		// remove the client and close its connection to clean reousrces
		c.Manager.RemoveClient(c)
	}()
	for {
		select {
		case msgEvent, ok := <-c.egress:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("websocket connection has been closed before sending the close-message (ERROR): %v \n", err)
				}
				// this will break the for-loop and trigger the defer func which is responsible for cleaning up the resources
				return
			}

			// try to send the message back (the server [go] -> client [react])
			// we are sending the message to the client so we need to marshal it to convert it from the go-type into JSON-type
			data, err := json.Marshal(msgEvent)
			if err != nil {
				log.Printf("error marshling the recieved event into msg to be sent to other clients : %v \n", err)
				break // to clean up resources
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("error while trying to send the text message to the client (ERROR): %v \n", err)
			}
			// we succedded
			log.Println("message send ..")
		case <-heartbeatTicker.C:
			// ping the client
			log.Println("[Heartbeat] | ping")
			// Send the Ping
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("[Heartbeat] | the connection is closed via the client : ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("[Heartbeat] | pong")
	return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
}
