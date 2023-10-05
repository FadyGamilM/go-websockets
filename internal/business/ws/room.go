package ws

type Room struct {
	ID      int64              `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}
