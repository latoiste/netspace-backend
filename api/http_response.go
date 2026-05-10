package api

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
