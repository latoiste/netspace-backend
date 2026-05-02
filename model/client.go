package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	Username string
	Age      int
	Interest []string
	RoomId   string
}

func NewClient(hub *Hub, conn *websocket.Conn, username string, age int, interest []string) *Client {
	return &Client{
		Hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		Username: username,
		Age:      age,
		Interest: interest,
		RoomId:   "",
	}
}

// buat ngehandle request dari frontend
func (c *Client) ReadPump() {
	defer c.conn.Close()

	for {
		// di sini jadiin timeoutnya infinite dulu,
		// nanti pake ping pongnya buat ngecek kalo client masih connect ato ngga
		c.conn.SetReadDeadline(time.Time{})

		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		if msgType != 0x1 {
			log.Println("MEssage is not json alsdlkasdnsd")
			continue
		}

		var req Request
		err = json.Unmarshal(msg, &req)
		if err != nil {
			log.Println(err)
			continue
		}

		switch req.Type {
		case "join":
		case "leave":
		case "logout":
		case "chat":
		}
	}
}

func (c *Client) WritePump() {

}
