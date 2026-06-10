package api

import (
	"encoding/json"

	"github.com/latoiste/netspace/model"
)

type WsEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// Server -> Client
type UserJoined struct {
	User UserDTO `json:"user"`
}

type UserLeft struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

type NewMessage struct {
	MessageId string `json:"id"`
	SenderId  string `json:"senderId"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	IsMine    bool   `json:"isMine"`
}

type NewPublicMessage struct {
	MessageId   string `json:"id"`
	SenderId    string `json:"senderId"`
	SenderName  string `json:"senderName"`
	SenderEmoji string `json:"senderEmoji"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
	IsMine      bool   `json:"isMine"`
	IsAdmin     bool   `json:"isAdmin"`
}

type PublicMessagesCleared struct {
	LocationSlug string `json:"locationSlug"`
	DeletedCount int64  `json:"deletedCount"`
}

type MessageSent struct {
	MessageId string `json:"id"`
	Timestamp string `json:"timestamp"`
}

type UserTyping struct {
	UserId string `json:"userId"`
}

type PublicUserTyping struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Emoji  string `json:"emoji"`
}

type TypingEvent struct {
	SenderId    string `json:"senderId"`
	RecipientId string `json:"recipientId"`
}

type NewGroupMessage struct {
	MessageId   string `json:"id"`
	GroupId     string `json:"groupId"`
	SenderId    string `json:"senderId"`
	SenderName  string `json:"senderName"`
	SenderEmoji string `json:"senderEmoji"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
	IsMine      bool   `json:"isMine"`
}

type GroupInvite struct {
	model.Notification
}

type GroupCreated struct {
	GroupId string `json:"groupId"`
	Name    string `json:"name"`
}

type GroupCreatedFailed struct {
	Message string
}

type MemberJoined struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Emoji  string `json:"emoji"`
	IsHost bool   `json:"isHost"`
}

type MemberLeft struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

type GroupDissolved struct {
	GroupId string `json:"groupId"`
}

// GroupRenamed tells every member a group was renamed so their open views and
// chat list update live.
type GroupRenamed struct {
	GroupId string `json:"groupId"`
	Name    string `json:"name"`
}

// ForceLogout tells a client an admin has ended its session. The frontend
// reacts by tearing the session down (blacklisting its token, stopping
// reconnection, and returning to the entry screen) instead of silently
// reconnecting — which would put the user straight back online.
type ForceLogout struct {
	Reason string `json:"reason"`
}

// MessagesRead notifies the sender that ReaderId has opened their chat and read
// the sender's messages, so the sender's bubbles flip to blue double-checks.
type MessagesRead struct {
	ReaderId string `json:"readerId"`
}

// ======================

// Client -> Server
type SendMessage struct {
	RecipientId string `json:"recipientId"`
	Message     string `json:"message"`
}

type TypingRequest struct {
	RecipientId string `json:"recipientId"`
}

type SendPublicMessage struct {
	LocationSlug string `json:"locationSlug"`
	Message      string `json:"message"`
}

type PublicTypingRequest struct {
	LocationSlug string `json:"locationSlug"`
}

type CreateGroup struct {
	Name      string   `json:"name"`
	MemberIds []string `json:"memberIds"`
}

type SendGroupMessage struct {
	GroupId string `json:"groupId"`
	Message string `json:"message"`
}

type LeaveGroup struct {
	GroupId string `json:"groupId"`
}

type InviteToGroup struct {
	GroupId string   `json:"groupId"`
	UserIds []string `json:"userIds"`
}

type RenameGroup struct {
	GroupId string `json:"groupId"`
	Name    string `json:"name"`
}

type GroupInviteResponse struct {
	GroupId string `json:"groupId"`
}

type BlockUser struct {
	UserId string `json:"userId"`
}

type UnblockUser struct {
	UserId string `json:"userId"`
}

// MarkRead is sent when the user opens a DM: "I've read the messages SenderId
// sent me." The hub flips those messages to read and tells SenderId.
type MarkRead struct {
	SenderId string `json:"senderId"`
}

type NotificationRead struct {
	NotificationId string `json:"notificationId"`
}

type NotificationDismiss struct {
	NotificationId string `json:"notificationId"`
}
