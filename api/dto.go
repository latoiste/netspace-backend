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

type InterestPercentageDTO struct {
	InterestDTO
	Percentage int `json:"Percentage"`
}

type NotificationDTO struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Emoji           string `json:"emoji"`
	AvatarGradient  string `json:"avatarGradient"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	TimestampString string `json:"timestamp"`
	Unread          bool   `json:"unread"`
	PrimaryLabel    string `json:"primaryLabel,omitempty"`
	SecondaryLabel  string `json:"secondaryLabel,omitempty"`
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

func ConstructNotificationDTO(notif model.Notification) NotificationDTO {
	return NotificationDTO{
		Id:              notif.Id,
		Type:            notif.Type,
		Emoji:           notif.Emoji,
		AvatarGradient:  notif.AvatarGradient,
		Title:           notif.Title,
		Description:     notif.Description,
		TimestampString: notif.Timestamp.Local().Format("15:04"),
		Unread:          notif.Unread,
		PrimaryLabel:    notif.PrimaryLabel,
		SecondaryLabel:  notif.SecondaryLabel,
	}
}
