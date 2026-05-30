package chat

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Hub) createGroup(groupName string, hostId string, memberIds []string) (*Group, error) {
	host, ok := h.Clients[hostId]
	if !ok {
		return nil, errors.New("Invalid host id")
	}

	for _, group := range h.groups {
		if group.name == groupName {
			return nil, errors.New("Group name is taken")
		}
	}

	groupId := GenerateGroupId()

	group := &Group{
		id:        groupId,
		name:      groupName,
		hostId:    hostId,
		memberIds: make(map[string]bool),
	}
	group.memberIds[hostId] = true
	h.groups[groupId] = group

	groupInviteNotif := model.NewGroupInviteNotif(
		host.Emoji,
		time.Now().UTC(),
		host.Name,
		groupName,
	)

	for _, memberId := range memberIds {
		if client, ok := h.Clients[memberId]; ok {
			h.sendNotification(client, groupInviteNotif)
		}
	}

	h.persistGroup <- model.Group{
		Id:       groupId,
		Name:     groupName,
		IsActive: true,
	}

	return group, nil
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

	h.broadcastGroupMemberLeft(group, sender)
}

func (h *Hub) handleInviteToGroup(msg api.InviteToGroup) {
	group, ok := h.groups[msg.GroupId]
	if !ok {
		log.Println("Group not found")
		return
	}

	host, ok := h.Clients[group.hostId]
	if !ok {
		log.Println("Invalid host id")
		return
	}

	groupInviteNotif := model.NewGroupInviteNotif(
		host.Emoji,
		time.Now().UTC(),
		host.Name,
		group.name,
	)

	for _, memberId := range msg.UserIds {
		client, ok := h.Clients[memberId]
		if !ok {
			log.Printf("Failed to send invite to %v, id not found\n", memberId)
			continue
		}
		log.Printf("Sent to %v", memberId)

		h.sendNotification(client, groupInviteNotif)
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
	log.Println("Member is added to group")
	fmt.Println(group)
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
