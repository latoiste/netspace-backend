package model

import "time"

// chat message
type Message struct {
	UserId  string `json:"userId"`
	RoomId  string `json:"roomId"`
	Content string `json:"content"`

	// field ini ga ush di send dari client,
	// tapi client akan menerima field ini di result
	Time time.Time `json:"time"`
}
