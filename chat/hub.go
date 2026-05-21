package chat

import (
	"encoding/json"
	"log"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

type Hub struct {
	Clients    map[*Client]bool
	rooms      map[string]map[*Client]bool
	Register   chan *Client
	unregister chan *Client
	broadcast  chan api.WsEvent
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan api.WsEvent),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.unregister:
			err := h.removeClient(client)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (h *Hub) broadcastLoop() {
	for event := range h.broadcast {
		log.Println("Broadcast received")
		msg, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
			continue
		}
		for client := range h.Clients {
			log.Println("hah??")
			client.Send <- msg
		}
	}
}

func (h *Hub) AddClient(client *Client, user *model.User) error {
	userDTO := api.ConstructUserDTO(*user)
	userJoined := api.UserJoined{
		User: userDTO,
	}
	data, err := json.Marshal(userJoined)
	if err != nil {
		return err
	}

	h.Register <- client

	h.broadcast <- api.WsEvent{
		Event: "user_joined",
		Data:  data,
	}
	return nil
}

func (h *Hub) removeClient(client *Client) error {
	delete(h.Clients, client)
	close(client.Send)

	userLeft := api.UserLeft{
		UserId: client.UserId,
		Name:   client.Name,
	}

	data, err := json.Marshal(userLeft)
	if err != nil {
		return err
	}

	h.broadcast <- api.WsEvent{
		Event: "user_left",
		Data:  data,
	}
	return nil
}
