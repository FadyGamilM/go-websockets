package ws

type Room struct {
	ID      int64              `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

func NewRoom(roomID int64, roomName string) *Room {
	return &Room{
		ID:      roomID,
		Name:    roomName,
		Clients: make(map[string]*Client),
	}
}
