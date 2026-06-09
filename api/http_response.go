package api

import (
	"github.com/latoiste/netspace/model"
)

type CreateUserResponse struct {
	UserId       string `json:"userId"`
	SessionToken string `json:"sessionToken"`
	LocationSlug string `json:"locationSlug"`
	LocationName string `json:"locationName"`
}

type GetUsersResponse struct {
	Users       []UserDTO `json:"users"`
	OnlineCount int       `json:"onlineCount"`
}

type GetNotificationsResponse struct {
	Notifications []NotificationDTO `json:"notifications"`
}

type GetActiveUsersResponse struct {
	Users []UserDTO
}

type GetLocationDetailResponse struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	PartnerId  string `json:"partnerId"`
	JoinedDate string `json:"joinedDate"`
	Capacity   string `json:"capacity"`
	Timezone   string `json:"timezone"`
	IsActive   bool   `json:"isActive"`
	QrToken    string `json:"qrToken"`
	QrLabel    string `json:"qrLabel"`
}

type GetChatListRespone struct {
	Chats []MessageDTO `json:"chats"`
}

type GetDMHistoryResponse struct {
	User     DMPartnerDTO     `json:"user"`
	Messages []ChatMessageDTO `json:"messages"`
}

type GetGroupHistoryResponse struct {
	Group    GroupSummaryDTO       `json:"group"`
	Members  []GroupMemberDTO      `json:"members"`
	Messages []GroupChatMessageDTO `json:"messages"`
}

type GetPublicHistoryResponse struct {
	Messages []PublicChatMessageDTO `json:"messages"`
}

type AdminLoginResponse struct {
	Success bool      `json:"success"`
	Admin   *AdminDTO `json:"admin,omitempty"`
	Token   string    `json:"token,omitempty"`
	Error   string    `json:"error,omitempty"`
}

func ConstructGetUsersResponse(users []model.User) GetUsersResponse {
	onlineCount := len(users)
	userDTOs := make([]UserDTO, len(users))

	for i, user := range users {
		userDTO := ConstructUserDTO(user)

		userDTOs[i] = userDTO
	}

	return GetUsersResponse{
		Users:       userDTOs,
		OnlineCount: onlineCount,
	}
}

func ConstructGetNotificationsResponse(notifs []model.Notification) GetNotificationsResponse {
	notifDTOs := make([]NotificationDTO, len(notifs))

	for i, notif := range notifs {
		notifDTO := ConstructNotificationDTO(notif)

		notifDTOs[i] = notifDTO
	}

	return GetNotificationsResponse{
		Notifications: notifDTOs,
	}
}

func ConstructGetActiveUsersResponse(users []model.User) []ActiveUserDTO {
	activeUserDTOs := make([]ActiveUserDTO, 0, len(users))

	for _, user := range users {
		activeUserDTO := ConstructActiveUserDTO(user)
		activeUserDTOs = append(activeUserDTOs, activeUserDTO)
	}

	return activeUserDTOs
}

func ConstructGetChatListResponse(privateMsgs []MessageDTO, groupMsgs []MessageDTO) GetChatListRespone {
	messages := append(privateMsgs, groupMsgs...)
	return GetChatListRespone{
		Chats: messages,
	}
}

func ConstructGetDMHistoryResponse(partner model.User, msgs []model.PrivateMessage, myId string) GetDMHistoryResponse {
	messageDTOs := make([]ChatMessageDTO, 0, len(msgs))
	for _, m := range msgs {
		messageDTOs = append(messageDTOs, ConstructChatMessageDTO(m, myId))
	}

	return GetDMHistoryResponse{
		User:     ConstructDMPartnerDTO(partner),
		Messages: messageDTOs,
	}
}

// senders maps userId -> user, covering current members plus any senders who
// have since left the group, so old messages still render a name/emoji.
func ConstructGetGroupHistoryResponse(group model.Group, members []model.User, msgs []model.GroupMessage, senders map[string]model.User, myId string) GetGroupHistoryResponse {
	memberDTOs := make([]GroupMemberDTO, 0, len(members))
	for _, m := range members {
		memberDTOs = append(memberDTOs, GroupMemberDTO{
			Id:     m.Id,
			Name:   m.Name,
			Emoji:  EmojiForUser(m.Id),
			IsHost: false,
		})
	}

	messageDTOs := make([]GroupChatMessageDTO, 0, len(msgs))
	for _, msg := range msgs {
		sender := senders[msg.SenderId]
		messageDTOs = append(messageDTOs, GroupChatMessageDTO{
			Id:          msg.MessageId,
			SenderId:    msg.SenderId,
			SenderName:  sender.Name,
			SenderEmoji: EmojiForUser(msg.SenderId),
			Message:     msg.Message,
			Timestamp:   msg.Timestamp.Local().Format("15:04"),
			IsMine:      msg.SenderId == myId,
		})
	}

	return GetGroupHistoryResponse{
		Group: GroupSummaryDTO{
			Name:  group.Name,
			Emoji: "☕",
		},
		Members:  memberDTOs,
		Messages: messageDTOs,
	}
}
