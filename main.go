package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/latoiste/netspace/db"
	"github.com/latoiste/netspace/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	dbConn := db.OpenDb(ctx)
	defer cancel()

	go db.MonitorDb(dbConn)

	handler.StartServer()
}
