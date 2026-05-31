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

type AdminLoginResponse struct {
	Success bool     `json:"success"`
	Admin   AdminDTO `json:"admin"`
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
