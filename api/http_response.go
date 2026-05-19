package api

import "github.com/latoiste/netspace/model"

type CreateUserResponse struct {
	UserId       string `json:"userId"`
	SessionToken string `json:"sessionToken"`
	LocationSlug string `json:"locationSlug"`
	LocationName string `json:"locationName"`
}

type GetUsersResponse struct {
	Users       []UserOutput `json:"users"`
	OnlineCount int          `json:"onlineCount"`
}

type UserOutput struct {
	Id        string           `json:"id"`
	Slug      string           `json:"slug"`
	Name      string           `json:"name"`
	Emoji     string           `json:"emoji"`
	Interests []InterestOutput `json:"interests"`
}

type InterestOutput struct {
	Emoji string `json:"emoji"`
	Label string `json:"label"`
}

func ConstructGetUsersResponse(users []model.User) GetUsersResponse {
	onlineCount := len(users)
	userOutputs := make([]UserOutput, 0)

	for _, user := range users {
		var userOutput UserOutput
		interestsOutputs := make([]InterestOutput, 0)

		for _, interest := range user.Interests {
			interestOutput := InterestOutput{
				Emoji: interest.Emoji,
				Label: interest.Label,
			}
			interestsOutputs = append(interestsOutputs, interestOutput)
		}

		userOutput = UserOutput{
			Id:        user.Id,
			Slug:      user.Slug,
			Name:      user.Name,
			Emoji:     "😮",
			Interests: interestsOutputs,
		}

		userOutputs = append(userOutputs, userOutput)
	}

	return GetUsersResponse{
		Users:       userOutputs,
		OnlineCount: onlineCount,
	}
}
