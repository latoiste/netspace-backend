package chat

// hadnles websocket connections manages client and rooms
type Hub struct {
	Clients map[*Client]bool
	rooms   map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
		rooms:   make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {}
	}
}
