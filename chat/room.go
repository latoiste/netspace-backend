package chat

import "github.com/latoiste/netspace/model"

type Room struct {
	id      string
	clients []*Client
	history []*model.Message
}

type RoomRequest struct {
	client *Client
	room   *Room
}
