package chat

import (
	"encoding/json"
	"log"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/db"
	"github.com/latoiste/netspace/model"
)

type Hub struct {
	Clients    map[string]*Client
	repo       *db.Repository
	groups     map[string]*Group
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

	inviteToGroup chan api.InviteToGroup
	acceptInvite  chan GroupInvite
	leaveGroup    chan GroupInvite
	sendGroupMsg  chan model.GroupMessage

	notificationRead chan api.NotificationRead

	persistPrivateMsg       chan model.PrivateMessage
	persistPublicMsg        chan model.PublicMessage
	persistGroupMsg         chan model.GroupMessage
	persistGroup            chan model.Group
	persistDissolveGroup    chan model.Group
	persistNotification     chan model.Notification
	persistNotificationRead chan api.NotificationRead
}

func NewHub(repo *db.Repository, locationId int) *Hub {
	return &Hub{
		Clients:                 make(map[string]*Client),
		locationId:              locationId,
		repo:                    repo,
		groups:                  make(map[string]*Group),
		Register:                make(chan *Client),
		unregister:              make(chan *Client),
		broadcast:               make(chan api.WsEvent),
		sendPrivateMsg:          make(chan model.PrivateMessage),
		typingStart:             make(chan api.TypingEvent),
		typingStop:              make(chan api.TypingEvent),
		sendPublicMsg:           make(chan model.PublicMessage),
		publicTypingStart:       make(chan api.PublicUserTyping),
		publicTypingStop:        make(chan api.PublicUserTyping),
		inviteToGroup:           make(chan api.InviteToGroup),
		acceptInvite:            make(chan GroupInvite),
		leaveGroup:              make(chan GroupInvite),
		sendGroupMsg:            make(chan model.GroupMessage),
		notificationRead:        make(chan api.NotificationRead, 20),
		persistPrivateMsg:       make(chan model.PrivateMessage, 20),
		persistPublicMsg:        make(chan model.PublicMessage, 20),
		persistGroupMsg:         make(chan model.GroupMessage, 20),
		persistGroup:            make(chan model.Group, 20),
		persistDissolveGroup:    make(chan model.Group, 20),
		persistNotification:     make(chan model.Notification, 20),
		persistNotificationRead: make(chan api.NotificationRead, 20),
	}
}

func (h *Hub) run() {
	go h.persistPrivateMsgLoop()
	go h.persistPublicMsgLoop()
	go h.persistGroupMsgLoop()
	go h.persistGroupLoop()
	go h.persistDissolveGroupLoop()
	go h.persistNotificationLoop()
	go h.persistNotificationReadLoop()

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
			h.handleSendPrivateMsg(msg)

		case msg := <-h.typingStart:
			h.notifyTyping("user_typing", msg)

		case msg := <-h.typingStop:
			h.notifyTyping("user_stopped_typing", msg)

		case msg := <-h.sendPublicMsg:
			h.handleSendPublicMsg(msg)

		case msg := <-h.publicTypingStart:
			h.notifyPublicTyping("public_user_typing", msg)

		case msg := <-h.publicTypingStop:
			h.notifyPublicTyping("public_user_stopped_typing", msg)

		case msg := <-h.inviteToGroup:
			h.handleInviteToGroup(msg)

		case invite := <-h.acceptInvite:
			h.handleAcceptInvite(invite)

		case invite := <-h.leaveGroup:
			h.handleLeaveGroup(invite)

		case msg := <-h.sendGroupMsg:
			h.handleSendGroupMsg(msg)

		case notif := <-h.notificationRead:
			h.persistNotificationRead <- notif
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

	for _, group := range h.groups {
		_, ok := group.memberIds[client.UserId]
		if !ok {
			continue
		}
		h.broadcastGroupMemberLeft(group, client)
	}

	h.broadcast <- api.WsEvent{
		Event: "user_left",
		Data:  data,
	}
	return nil
}
