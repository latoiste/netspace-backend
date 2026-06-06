package main

import (
	"fmt"
	"log"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/app"
	"github.com/latoiste/netspace/handler"
)

func main() {
	env := app.NewEnv()

	handler := handler.NewHandler(
		env.Repo,
		env.Auth,
		env.Blacklist,
		env.Manager,
	)

	go env.Repo.MonitorDb()

	now := time.Now()

	todayStart := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
	todayCurrent := now

	yesterdayStart := todayStart.AddDate(0, 0, -1)

	log.Println(todayStart.UTC())
	log.Println(todayCurrent.UTC())
	log.Println(yesterdayStart.UTC())

	log.Println(api.ConstructAnalyticsDTO(10, 23, "Check-in Hari Ini", "up"))

	fmt.Println(env.Auth.GenerateJWT("4301cfd6-4337-4b90-973c-360b96daa811"))
	fmt.Println(env.Auth.GenerateJWT("4e92049c-8900-413e-82ce-68b85a393e2a"))
	fmt.Println(env.Auth.GenerateJWT("f810831d-1eae-483f-b484-22e40aabbb93"))
	fmt.Println(env.Auth.GenerateJWT("2424a831-c983-4231-93cf-92624dcfbff4"))

	handler.StartServer()
}
