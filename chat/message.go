package chat

import (
	"time"

	"github.com/google/uuid"
)

// chat message
type PrivateMessage struct {
	MessageId   string    `json:"messageId"`
	LocationId  int       `json:"locationId"`
	SenderId    string    `json:"senderId"`
	RecipientId string    `json:"recipientId"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
}

func generateMessageId() string {
	id := uuid.New()
	return id.String()
}
