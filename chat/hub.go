package chat

import (
	"encoding/json"
	"log"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

type Hub struct {
	Clients    map[string]*Client
	rooms      map[string]map[*Client]bool
	Register   chan *Client
	unregister chan *Client

	broadcast      chan api.WsEvent
	sendPrivateMsg chan PrivateMessage
	typingStart    chan TypingEvent
	typingStop     chan TypingEvent
}

func NewHub() *Hub {
	return &Hub{
		Clients:        make(map[string]*Client),
		rooms:          make(map[string]map[*Client]bool),
		Register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan api.WsEvent),
		sendPrivateMsg: make(chan PrivateMessage),
		typingStart:    make(chan TypingEvent),
		typingStop:     make(chan TypingEvent),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserId] = client
		case client := <-h.unregister:
			err := h.removeClient(client)
			if err != nil {
				log.Println(err)
				continue
			}
		case msg := <-h.sendPrivateMsg:
			sender, ok := h.Clients[msg.SenderId]
			if !ok {
				log.Println("Sender id not found")
				continue
			}

			recipient, ok := h.Clients[msg.RecipientId]
			if !ok {
				log.Println("Recipient id not found")
				continue
			}

			newMsg := api.NewMessage{
				MessageId: msg.MessageId,
				Message:   msg.Message,
				Timestamp: msg.Timestamp,
			}

			msgSent := api.MessageSent{
				MessageId: msg.MessageId,
				Timestamp: msg.Timestamp,
			}

			senderMsg := newMsg
			senderMsg.IsMine = true

			recipientMsg := newMsg
			recipientMsg.IsMine = false

			err := recipient.sendEvent("new_message", recipientMsg)
			if err != nil {
				log.Println(err)
				continue
			}

			err = sender.sendEvent("new_message", senderMsg)
			if err != nil {
				log.Println(err)
				continue
			}

			err = sender.sendEvent("message_sent", msgSent)
			if err != nil {
				log.Println(err)
				continue
			}
		case msg := <-h.typingStart:
			h.notifyTyping("user_typing", msg)
		case msg := <-h.typingStop:
			h.notifyTyping("user_stopped_typing", msg)
		}
	}
}

func (h *Hub) broadcastLoop() {
	for event := range h.broadcast {
		log.Println("Broadcast received")
		message, err := json.Marshal(event)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, client := range h.Clients {
			client.Send <- message
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
	delete(h.Clients, client.UserId)
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

func (h *Hub) notifyTyping(event string, msg TypingEvent) {
	recipient, ok := h.Clients[msg.recipientId]
	if !ok {
		log.Println("Recipient not found")
		return
	}

	userTyping := api.UserTyping{
		UserId: msg.senderId,
	}

	err := recipient.sendEvent(event, userTyping)
	if err != nil {
		log.Println(err)
		return
	}
}
