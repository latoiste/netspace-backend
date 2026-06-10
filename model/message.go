package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseMessage struct {
	MessageId  string
	LocationId int
	SenderId   string
	Message    string
	Timestamp  time.Time
}

type PrivateMessage struct {
	BaseMessage
	RecipientId string
	// IsRead is the read receipt: true once the recipient has opened the chat
	// and seen this message. Only meaningful from the sender's perspective.
	IsRead bool
}

type PublicMessage struct {
	BaseMessage
	AdminId     string
	SenderName  string
	SenderEmoji string
	IsAdmin     bool
}

type GroupMessage struct {
	BaseMessage
	GroupId string
}

func GenerateMessageId() string {
	id := uuid.New()
	return id.String()
}
