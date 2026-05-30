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
