package chat

import (
	"encoding/json"
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
		case "send_message":
			var data api.SendMessage
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				log.Println(err)
				continue
			}

			c.Hub.sendPrivateMsg <- PrivateMessage{
				MessageId:   generateMessageId(),
				SenderId:    c.UserId,
				RecipientId: data.RecipientId,
				Message:     data.Message,
				Timestamp:   time.Now(),
			}
		case "typing_start":
			msg := c.handleTypingRequest(req)
			if msg != (TypingEvent{}) {
				c.Hub.typingStart <- msg
			}
		case "typing_stop":
			msg := c.handleTypingRequest(req)
			if msg != (TypingEvent{}) {
				c.Hub.typingStop <- msg
			}
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

func (c *Client) sendEvent(event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	wsEvent := api.WsEvent{
		Event: event,
		Data:  payload,
	}

	message, err := json.Marshal(wsEvent)
	if err != nil {
		return err
	}

	c.Send <- message

	return nil
}

func (c *Client) handleTypingRequest(req api.WsEvent) TypingEvent {
	var data api.TypingRequest
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		log.Println(err)
		return TypingEvent{}
	}

	return TypingEvent{
		senderId:    c.UserId,
		recipientId: data.RecipientId,
	}
}
