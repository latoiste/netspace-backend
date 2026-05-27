package api

import (
	"encoding/json"
	"time"
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
	MessageId string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	IsMine    bool      `json:"isMine"`
}

type NewPublicMessage struct {
	MessageId   string    `json:"id"`
	SenderId    string    `json:"senderId"`
	SenderName  string    `json:"senderName"`
	SenderEmoji string    `json:"senderEmoji"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	IsMine      bool      `json:"isMine"`
}

type MessageSent struct {
	MessageId string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
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
