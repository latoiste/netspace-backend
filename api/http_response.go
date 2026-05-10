package api

type CreateUserResponse struct {
	UserId       string `json:"userId"`
	SessionToken string `json:"sessionToken"`
	LocationSlug string `json:"locationSlug"`
	LocationName string `json:"locationName"`
}
