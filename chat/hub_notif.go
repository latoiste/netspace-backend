package chat

import (
	"log"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Hub) sendNotification(receiver *Client, notif model.Notification) {
	notif.UserId = receiver.UserId

	notifDTO := api.ConstructNotificationDTO(notif)

	err := receiver.sendEvent("new_notification", notifDTO)
	if err != nil {
		log.Println(err)
	}
	h.persistNotification <- notif
}
