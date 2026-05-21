package api

import "encoding/json"

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

// ======================

// Client -> Server
