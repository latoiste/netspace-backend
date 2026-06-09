package chat

import (
	"context"
	"log"
	"time"

	"github.com/latoiste/netspace/api"
)

// markReadRequest carries a DM read receipt from a client's ReadPump into the
// hub goroutine: readerId opened the chat with senderId and read their messages.
type markReadRequest struct {
	readerId string
	senderId string
}

// handleMarkRead flips the sender's messages to "read" in the DB and tells the
// sender (if they're online) so their bubbles turn blue immediately. Runs on the
// run() goroutine, so the h.Clients lookup is race-free.
func (h *Hub) handleMarkRead(req markReadRequest) {
	if req.readerId == "" || req.senderId == "" || req.readerId == req.senderId {
		return
	}

	// Persist the flip off the run() goroutine so the DB timeout can't stall the
	// hub. The live notification below is what makes the receipt feel instant;
	// the DB write just keeps history correct after a re-fetch.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if _, err := h.repo.MarkPrivateMessagesRead(req.senderId, req.readerId, ctx); err != nil {
			log.Println("mark read:", err)
		}
	}()

	if sender, ok := h.Clients[req.senderId]; ok {
		if err := sender.sendEvent("messages_read", api.MessagesRead{ReaderId: req.readerId}); err != nil {
			log.Println("messages_read send:", err)
		}
	}
}
