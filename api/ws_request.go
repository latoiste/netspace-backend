package api

import (
	"encoding/json"
)

type WsRequest struct {
	// typenya bisa: "register", "logout", "join", "leave", "chat"
	Type string `json:"type"`

	// ini implementasi request dari setiap request type
	Body json.RawMessage `json:"body"`
}

// untuk register/logout
type RegisterBody struct {
	Username string   `json:"username"`
	Id       string   `json:"id"`
	Age      int      `json:"age"`
	Interest []string `json:"interest"`
}

// type LogoutBody struct {
// 	Id string `json:"id"`
// }

// untuk join/leave
// type RoomBody struct {
// 	Username string `json:"username"`
// 	RoomId   string `json:"roomId"`
// }

// type ChatBody struct {
// 	Username string `json:"username"`
// 	RoomId   string `json:"roomId"`
// 	Msg      string `json:"msg"`
// }
