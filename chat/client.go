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
	Name         string
	LocationSlug string
	RoomId       string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId string, name string, locationSlug string) *Client {
	return &Client{
		Hub:          hub,
		conn:         conn,
		Send:         make(chan []byte, 512),
		UserId:       userId,
		Name:         name,
		LocationSlug: locationSlug,
		RoomId:       "",
	}
}

func (c *Client) ReadPump() {
	defer func() {
		log.Println("Connection is clsoed")
		c.Hub.unregister <- c
		c.conn.Close()
	}()

	for {
		// TODO: di sini jadiin timeoutnya infinite dulu,
		// nanti pake ping pongnya buat ngecek kalo client masih connect ato ngga
		c.conn.SetReadDeadline(time.Time{})

		var req api.WsEvent
		err := c.conn.ReadJSON(&req)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Connection is closed normally")
				break
			}
			log.Println(err)
			break
		}

		switch req.Event {

		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		msg, ok := <-c.Send
		if !ok {
			log.Println("Connection is closed by Hub")
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("written to client")
	}
}
