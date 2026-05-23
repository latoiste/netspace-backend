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

type MessageSent struct {
	MessageId string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// ======================

// Client -> Server
type SendMessage struct {
	RecipientId string `json:"recipientId"`
	Message     string `json:"message"`
}
