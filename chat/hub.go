package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/db"
	"github.com/latoiste/netspace/model"
)

type Hub struct {
	Clients    map[string]*Client
	repo       *db.Repository
	rooms      map[string]map[*Client]bool
	Register   chan *Client
	unregister chan *Client
	locationId int

	broadcast      chan api.WsEvent
	sendPrivateMsg chan model.PrivateMessage
	typingStart    chan api.TypingEvent
	typingStop     chan api.TypingEvent

	sendPublicMsg     chan model.PublicMessage
	publicTypingStart chan api.PublicUserTyping
	publicTypingStop  chan api.PublicUserTyping

	persistPrivateMsg chan model.PrivateMessage
}

func NewHub(repo *db.Repository, locationId int) *Hub {
	return &Hub{
		Clients:           make(map[string]*Client),
		locationId:        locationId,
		repo:              repo,
		rooms:             make(map[string]map[*Client]bool),
		Register:          make(chan *Client),
		unregister:        make(chan *Client),
		broadcast:         make(chan api.WsEvent),
		sendPrivateMsg:    make(chan model.PrivateMessage),
		typingStart:       make(chan api.TypingEvent),
		typingStop:        make(chan api.TypingEvent),
		sendPublicMsg:     make(chan model.PublicMessage),
		publicTypingStart: make(chan api.PublicUserTyping),
		publicTypingStop:  make(chan api.PublicUserTyping),
		persistPrivateMsg: make(chan model.PrivateMessage, 20),
	}
}

func (h *Hub) run() {
	go func() {
		for msg := range h.persistPrivateMsg {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			err := h.repo.InsertPrivateMessage(msg, ctx)
			if err != nil {
				log.Println(err)
			}
			cancel()
		}
	}()

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

			h.persistPrivateMsg <- msg

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
		case msg := <-h.sendPublicMsg:
			sender, ok := h.Clients[msg.SenderId]
			if !ok {
				log.Println("Sender id not found")
				continue
			}

			newMsg := api.NewPublicMessage{
				MessageId:   msg.MessageId,
				SenderId:    msg.SenderId,
				SenderName:  sender.Name,
				SenderEmoji: sender.Emoji,
				Message:     msg.Message,
				Timestamp:   msg.Timestamp,
			}

			senderMsg := newMsg
			senderMsg.IsMine = true

			otherMsg := newMsg
			otherMsg.IsMine = false

			//TODO: persist public message

			err := sender.sendEvent("new_public_message", senderMsg)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, client := range h.Clients {
				if client.UserId != msg.SenderId {
					err = client.sendEvent("new_public_message", otherMsg)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		case msg := <-h.publicTypingStart:
			h.notifyPublicTyping("public_user_typing", msg)
		case msg := <-h.publicTypingStop:
			h.notifyPublicTyping("public_user_stopped_typing", msg)
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

func (h *Hub) notifyTyping(event string, msg api.TypingEvent) {
	recipient, ok := h.Clients[msg.RecipientId]
	if !ok {
		log.Println("Recipient not found")
		return
	}

	userTyping := api.UserTyping{
		UserId: msg.SenderId,
	}

	err := recipient.sendEvent(event, userTyping)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Hub) notifyPublicTyping(event string, msg api.PublicUserTyping) {
	for _, client := range h.Clients {
		if client.UserId != msg.UserId {
			err := client.sendEvent(event, msg)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
