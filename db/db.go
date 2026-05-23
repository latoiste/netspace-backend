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

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func OpenDb(ctx context.Context) *sql.DB {
	datasourceName := os.Getenv("DB_CONN_STRING")
	if datasourceName == "" {
		log.Fatal("DB_CONN_STRING key not defined in .env file")
	}
	db, err := sql.Open("postgres", datasourceName)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("yay")

	return db
}

func (r *Repository) MonitorDb() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		if err := r.db.PingContext(ctx); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 10)

		cancel()
	}
}
