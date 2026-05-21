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

type UserDTO struct {
	Id        string        `json:"id"`
	Slug      string        `json:"slug"`
	Name      string        `json:"name"`
	Emoji     string        `json:"emoji"`
	Interests []InterestDTO `json:"interests"`
}

type InterestDTO struct {
	Emoji string `json:"emoji"`
	Label string `json:"label"`
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

func ConstructUserDTO(user model.User) UserDTO {
	interestDTOs := make([]InterestDTO, 0, len(user.Interests))

	for _, interest := range user.Interests {
		interestDTO := InterestDTO{
			Emoji: interest.Emoji,
			Label: interest.Label,
		}
		interestDTOs = append(interestDTOs, interestDTO)
	}

	return UserDTO{
		Id:        user.Id,
		Slug:      user.Slug,
		Name:      user.Name,
		Emoji:     "😮",
		Interests: interestDTOs,
	}
}
