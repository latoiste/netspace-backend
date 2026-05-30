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
	ChatRequestNotif
	SystemNotif
)

type Notification struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	Emoji          string `json:"emoji"`
	AvatarGradient string `json:"avatarGradient"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Timestamp      string `json:"timestamp"`
	Unread         bool   `json:"unread"`
	PrimaryLabel   string `json:"primaryLabel"`
	SecondaryLabel string `json:"secondaryLabel"`
}

func newBaseNotif(
	emoji string,
	timestamp time.Time,
) Notification {
	return Notification{
		Id:        GenerateNotifId(),
		Emoji:     emoji,
		Unread:    true,
		Timestamp: timestamp.Local().Format("15:04"),
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
