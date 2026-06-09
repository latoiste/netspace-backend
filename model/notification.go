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
	// GroupId links a group_invite notification to its group so the recipient
	// can accept it. Persisted, so it survives a re-fetch from /api/notifications;
	// empty for non-group notifs.
	GroupId string
	// SenderId is the user id of whoever triggered the notification (e.g. the
	// person who messaged you). Lets the client open the DM with them when the
	// notification is tapped. Persisted; empty for notifs with no actor.
	SenderId string
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

func NewGroupInviteNotif(emoji string, timestamp time.Time, senderName string, groupName string, groupId string) Notification {
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
		GroupId:        groupId,
	}
}

// NewGroupRejectNotif tells the host that an invitee declined their group
// invite, so the rejection surfaces in the host's notification list.
func NewGroupRejectNotif(emoji string, timestamp time.Time, declinerName string, groupName string) Notification {
	baseNotif := newBaseNotif(emoji, timestamp)

	return Notification{
		Id:             baseNotif.Id,
		Emoji:          baseNotif.Emoji,
		Timestamp:      baseNotif.Timestamp,
		Unread:         baseNotif.Unread,
		Type:           SystemNotif.String(),
		AvatarGradient: "linear-gradient(135deg, #64748b, #94a3b8)",
		Title:          "Undangan ditolak",
		Description:    fmt.Sprintf("%v menolak undangan ke \"%v\"", declinerName, groupName),
	}
}

func NewMessageNotif(emoji string, timestamp time.Time, senderName string, message string, senderId string) Notification {
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
		SenderId:       senderId,
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
