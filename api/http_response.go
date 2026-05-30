package api

import "github.com/latoiste/netspace/model"

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

func ConstructGetUsersResponse(users []model.User) GetUsersResponse {
	onlineCount := len(users)
	userDTOs := make([]UserDTO, 0)

	for _, user := range users {
		userDTO := ConstructUserDTO(user)

		userDTOs = append(userDTOs, userDTO)
	}

	return GetUsersResponse{
		Users:       userDTOs,
		OnlineCount: onlineCount,
	}
}
