package chat

import "log"

// blockRequest carries a session-scoped block from a client's ReadPump into the
// hub goroutine, where the shared block map is safe to mutate.
type blockRequest struct {
	blockerId string
	blockedId string
}

// handleBlockUser records that blockerId no longer wants blockedId's messages.
// Runs in the hub goroutine, so writing h.blocks is race-free.
func (h *Hub) handleBlockUser(req blockRequest) {
	if req.blockerId == "" || req.blockedId == "" {
		return
	}
	if h.blocks[req.blockerId] == nil {
		h.blocks[req.blockerId] = make(map[string]bool)
	}
	h.blocks[req.blockerId][req.blockedId] = true
	log.Printf("User %v blocked %v\n", req.blockerId, req.blockedId)
}

// handleUnblockUser reverses a block: blockerId once again receives blockedId's
// messages/notifications. Runs in the hub goroutine, so writing h.blocks is
// race-free. No-op if the block wasn't there.
func (h *Hub) handleUnblockUser(req blockRequest) {
	if req.blockerId == "" || req.blockedId == "" {
		return
	}
	blocked, ok := h.blocks[req.blockerId]
	if !ok {
		return
	}
	delete(blocked, req.blockedId)
	if len(blocked) == 0 {
		delete(h.blocks, req.blockerId)
	}
	log.Printf("User %v unblocked %v\n", req.blockerId, req.blockedId)
}

// isBlocked reports whether blockerId has blocked blockedId. Read only from the
// hub goroutine (the message-delivery handlers).
func (h *Hub) isBlocked(blockerId string, blockedId string) bool {
	blocked, ok := h.blocks[blockerId]
	if !ok {
		return false
	}
	return blocked[blockedId]
}
