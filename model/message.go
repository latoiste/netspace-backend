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
}

type PublicMessage struct {
	BaseMessage
}

type GroupMessage struct {
	BaseMessage
	GroupId string
}

func GenerateMessageId() string {
	id := uuid.New()
	return id.String()
}
