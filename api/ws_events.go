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

type GroupInviteResponse struct {
	GroupId string `json:"groupId"`
}

type NotificationRead struct {
	NotificationId string `json:"notificationId"`
}
