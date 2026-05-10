package chat

import (
	"fmt"
	"log"

	"github.com/latoiste/netspace/model"
)

// hadnles websocket connections manages client and rooms
type Hub struct {
	clients   map[*Client]bool
	rooms     map[string]map[*Client]bool
	Register  chan *Client
	Logout    chan *Client
	JoinRoom  chan *RoomRequest
	LeaveRoom chan *RoomRequest
	Broadcast chan *model.Message
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		rooms:     make(map[string]map[*Client]bool),
		Register:  make(chan *Client),
		Logout:    make(chan *Client),
		JoinRoom:  make(chan *RoomRequest),
		LeaveRoom: make(chan *RoomRequest),
		Broadcast: make(chan *model.Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client, ok := <-h.Register:
			if ok {
				log.Println("Register request received")
				h.clients[client] = true
				fmt.Println(h.clients)
			}
		case client, ok := <-h.Logout:
			if ok {
				log.Println("Logout request received")
				delete(h.clients, client)
				// client.logout <- true

				// TODO: remove dari room, close room kalo tinggal <= 1 org (add timeout dulu ga langsung close room)

			}
		case req := <-h.JoinRoom:
			log.Println("Join request received")
			h.handleJoinRoom(req)
		case req := <-h.LeaveRoom:
			log.Println("Leave request received")
			h.handleLeaveRoom(req)
		case _ = <-h.Broadcast:
			log.Println("broadcast mesasge wow")
		}
	}
}

func (h *Hub) handleJoinRoom(request *RoomRequest) {

}

func (h *Hub) handleLeaveRoom(request *RoomRequest) {

}
