package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Manager struct {
	Clients map[*Client]bool
	//
	sync.RWMutex
}

func NewManager(r *gin.Engine) *Manager {
	wsRoutes := r.Group("/api/ws")
	m := &Manager{
		Clients: make(map[*Client]bool),
	}
	wsRoutes.GET("/", m.ServerWebSockets)
	wsRoutes.POST("/", m.ServerWebSockets)
	wsRoutes.GET("/clients", m.HandleGetClients)
	return m
}

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
