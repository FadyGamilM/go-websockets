package ws

import "github.com/gorilla/websocket"

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
