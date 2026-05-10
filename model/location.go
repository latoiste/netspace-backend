package model

type Location struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}
