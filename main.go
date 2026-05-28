package main

import (
	"github.com/latoiste/netspace/app"
	"github.com/latoiste/netspace/handler"
)

func main() {
	env := app.NewEnv()

	handler := handler.NewHandler(
		env.Repo,
		env.Auth,
		env.Blacklist,
	)

	go env.Repo.MonitorDb()

	handler.StartServer()
}
