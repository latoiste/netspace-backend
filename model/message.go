package model

import (
	"time"

	"github.com/google/uuid"
)

type PrivateMessage struct {
	MessageId   string
	LocationId  int
	SenderId    string
	RecipientId string
	Message     string
	Timestamp   time.Time
}

type PublicMessage struct {
	MessageId  string
	LocationId int
	SenderId   string
	Message    string
	Timestamp  time.Time
}

type GroupMessage struct {
}

func GenerateMessageId() string {
	id := uuid.New()
	return id.String()
}
