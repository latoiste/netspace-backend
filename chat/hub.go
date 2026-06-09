package chat

import (
	"encoding/json"
	"log"
	"sync/atomic"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/db"
	"github.com/latoiste/netspace/model"
)

type Hub struct {
	Clients    map[string]*Client
	repo       *db.Repository
	groups     map[string]*Group
	Register   chan *Client
	unregister chan *Client
	locationId int

	// onlineCount mirrors len(Clients) — the number of users with a live socket
	// right now. It's written only from run() (the map's single owner) and read
	// lock-free by admin analytics, so "User Aktif Sekarang" reflects actual
	// connections instead of the DB isactive flag, which drifts (e.g. a ghost is
	// left behind when the process is killed before its disconnect handler runs).
	onlineCount atomic.Int64

	broadcast      chan api.WsEvent
	sendPrivateMsg chan model.PrivateMessage
	typingStart    chan api.TypingEvent
	typingStop     chan api.TypingEvent

	sendPublicMsg     chan model.PublicMessage
	publicTypingStart chan api.PublicUserTyping
	publicTypingStop  chan api.PublicUserTyping

	// createGroupReq carries create-group actions onto run() (with a reply
	// channel) so group/Client state is only ever mutated there — never from a
	// client's ReadPump goroutine, which would race.
	createGroupReq chan createGroupRequest
	inviteToGroup  chan inviteRequest
	renameGroup    chan renameRequest
	acceptInvite   chan GroupInvite
	rejectInvite   chan GroupInvite
	leaveGroup     chan GroupInvite
	sendGroupMsg   chan model.GroupMessage

	// forceLogout carries a user id whose live socket should be dropped (admin
	// kick). Handled on run() so the Clients map is only ever touched there.
	forceLogout chan string

	// blocks records session-scoped blocks: blocks[blockerId][blockedId] = true
	// means blockerId will not be delivered blockedId's messages/notifications.
	// Mutated and read only from the run() goroutine to stay race-free.
	blocks      map[string]map[string]bool
	blockUser   chan blockRequest
	unblockUser chan blockRequest

	// markRead carries a DM read receipt: readerId opened the chat with senderId
	// and read their messages. Handled on run() so the sender lookup is race-free.
	markRead chan markReadRequest

	notificationRead chan api.NotificationRead

	// notificationDismiss carries a "user tapped/dismissed a notification" action
	// (its id + the owner) so it's deleted for good. Buffered + handled off the
	// run() goroutine since it only touches the DB.
	notificationDismiss chan notificationDismissReq

	persistPrivateMsg       chan model.PrivateMessage
	persistPublicMsg        chan model.PublicMessage
	persistGroupMsg         chan model.GroupMessage
	persistGroup            chan persistGroupReq
	persistDissolveGroup    chan model.Group
	persistNotification     chan model.Notification
	persistNotificationRead chan api.NotificationRead
}

func NewHub(repo *db.Repository, locationId int) *Hub {
	return &Hub{
		Clients:                 make(map[string]*Client),
		locationId:              locationId,
		repo:                    repo,
		groups:                  make(map[string]*Group),
		Register:                make(chan *Client),
		unregister:              make(chan *Client),
		broadcast:               make(chan api.WsEvent),
		sendPrivateMsg:          make(chan model.PrivateMessage),
		typingStart:             make(chan api.TypingEvent),
		typingStop:              make(chan api.TypingEvent),
		sendPublicMsg:           make(chan model.PublicMessage),
		publicTypingStart:       make(chan api.PublicUserTyping),
		publicTypingStop:        make(chan api.PublicUserTyping),
		createGroupReq:          make(chan createGroupRequest),
		inviteToGroup:           make(chan inviteRequest),
		renameGroup:             make(chan renameRequest),
		acceptInvite:            make(chan GroupInvite),
		rejectInvite:            make(chan GroupInvite),
		leaveGroup:              make(chan GroupInvite),
		sendGroupMsg:            make(chan model.GroupMessage),
		forceLogout:             make(chan string, 8),
		blocks:                  make(map[string]map[string]bool),
		blockUser:               make(chan blockRequest),
		unblockUser:             make(chan blockRequest),
		markRead:                make(chan markReadRequest, 20),
		notificationRead:        make(chan api.NotificationRead, 20),
		notificationDismiss:     make(chan notificationDismissReq, 20),
		persistPrivateMsg:       make(chan model.PrivateMessage, 20),
		persistPublicMsg:        make(chan model.PublicMessage, 20),
		persistGroupMsg:         make(chan model.GroupMessage, 20),
		persistGroup:            make(chan persistGroupReq, 20),
		persistDissolveGroup:    make(chan model.Group, 20),
		persistNotification:     make(chan model.Notification, 20),
		persistNotificationRead: make(chan api.NotificationRead, 20),
	}
}

func (h *Hub) run() {
	go h.persistPrivateMsgLoop()
	go h.persistPublicMsgLoop()
	go h.persistGroupMsgLoop()
	go h.persistGroupLoop()
	go h.persistDissolveGroupLoop()
	go h.persistNotificationLoop()
	go h.persistNotificationReadLoop()
	go h.persistNotificationDismissLoop()

	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserId] = client
			h.onlineCount.Store(int64(len(h.Clients)))

		case client := <-h.unregister:
			err := h.removeClient(client)
			if err != nil {
				log.Println(err)
				continue
			}

		case event := <-h.broadcast:
			h.fanout(event)
		case msg := <-h.sendPrivateMsg:
			h.handleSendPrivateMsg(msg)

		case msg := <-h.typingStart:
			h.notifyTyping("user_typing", msg)

		case msg := <-h.typingStop:
			h.notifyTyping("user_stopped_typing", msg)

		case msg := <-h.sendPublicMsg:
			h.handleSendPublicMsg(msg)

		case msg := <-h.publicTypingStart:
			h.notifyPublicTyping("public_user_typing", msg)

		case msg := <-h.publicTypingStop:
			h.notifyPublicTyping("public_user_stopped_typing", msg)

		case req := <-h.createGroupReq:
			group, err := h.createGroup(req.name, req.hostId, req.memberIds)
			req.reply <- createGroupResult{group: group, err: err}

		case msg := <-h.inviteToGroup:
			h.handleInviteToGroup(msg)

		case req := <-h.renameGroup:
			h.handleRenameGroup(req)

		case invite := <-h.acceptInvite:
			h.handleAcceptInvite(invite)

		case invite := <-h.rejectInvite:
			h.handleRejectInvite(invite)

		case req := <-h.blockUser:
			h.handleBlockUser(req)

		case req := <-h.unblockUser:
			h.handleUnblockUser(req)

		case req := <-h.markRead:
			h.handleMarkRead(req)

		case invite := <-h.leaveGroup:
			h.handleLeaveGroup(invite)

		case msg := <-h.sendGroupMsg:
			h.handleSendGroupMsg(msg)

		case notif := <-h.notificationRead:
			h.persistNotificationRead <- notif

		case uid := <-h.forceLogout:
			if c, ok := h.Clients[uid]; ok {
				c.forceLogout()
			}
		}
	}
}

// fanout marshals an event once and pushes it to every connected client. It
// runs ONLY on the run() goroutine — the single owner of h.Clients — so there's
// no data race on the map and no send-on-closed-channel panic (the same
// goroutine is the only one that closes Send, via removeClient).
//
// Sends are non-blocking: a client whose 512-deep buffer is full is genuinely
// wedged (its WritePump isn't draining), so we skip it rather than stall the
// entire hub — which would freeze presence updates and message delivery for
// everyone. The wedged socket is reaped later by its own idle timeout.
func (h *Hub) fanout(event api.WsEvent) {
	message, err := json.Marshal(event)
	if err != nil {
		log.Println("fanout marshal:", err)
		return
	}
	for uid, client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			log.Printf("fanout: skipping backlogged client %s", uid)
		}
	}
}

// OnlineCount returns the number of users with a live socket on this hub right
// now. Read lock-free via the atomic mirror of len(Clients), so admin analytics
// can poll it without touching the run()-owned map.
func (h *Hub) OnlineCount() int {
	return int(h.onlineCount.Load())
}

func (h *Hub) AddClient(client *Client, user *model.User) error {
	userDTO := api.ConstructUserDTO(*user)
	userJoined := api.UserJoined{
		User: userDTO,
	}
	data, err := json.Marshal(userJoined)
	if err != nil {
		return err
	}

	h.Register <- client

	h.broadcast <- api.WsEvent{
		Event: "user_joined",
		Data:  data,
	}
	return nil
}

func (h *Hub) removeClient(client *Client) error {
	// A user can briefly hold two live sockets — a reconnect after a network
	// blip, or a second browser tab. Registration overwrites h.Clients[UserId]
	// with the newest socket, so an older socket's *delayed* unregister must not
	// evict the entry: doing so would drop the live client, broadcast a bogus
	// "user_left", and (via the DB flip below) mark the user offline — making
	// them vanish from everyone's roster while still connected. Only the socket
	// that is currently the live one is allowed to tear the presence down.
	if current, ok := h.Clients[client.UserId]; !ok || current != client {
		client.closeSend()
		return nil
	}

	delete(h.Clients, client.UserId)
	h.onlineCount.Store(int64(len(h.Clients)))
	client.closeSend()

	// Mark the user offline in the DB so a fresh roster fetch
	// (/api/locations/{slug}/users filters isActive = true) excludes them.
	// Done off the run() goroutine so the 2s DB timeout can't stall the hub.
	go client.handleClientDisconnect()

	userLeft := api.UserLeft{
		UserId: client.UserId,
		Name:   client.Name,
	}

	data, err := json.Marshal(userLeft)
	if err != nil {
		return err
	}

	for _, group := range h.groups {
		_, ok := group.memberIds[client.UserId]
		if !ok {
			continue
		}
		h.broadcastGroupMemberLeft(group, client)
	}

	// removeClient already runs on the run() goroutine, so fan out directly.
	// Sending to h.broadcast here would deadlock — run() is the only receiver
	// and it's busy executing this function.
	h.fanout(api.WsEvent{
		Event: "user_left",
		Data:  data,
	})
	return nil
}
