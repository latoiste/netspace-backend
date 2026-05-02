package model

import (
	"log"
)

// hadnles websocket connections manages client and rooms
type Hub struct {
	clients   map[*Client]bool
	rooms     map[string]map[*Client]bool
	Register  chan *Client
	Logout    chan *Client
	JoinRoom  chan *RoomRequest
	LeaveRoom chan *RoomRequest
	Broadcast chan *Message
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		rooms:     make(map[string]map[*Client]bool),
		Register:  make(chan *Client),
		Logout:    make(chan *Client),
		JoinRoom:  make(chan *RoomRequest),
		LeaveRoom: make(chan *RoomRequest),
		Broadcast: make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			log.Println("Register request: ", client.Username)
		case client := <-h.Logout:
			log.Println("Logout request: ", client.Username)
		case _ = <-h.JoinRoom:
			log.Printf("Join request received")
		case _ = <-h.LeaveRoom:
			log.Printf("Leave request received")
		case _ = <-h.Broadcast:
			log.Printf("broadcast mesasge wow")
		}
	}
}

func (h *Hub) handleRegister(client *Client) {

}

func (h *Hub) handleLogout(client *Client) {

}

func (h *Hub) handleClientJoin(request *RoomRequest) {

}

func (h *Hub) handleClientLeave(request *RoomRequest) {

}
