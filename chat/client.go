package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/api"
)

type Client struct {
	Hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	Username string
	Age      int
	Interest []string
	RoomId   string
	// logout   chan bool
}

func NewClient(hub *Hub, conn *websocket.Conn, username string, age int, interest []string) *Client {
	return &Client{
		Hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 512),
		Username: username,
		Age:      age,
		Interest: interest,
		RoomId:   "",
		// logout:   make(chan bool),
	}
}

// buat ngehandle request dari frontend
func (c *Client) ReadPump() {
	defer func() {
		log.Println("Connection is clsoed")
		c.conn.Close()
	}()

	for {
		// di sini jadiin timeoutnya infinite dulu,
		// nanti pake ping pongnya buat ngecek kalo client masih connect ato ngga
		c.conn.SetReadDeadline(time.Time{})

		var req api.WsRequest
		err := c.conn.ReadJSON(&req)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Connection is closed normally")
				break
			}
			log.Println("Connection is already closed", err)
			break
		}

		switch req.Type {
		case "join":
		case "leave":
		case "logout":
			c.Hub.Logout <- c
			close(c.send)
			// <-c.logout
			return
		case "chat":
		}
	}
}

func (c *Client) WritePump() {

}
