package handlers

import (
	"fmt"
	"net/http"

	"github.com/FadyGamilM/go-websockets/internal/business/ws"
	"github.com/FadyGamilM/go-websockets/internal/core"
	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&createRoomDto); err != nil {
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
