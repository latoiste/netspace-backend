package api

import "github.com/latoiste/netspace/model"

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
