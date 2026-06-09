package model

import (
	"time"

	"github.com/google/uuid"
)

type Interest struct {
	Id       int    `json:"id"`
	Emoji    string `json:"emoji"`
	Label    string `json:"label"`
	IsCustom bool   `json:"isCustom"`
}

type User struct {
	Id         string
	Name       string
	Age        int
	Gender     string
	Occupation string
	Slug       string
	Emoji      string
	LocationId int
	CreatedAt  time.Time
	IsActive   bool
	Interests  []Interest
}

func GenerateUserId() string {
	id := uuid.New()
	return id.String()
}
