package chat

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

// createGroup builds a new group, invites the initial members, and persists it.
// Runs ONLY on the run() goroutine (driven by createGroupReq), so it can touch
// h.groups / h.Clients without a data race.
func (h *Hub) createGroup(groupName string, hostId string, memberIds []string) (*Group, error) {
	host, ok := h.Clients[hostId]
	if !ok {
		return nil, errors.New("Invalid host id")
	}

	name := strings.TrimSpace(groupName)
	if name == "" {
		return nil, errors.New("Group name is empty")
	}

	for _, group := range h.groups {
		if group.name == name {
			return nil, errors.New("Group name is taken")
		}
	}

	groupId := GenerateGroupId()

	group := &Group{
		id:             groupId,
		name:           name,
		hostId:         hostId,
		memberIds:      make(map[string]bool),
		pendingInvites: make(map[string]bool),
	}
	group.memberIds[hostId] = true
	h.groups[groupId] = group

	// The creator is the inviter for the initial invites.
	for _, memberId := range memberIds {
		h.inviteMember(group, host, memberId)
	}

	h.persistGroup <- persistGroupReq{
		group:  model.Group{Id: groupId, Name: name, IsActive: true},
		hostId: hostId,
	}

	return group, nil
}

// inviteMember sends a single group invite — idempotently. It skips the inviter,
// existing members, and anyone already invited (so nobody ever gets two invite
// notifications for the same group), records the invite as pending, and only
// delivers to a user who is currently online. Each invite gets its own
// notification id (so per-user persistence doesn't collide).
func (h *Hub) inviteMember(group *Group, inviter *Client, memberId string) {
	if memberId == "" || memberId == inviter.UserId {
		return
	}
	if group.memberIds[memberId] || group.pendingInvites[memberId] {
		return
	}

	group.pendingInvites[memberId] = true

	client, ok := h.Clients[memberId]
	if !ok {
		log.Printf("invite: user %v not online, skipped\n", memberId)
		return
	}

	notif := model.NewGroupInviteNotif(
		inviter.Emoji,
		time.Now().UTC(),
		inviter.Name,
		group.name,
		group.id,
	)
	h.sendNotification(client, notif)
}

func (h *Hub) handleLeaveGroup(invite GroupInvite) {
	sender, ok := h.Clients[invite.userId]
	if !ok {
		log.Println("Sender not found")
		return
	}

	group, ok := h.groups[invite.groupId]
	if !ok {
		log.Println("Group not found")
		return
	}

	_, ok = group.memberIds[sender.UserId]
	if !ok {
		log.Printf("User %v is trying to leave but not in group %v\n", sender.UserId, group.name)
		return
	}

	// Drop the leaver's membership row so the group disappears from THEIR chat
	// list, while staying intact for everyone else. Done off the run() goroutine.
	go func(gid, uid string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := h.repo.RemoveGroupMember(gid, uid, ctx); err != nil {
			log.Println("remove group member:", err)
		}
	}(group.id, sender.UserId)

	h.broadcastGroupMemberLeft(group, sender)
}

// handleInviteToGroup is the invite-more-people path (the invite sheet). The
// inviter (whoever clicked invite) must be an online member; invites use the
// inviter's name/emoji and don't depend on the host being online.
func (h *Hub) handleInviteToGroup(req inviteRequest) {
	group, ok := h.groups[req.groupId]
	if !ok {
		log.Println("Group not found")
		return
	}

	inviter, ok := h.Clients[req.inviterId]
	if !ok {
		log.Println("Inviter not found")
		return
	}

	if !group.memberIds[req.inviterId] {
		log.Printf("Non-member %v tried to invite to group %v\n", req.inviterId, group.name)
		return
	}

	for _, memberId := range req.userIds {
		h.inviteMember(group, inviter, memberId)
	}
}

func (h *Hub) handleAcceptInvite(invite GroupInvite) {
	sender, ok := h.Clients[invite.userId]
	if !ok {
		log.Println("Sender not found")
		return
	}

	group, ok := h.groups[invite.groupId]
	if !ok {
		log.Println("Group not found")
		return
	}

	// Either way, the invite is no longer pending.
	delete(group.pendingInvites, sender.UserId)
	h.dismissGroupInviteNotif(sender.UserId, group.id)

	_, ok = group.memberIds[sender.UserId]
	if ok {
		log.Printf("User %v is trying to join but already in group %v\n", sender.UserId, group.name)
		return
	}

	newMember := api.MemberJoined{
		Id:     sender.UserId,
		Name:   sender.Name,
		Emoji:  sender.Emoji,
		IsHost: false,
	}

	for memberId := range group.memberIds {
		client, ok := h.Clients[memberId]
		if !ok {
			log.Printf("Sender id %v not found in group %v\n", memberId, group.name)
			continue
		}
		client.sendEvent("member_joined", newMember)
	}
	group.memberIds[sender.UserId] = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := h.repo.InsertGroupMember(group.id, newMember.Id, ctx)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Member is added to group")
}

// handleRejectInvite notifies the group host that an invitee declined their
// invite. The decliner's own notification is dismissed client-side.
func (h *Hub) handleRejectInvite(invite GroupInvite) {
	decliner, ok := h.Clients[invite.userId]
	if !ok {
		log.Println("Decliner not found")
		return
	}

	group, ok := h.groups[invite.groupId]
	if !ok {
		log.Println("Group not found")
		return
	}

	// No longer a pending invite.
	delete(group.pendingInvites, decliner.UserId)
	h.dismissGroupInviteNotif(decliner.UserId, group.id)

	host, ok := h.Clients[group.hostId]
	if !ok {
		log.Println("Host not found")
		return
	}

	// No point telling the host they declined their own invite.
	if host.UserId == decliner.UserId {
		return
	}

	rejectNotif := model.NewGroupRejectNotif(
		decliner.Emoji,
		time.Now().UTC(),
		decliner.Name,
		group.name,
	)
	h.sendNotification(host, rejectNotif)
}

// handleRenameGroup renames a group (any member may rename it, WhatsApp-style),
// persists the new name, and tells every member live so open views + chat lists
// update. Runs on the run() goroutine.
func (h *Hub) handleRenameGroup(req renameRequest) {
	group, ok := h.groups[req.groupId]
	if !ok {
		log.Println("Group not found")
		return
	}

	if !group.memberIds[req.userId] {
		log.Printf("Non-member %v tried to rename group %v\n", req.userId, group.name)
		return
	}

	name := strings.TrimSpace(req.name)
	if name == "" || name == group.name {
		return
	}

	group.name = name

	go func(gid, n string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := h.repo.UpdateGroupName(gid, n, ctx); err != nil {
			log.Println("rename group:", err)
		}
	}(group.id, name)

	renamed := api.GroupRenamed{GroupId: group.id, Name: name}
	for memberId := range group.memberIds {
		if client, ok := h.Clients[memberId]; ok {
			client.sendEvent("group_renamed", renamed)
		}
	}
}

// dismissGroupInviteNotif deletes the user's invite notification for a group
// once they've resolved it (accepted/rejected), so it doesn't reappear on a
// notifications re-fetch. Runs off the run() goroutine.
func (h *Hub) dismissGroupInviteNotif(userId string, groupId string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := h.repo.DeleteGroupInviteNotif(userId, groupId, ctx); err != nil {
			log.Println("dismiss group invite notif:", err)
		}
	}()
}

func (h *Hub) broadcastGroupMemberLeft(group *Group, leavingClient *Client) {
	memberLeft := api.MemberLeft{
		UserId: leavingClient.UserId,
		Name:   leavingClient.Name,
	}

	for memberId := range group.memberIds {
		if memberId == leavingClient.UserId {
			continue
		}
		client, ok := h.Clients[memberId]
		if !ok {
			log.Printf("Sender id not found")
			continue
		}
		client.sendEvent("member_left", memberLeft)
	}

	delete(group.memberIds, leavingClient.UserId)
	delete(group.pendingInvites, leavingClient.UserId)

	if len(group.memberIds) <= 1 {
		for memberId := range group.memberIds {
			client, ok := h.Clients[memberId]
			if !ok {
				continue
			}
			groupDissolved := api.GroupDissolved{
				GroupId: group.id,
			}

			client.sendEvent("group_dissolved", groupDissolved)
			delete(group.memberIds, memberId)
		}
		delete(h.groups, group.id)

		h.persistDissolveGroup <- model.Group{
			Id:       group.id,
			Name:     group.name,
			IsActive: false,
		}
	}
}
