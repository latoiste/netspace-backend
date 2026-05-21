package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/api"
)

type Client struct {
	Hub          *Hub
	conn         *websocket.Conn
	Send         chan []byte
	UserId       string
	LocationSlug string
	RoomId       string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId string, locationSlug string) *Client {
	return &Client{
		Hub:          hub,
		conn:         conn,
		Send:         make(chan []byte, 512),
		UserId:       userId,
		LocationSlug: locationSlug,
		RoomId:       "",
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

		var req api.WsEvent
		err := c.conn.ReadJSON(&req)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Connection is closed normally")
				break
			}
			log.Println("Connection is already closed", err)
			break
		}

		switch req.Event {

		}
	}
}

func (c *Client) WritePump() {

}
