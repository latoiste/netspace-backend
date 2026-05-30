package chat

import (
	"log"

	"github.com/latoiste/netspace/api"
)

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
