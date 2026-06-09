package api

import "github.com/latoiste/netspace/model"

type CreateUserRequest struct {
	LocationSlug string           `json:"locationSlug"`
	Name         string           `json:"name"`
	Slug         string           `json:"slug"`
	Age          int              `json:"age"`
	Gender       string           `json:"gender"`
	Occupation   string           `json:"occupation"`
	Interests    []model.Interest `json:"interests"`
}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
