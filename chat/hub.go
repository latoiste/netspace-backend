package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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

	persistPrivateMsg chan model.PrivateMessage
	persistPublicMsg  chan model.PublicMessage
}

func NewHub(repo *db.Repository, locationId int) *Hub {
	return &Hub{
		Clients:           make(map[string]*Client),
		locationId:        locationId,
		repo:              repo,
		groups:            make(map[string]*Group),
		Register:          make(chan *Client),
		unregister:        make(chan *Client),
		broadcast:         make(chan api.WsEvent),
		sendPrivateMsg:    make(chan model.PrivateMessage),
		typingStart:       make(chan api.TypingEvent),
		typingStop:        make(chan api.TypingEvent),
		sendPublicMsg:     make(chan model.PublicMessage),
		publicTypingStart: make(chan api.PublicUserTyping),
		publicTypingStop:  make(chan api.PublicUserTyping),
		inviteToGroup:     make(chan api.InviteToGroup),
		acceptInvite:      make(chan GroupInvite),
		leaveGroup:        make(chan GroupInvite),
		persistPrivateMsg: make(chan model.PrivateMessage, 20),
		persistPublicMsg:  make(chan model.PublicMessage, 20),
	}
}

func (h *Hub) run() {
	go h.persistPrivateMsgLoop()
	go h.persistPublicMsgLoop()

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

			h.persistPublicMsg <- msg

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
		case msg := <-h.inviteToGroup:
			// buat group invite nanti
			_, ok := h.groups[msg.GroupId]
			if !ok {
				log.Println("Group not found")
				continue
			}

			for _, memberId := range msg.UserIds {
				_, ok := h.Clients[memberId]
				if !ok {
					log.Printf("Failed to send invite to %v, id not found\n", memberId)
					continue
				}
				log.Printf("Sent to %v", memberId)
				// TODO: send group invite
			}
		case invite := <-h.acceptInvite:
			sender, ok := h.Clients[invite.userId]
			if !ok {
				log.Println("Sender not found")
				continue
			}

			group, ok := h.groups[invite.groupId]
			if !ok {
				log.Println("Group not found")
				continue
			}

			_, ok = group.memberIds[sender.UserId]
			if ok {
				log.Printf("User %v is trying to join but already in group %v\n", sender.UserId, group.name)
				continue
			}

			newMember := api.MemberJoined{
				Id:     sender.UserId,
				Name:   sender.Name,
				Emoji:  sender.Emoji,
				IsHost: false,
			}

			for memberId := range group.memberIds {
				client, ok := h.Clients[memberId]
				if !ok {
					log.Printf("Sender id %v not found in group %v\n", memberId, group.name)
					continue
				}
				client.sendEvent("member_joined", newMember)
			}
			group.memberIds[sender.UserId] = true
			log.Println("Member is added to group")
			fmt.Println(group)
		case invite := <-h.leaveGroup:
			sender, ok := h.Clients[invite.userId]
			if !ok {
				log.Println("Sender not found")
				continue
			}

			group, ok := h.groups[invite.groupId]
			if !ok {
				log.Println("Group not found")
				continue
			}

			_, ok = group.memberIds[sender.UserId]
			if !ok {
				log.Printf("User %v is trying to leave but not in group %v\n", sender.UserId, group.name)
				continue
			}

			h.broadcastGroupMemberLeft(group, sender)
			if len(group.memberIds) <= 1 {
				// TODO: send group_dissolve
			}
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

func (h *Hub) persistPrivateMsgLoop() {
	for msg := range h.persistPrivateMsg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.InsertPrivateMessage(msg, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}
}

func (h *Hub) persistPublicMsgLoop() {
	for msg := range h.persistPublicMsg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.InsertPublicMessage(msg, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
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

func (h *Hub) createGroup(groupName string, hostId string, memberIds []string) (*Group, error) {
	if _, ok := h.Clients[hostId]; !ok {
		return nil, errors.New("Invalid host id")
	}

	for _, group := range h.groups {
		if group.name == groupName {
			return nil, errors.New("Group name is taken")
		}
	}

	validMemberIds := make([]string, len(memberIds))

	for _, memberId := range memberIds {
		if _, ok := h.Clients[memberId]; ok {
			// TODO: send group invite
			validMemberIds = append(validMemberIds, memberId)
		}
	}
	groupId := GenerateGroupId()

	group := &Group{
		id:        groupId,
		name:      groupName,
		hostId:    hostId,
		memberIds: make(map[string]bool),
	}
	group.memberIds[hostId] = true
	h.groups[groupId] = group

	return group, nil
}

func (h *Hub) broadcastGroupMemberLeft(group *Group, leavingClient *Client) {
	memberLeft := api.MemberLeft{
		UserId: leavingClient.UserId,
		Name:   leavingClient.Name,
	}

	for memberId := range group.memberIds {
		if memberId == leavingClient.UserId {
			continue
		}
		client, ok := h.Clients[memberId]
		if !ok {
			log.Printf("Sender id not found")
			continue
		}
		client.sendEvent("member_left", memberLeft)
	}

	delete(group.memberIds, leavingClient.UserId)
	fmt.Println(group)
}
