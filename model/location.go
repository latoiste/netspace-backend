package model

import (
	"fmt"
	"time"
)

type Location struct {
	Id        int
	Slug      string
	Name      string
	Address   string
	PartnerId string
	JoinDate  time.Time
	Capacity  int
	Timezone  string
	IsActive  bool
	QrToken   string
	QrLabel   string
}

func (l *Location) FormatTimezoneLabel() string {
	loc, _ := time.LoadLocation(l.Timezone)

	name, offset := time.Now().In(loc).Zone()
	hours := offset / 3600

	return fmt.Sprintf("%s · UTC%+d", name, hours)
}
