package main

import (
	"context"
	"time"

	"github.com/latoiste/netspace/db"
	"github.com/latoiste/netspace/handlers"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	dbConn := db.OpenDb(ctx)
	defer cancel()

	go db.MonitorDb(dbConn)

	handlers.StartServer()
}
