package chat

import (
	"log"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Hub) handleSendPrivateMsg(msg model.PrivateMessage) {
	sender, ok := h.Clients[msg.SenderId]
	if !ok {
		log.Println("Sender id not found")
		return
	}

	recipient, ok := h.Clients[msg.RecipientId]
	if !ok {
		log.Println("Recipient id not found")
		return
	}

	newMsg := api.NewMessage{
		MessageId: msg.MessageId,
		Message:   msg.Message,
		Timestamp: msg.Timestamp.Local().Format("15:04"),
	}

	msgSent := api.MessageSent{
		MessageId: msg.MessageId,
		Timestamp: msg.Timestamp.Local().Format("15:04"),
	}

	senderMsg := newMsg
	senderMsg.IsMine = true

	recipientMsg := newMsg
	recipientMsg.IsMine = false

	h.persistPrivateMsg <- msg

	err := recipient.sendEvent("new_message", recipientMsg)
	if err != nil {
		log.Println(err)
		return
	}

	err = sender.sendEvent("new_message", senderMsg)
	if err != nil {
		log.Println(err)
		return
	}

	err = sender.sendEvent("message_sent", msgSent)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Hub) handleSendPublicMsg(msg model.PublicMessage) {
	sender, ok := h.Clients[msg.SenderId]
	if !ok {
		log.Println("Sender id not found")
		return
	}

	newMsg := api.NewPublicMessage{
		MessageId:   msg.MessageId,
		SenderId:    msg.SenderId,
		SenderName:  sender.Name,
		SenderEmoji: sender.Emoji,
		Message:     msg.Message,
		Timestamp:   msg.Timestamp.Local().Format("15:04"),
	}

	senderMsg := newMsg
	senderMsg.IsMine = true

	otherMsg := newMsg
	otherMsg.IsMine = false

	h.persistPublicMsg <- msg

	err := sender.sendEvent("new_public_message", senderMsg)
	if err != nil {
		log.Println(err)
		return
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
}

func (h *Hub) handleSendGroupMsg(msg model.GroupMessage) {
	sender, ok := h.Clients[msg.SenderId]
	if !ok {
		log.Println("Sender id not foudn")
		return
	}

	group, ok := h.groups[msg.GroupId]
	if !ok {
		log.Println("Group id not found")
		return
	}

	newMsg := api.NewPublicMessage{
		MessageId:   msg.MessageId,
		SenderId:    msg.SenderId,
		SenderName:  sender.Name,
		SenderEmoji: sender.Emoji,
		Message:     msg.Message,
		Timestamp:   msg.Timestamp.Local().Format("15:04"),
	}

	senderMsg := newMsg
	senderMsg.IsMine = true

	otherMsg := newMsg
	otherMsg.IsMine = false

	h.persistGroupMsg <- msg

	err := sender.sendEvent("new_group_message", senderMsg)
	if err != nil {
		log.Println(err)
		return
	}

	for memberId := range group.memberIds {
		if memberId == sender.UserId {
			continue
		}
		client, ok := h.Clients[memberId]
		if !ok {
			log.Printf("client %v not found in group %v\n", memberId, group.name)
			continue
		}
		err := client.sendEvent("new_group_message", otherMsg)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
