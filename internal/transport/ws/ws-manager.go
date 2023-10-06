package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// if we have db repos we will have it here ..
type Manager struct {
	Clients map[*Client]bool
	// to lock the Clients map from deadlock
	sync.RWMutex
	handlers map[string]EventHandler
}

// factory method
func NewManager(r *gin.Engine) *Manager {
	wsRoutes := r.Group("/api/ws")
	m := &Manager{
		Clients:  make(map[*Client]bool),
		handlers: make(map[string]EventHandler),
	}
	wsRoutes.GET("/", m.ServerWebSockets)
	wsRoutes.POST("/", m.ServerWebSockets)
	wsRoutes.GET("/clients", m.HandleGetClients)
	// setup the event handlers
	m.setupEventHandlers()
	return m
}

// setup each handler to its appropriate event type
func (m *Manager) setupEventHandlers() {
	// setup the event handler for send-msg event
	m.handlers[Event_SendMessage] = SendMessageEventHandler
}

func (m *Manager) RouteEventsToHandlers(e Event, c *Client) error {
	// if this type of events is supported , handle it
	handler, ok := m.handlers[e.Type]
	if ok {
		return handler(e, c)
	} else {
		// if not , return an error unsupported events
		return errors.New("unsupported event type")
	}
}

// endpoint of getting all the clients managed by the manager (All the connected clients)
func (m *Manager) HandleGetClients(c *gin.Context) {
	if len(m.Clients) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"clients": []struct{}{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"clients": len(m.Clients),
		})
	}
}

/*
@ Logic
  - upgrade the connection protocol to 101
  - add a client for the current ws request
  - run a read and write messages for this client on a concurrent go-routine
*/
func (m *Manager) ServerWebSockets(c *gin.Context) {
	// upgrade the connection protocol
	// upgrade the request to websocket protocol
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erorr trying to upgrade the request protocol",
		})
	}
	log.Printf("new websocket connection ..")

	// create a client for the currenct request
	client := NewClient(m, conn)

	m.AddClient(client)

	go client.ReadMessages()
	go client.WriteMessages()
}

func (m *Manager) AddClient(c *Client) {
	// when two requests coming at the same time we want to save the manager.clients data structure from being deadly-locked
	m.RWMutex.Lock()

	// add the client to the clients
	if _, ok := m.Clients[c]; ok {
		return
	}

	// add the client
	m.Clients[c] = true

	// unlock the data structure
	defer m.RWMutex.Unlock()
}

func (m *Manager) RemoveClient(c *Client) {
	// when two requests coming at the same time we want to save the manager.clients data structure from being deadly-locked
	m.RWMutex.Lock()

	// add the client to the clients
	if _, ok := m.Clients[c]; !ok {
		return
	}

	// add the client
	c.Conn.Close()
	delete(m.Clients, c)

	// unlock the data structure
	defer m.RWMutex.Unlock()
}
