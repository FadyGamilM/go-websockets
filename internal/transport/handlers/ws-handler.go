package handlers

import (
	"fmt"
	"net/http"

	"github.com/FadyGamilM/go-websockets/internal/business/ws"
	"github.com/FadyGamilM/go-websockets/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsHandler struct {
	hub *ws.Hub
}

type WsHandlerConfig struct {
	R   *gin.Engine
	Hub *ws.Hub
}

func NewWsHandler(wshc *WsHandlerConfig) *wsHandler {
	handler := wsHandler{
		hub: wshc.Hub,
	}

	wsRoutes := wshc.R.Group("/api/ws")
	wsRoutes.POST("/room", handler.HandleCreateRoom)

	return &handler
}

func (wsh *wsHandler) HandleCreateRoom(c *gin.Context) {
	var createRoomDto core.CreateRoomReqDto
	if err := c.ShouldBind(&createRoomDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	room := ws.NewRoom(createRoomDto.ID, createRoomDto.Name)
	// add the room to the hub
	wsh.hub.Rooms[fmt.Sprintf("%v", room.ID)] = room

	c.JSON(http.StatusCreated, gin.H{
		"data": wsh.hub,
	})
}

type JoinRoomReqDto struct {
	RoomID   int64  `uri:"room_id" binding:"min=0, required"`
	ClientID int64  `form:"client_id" binding:"min=0, required"`
	Username string `form:"username" binding:"required"`
}

// api/ws/rooms/:roomID?userid=1&username=fady
func (wsh *wsHandler) HandleJoinRoom(c *gin.Context) {
	// upgrade the request to websocket protocol
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not upgrade the request to websocket protocol",
		})
	}

	var dto JoinRoomReqDto
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "bad request",
			},
		)
	}
	// create a client
	client := ws.Client{
		ClientID:   dto.ClientID,
		Username:   dto.Username,
		RoomID:     dto.RoomID,
		Connection: connection,
		Message:    make(chan *ws.Message, 10),
	}

	messageToBeBroadcasted := &ws.Message{
		Content:  "a new user has joined the room, say hi to " + client.Username,
		Username: client.Username,
		RoomID:   client.RoomID,
	}

	// register the client to the register-channel
	// broadcast the message
}
