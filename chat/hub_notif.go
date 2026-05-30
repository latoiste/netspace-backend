package chat

import (
	"log"

	"github.com/latoiste/netspace/model"
)

func (h *Hub) sendNotification(receiver *Client, notif model.Notification) {
	notif.UserId = receiver.UserId

	err := receiver.sendEvent("new_notification", notif)
	if err != nil {
		log.Println(err)
	}
	h.persistNotification <- notif
}
