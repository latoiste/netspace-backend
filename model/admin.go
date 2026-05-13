package model

import "time"

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
	Date          string `json:"date"`
	ActiveUsers   int    `json:"activeUsers"`
	TotalSessions int    `json:"totalSessions"`
	MessageCount  int    `json:"messageCount"`
}