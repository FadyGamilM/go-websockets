package ws

type Hub struct {
	Rooms map[string]*Room `json:"room"`
}
