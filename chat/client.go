package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

type Client struct {
	Hub          *Hub
	conn         *websocket.Conn
	Send         chan []byte
	UserId       string
	Name         string
	Emoji        string
	LocationSlug string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId string, name string, emoji string, locationSlug string) *Client {
	return &Client{
		Hub:          hub,
		conn:         conn,
		Send:         make(chan []byte, 512),
		UserId:       userId,
		Name:         name,
		Emoji:        emoji,
		LocationSlug: locationSlug,
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
				log.Println("On send_message", err)
				continue
			}

			c.Hub.sendPrivateMsg <- model.PrivateMessage{
				BaseMessage: model.BaseMessage{
					MessageId:  model.GenerateMessageId(),
					LocationId: c.Hub.locationId,
					SenderId:   c.UserId,
					Message:    data.Message,
					Timestamp:  time.Now().UTC(),
				},
				RecipientId: data.RecipientId,
			}
		case "typing_start":
			msg := c.handleTypingRequest(req)
			if msg != (api.TypingEvent{}) {
				c.Hub.typingStart <- msg
			}
		case "typing_stop":
			msg := c.handleTypingRequest(req)
			if msg != (api.TypingEvent{}) {
				c.Hub.typingStop <- msg
			}
		case "send_public_message":
			var data api.SendPublicMessage
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				log.Println("On send_public_message", err)
				continue
			}

			c.Hub.sendPublicMsg <- model.PublicMessage{
				BaseMessage: model.BaseMessage{
					MessageId:  model.GenerateMessageId(),
					LocationId: c.Hub.locationId,
					SenderId:   c.UserId,
					Message:    data.Message,
					Timestamp:  time.Now().UTC(),
				},
			}
		case "public_typing_start":
			msg := c.handlePublicTypingRequest(req)
			if msg != (api.PublicUserTyping{}) {
				c.Hub.publicTypingStart <- msg
			}
		case "public_typing_stop":
			msg := c.handlePublicTypingRequest(req)
			if msg != (api.PublicUserTyping{}) {
				c.Hub.publicTypingStop <- msg
			}
		case "create_group":
			var data api.CreateGroup
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				log.Println("On create_group", err)
				continue
			}
			// add group to hub, check name
			group, err := c.Hub.createGroup(data.Name, c.UserId, data.MemberIds)
			if err != nil {
				log.Println(err)
				continue
			}
			// send back group_created
			groupCreated := api.GroupCreated{
				GroupId: group.id,
				Name:    group.name,
			}
			c.sendEvent("group_created", groupCreated)
		case "invite_to_group":
			var data api.InviteToGroup
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				log.Println(err)
				continue
			}

			c.Hub.inviteToGroup <- data
		case "accept_group_invite":
			var data api.GroupInviteResponse
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				log.Println(err)
				continue
			}

			c.Hub.acceptInvite <- GroupInvite{
				groupId: data.GroupId,
				userId:  c.UserId,
			}
		// case "reject_group_invite":
		// 	var data api.GroupInviteResponse
		case "leave_group":
			var data api.LeaveGroup
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				log.Println(err)
				continue
			}

			c.Hub.leaveGroup <- GroupInvite{
				groupId: data.GroupId,
				userId:  c.UserId,
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

func (c *Client) handleTypingRequest(req api.WsEvent) api.TypingEvent {
	var data api.TypingRequest
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		log.Println(err)
		return api.TypingEvent{}
	}

	return api.TypingEvent{
		SenderId:    c.UserId,
		RecipientId: data.RecipientId,
	}
}

func (c *Client) handlePublicTypingRequest(req api.WsEvent) api.PublicUserTyping {
	var data api.PublicTypingRequest

	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		log.Println(err)
		return api.PublicUserTyping{}
	}

	return api.PublicUserTyping{
		UserId: c.UserId,
		Name:   c.Name,
		Emoji:  c.Emoji,
	}
}
