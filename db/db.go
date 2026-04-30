package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func OpenDb(ctx context.Context) *sql.DB {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 dbname=Netspace user=postgres connect_timeout=10 sslmode=prefer")

	if err != nil {
		log.Fatal("fuck ", err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Fatal("fuck ", err)
	}
	fmt.Println("yay")

	return db
}

func MonitorDb(db *sql.DB) {
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		if err := db.PingContext(ctx); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 10)
		log.Print("yay")
	}
}
