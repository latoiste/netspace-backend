package chat

import (
	"github.com/google/uuid"
	"github.com/latoiste/netspace/model"
)

type Group struct {
	id        string
	name      string
	hostId    string
	memberIds map[string]bool
	// pendingInvites are user ids that have been invited but haven't accepted or
	// rejected yet. Used to make invites idempotent: re-inviting someone who's
	// already a member or already pending is a no-op, so they never get two
	// invite notifications for the same group.
	pendingInvites map[string]bool
}

type GroupInvite struct {
	groupId string
	userId  string
}

// inviteRequest carries an invite-more-people action into the run() goroutine.
// inviterId is whoever clicked invite (must be an online member); the notif uses
// their name/emoji and doesn't depend on the group host being online.
type inviteRequest struct {
	groupId   string
	inviterId string
	userIds   []string
}

// renameRequest carries a group rename into the run() goroutine.
type renameRequest struct {
	groupId string
	userId  string
	name    string
}

// createGroupRequest carries a create-group action into the run() goroutine and
// carries the result back on reply, so the calling ReadPump never touches hub
// state directly (which would be a data race).
type createGroupRequest struct {
	name      string
	hostId    string
	memberIds []string
	reply     chan createGroupResult
}

type createGroupResult struct {
	group *Group
	err   error
}

// persistGroupReq persists a new group plus its host membership, in order
// (group row first, then the FK-referencing member row), without the persist
// loop having to read the run()-owned h.groups map.
type persistGroupReq struct {
	group  model.Group
	hostId string
}

func GenerateGroupId() string {
	id := uuid.New()
	return id.String()
}
