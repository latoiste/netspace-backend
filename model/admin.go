package model

import (
	"time"
)

type Admin struct {
	Id       string
	Username string
	Password string
	Role     string
	Plan     string
	Avatar   string
	Name     string
}

type ActiveSession struct {
	UserID       int       `json:"userId"`
	Name         string    `json:"name"`
	Age          int       `json:"age"`
	Gender       string    `json:"gender"`
	TableNumber  string    `json:"tableNumber"`
	Interest     []string  `json:"interest"`
	CurrentJob   string    `json:"currentJob"`
	LocationID   int       `json:"locationId"`
	LastActiveAt time.Time `json:"lastActiveAt"`
}

type AnalyticsData struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	Delta     string `json:"delta"`
	DeltaType string `json:"deltaType"`
}
