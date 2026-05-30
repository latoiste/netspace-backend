package chat

import (
	"context"
	"log"
	"time"
)

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

func (h *Hub) persistGroupMsgLoop() {
	for msg := range h.persistGroupMsg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.InsertGroupMessage(msg, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}
}

func (h *Hub) persistGroupLoop() {
	for group := range h.persistGroup {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.InsertGroup(group, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}
}

func (h *Hub) persistDissolveGroupLoop() {
	for group := range h.persistDissolveGroup {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.UpdateGroupIsActive(group.Id, group.IsActive, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}
}

func (h *Hub) persistNotificationLoop() {
	for notif := range h.persistNotification {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.InsertNotification(notif, notif.UserId, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}
}

func (h *Hub) persistNotificationReadLoop() {
	for notif := range h.persistNotificationRead {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		err := h.repo.UpdateNotificationUnread(notif.NotificationId, false, ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}
}
