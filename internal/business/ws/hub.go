package ws

type Hub struct {
	Rooms map[string]*Room `json:"room"`
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}
