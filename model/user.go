package model

import "github.com/google/uuid"

type Interest struct {
	Id       int    `json:"id"`
	Emoji    string `json:"emoji"`
	Label    string `json:"label"`
	IsCustom bool   `json:"isCustom"`
}

type User struct {
	Id        string
	Name      string
	Slug      string
	Age       int
	Gender    string
	Interests []Interest
}

func GenerateUserId() string {
	id := uuid.New()
	return id.String()
}
