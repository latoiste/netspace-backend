package model

import (
	"fmt"
	"time"
)

type Location struct {
	Id        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	PartnerId string    `json:"partnerId"`
	JoinDate  time.Time `json:"joinDate"`
	Capacity  int       `json:"capacity"`
	Timezone  string    `json:"timezone"`
	IsActive  bool      `json:"isActive"`
	QrToken   string    `json:"qrToken"`
	QrLabel   string    `json:"qrLabel"`
}

func (l *Location) FormatTimezoneLabel() string {
	loc, _ := time.LoadLocation(l.Timezone)

	name, offset := time.Now().In(loc).Zone()
	hours := offset / 3600

	return fmt.Sprintf("%s · UTC%+d", name, hours)
}
