package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Env struct {
	db *sql.DB
}

func OpenDb(ctx context.Context) *Env {
	if os.Getenv("DB_CONN_STRING") == "" {
		log.Fatal("DB_CONN_STRING key not defined in .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("DB_CONN_STRING"))

	if err != nil {
		log.Fatal("fuck ", err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Fatal("fuck ", err)
	}
	fmt.Println("yay")

	return &Env{db: db}
}

func (e *Env) MonitorDb() {
	db := e.db
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		if err := db.PingContext(ctx); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 10)

		cancel()
	}
}
