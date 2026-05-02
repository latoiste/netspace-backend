package model

type Room struct {
	id      string
	clients []*Client
	history []*Message
}

type RoomRequest struct {
	client *Client
	room   *Room
}
