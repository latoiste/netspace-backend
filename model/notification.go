package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationType int

const (
	MessageNotif NotificationType = iota
	GroupInviteNotif
)

type Notification struct {
	Id             string
	UserId         string
	Type           string
	Emoji          string
	AvatarGradient string
	Title          string
	Timestamp      time.Time
	Description    string
	Unread         bool
	PrimaryLabel   string
	SecondaryLabel string
}

func newBaseNotif(
	emoji string,
	timestamp time.Time,
) Notification {
	return Notification{
		Id:        GenerateNotifId(),
		Emoji:     emoji,
		Unread:    true,
		Timestamp: timestamp,
	}
}

func NewGroupInviteNotif(emoji string, timestamp time.Time, senderName string, groupName string) Notification {
	baseNotif := newBaseNotif(emoji, timestamp)

	return Notification{
		Id:             baseNotif.Id,
		Emoji:          baseNotif.Emoji,
		Timestamp:      baseNotif.Timestamp,
		Unread:         baseNotif.Unread,
		Type:           GroupInviteNotif.String(),
		AvatarGradient: "linear-gradient(135deg, #6366f1, #8b5cf6)",
		Title:          "Undangan Group Session",
		Description:    fmt.Sprintf("%v mengundangmu ke \"%v\"", senderName, groupName),
		PrimaryLabel:   "Gabung",
		SecondaryLabel: "Tolak",
	}
}

func NewMessageNotif(emoji string, timestamp time.Time, senderName string, message string) Notification {
	baseNotif := newBaseNotif(emoji, timestamp)

	return Notification{
		Id:             baseNotif.Id,
		Emoji:          baseNotif.Emoji,
		Timestamp:      baseNotif.Timestamp,
		Unread:         baseNotif.Unread,
		Type:           MessageNotif.String(),
		AvatarGradient: "linear-gradient(135deg, #f97316, #ec4899)",
		Title:          fmt.Sprintf("%v mengirim pesan", senderName),
		Description:    fmt.Sprintf("\"%v\"", message),
	}
}

func (nt NotificationType) String() string {
	strings := [...]string{"message", "group_invite", "chat_request", "system"}

	return strings[nt]
}

func GenerateNotifId() string {
	id := uuid.New()
	return id.String()
}
