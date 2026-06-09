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
	for req := range h.persistGroup {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		// Insert the group row first, then the host's membership (which has an FK
		// to it). hostId comes in on the request, so this loop never reads the
		// run()-owned h.groups map.
		if err := h.repo.InsertGroup(req.group, ctx); err != nil {
			log.Println(err)
			cancel()
			continue
		}

		if err := h.repo.InsertGroupMember(req.group.Id, req.hostId, ctx); err != nil {
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

// notificationDismissReq is a "delete my notification" action: the notification
// id plus the owner (so a user can only dismiss their own).
type notificationDismissReq struct {
	notificationId string
	userId         string
}

func (h *Hub) persistNotificationDismissLoop() {
	for req := range h.notificationDismiss {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

		if err := h.repo.DeleteNotification(req.notificationId, req.userId, ctx); err != nil {
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
