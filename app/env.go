package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/db"
)

type Env struct {
	Repo *db.Repository
	Auth *auth.Auth
}

func NewEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	dbConn := db.OpenDb(ctx)
	defer cancel()

	jwtKey := os.Getenv("JWT_SECRET_KEY")

	auth := auth.NewAuth([]byte(jwtKey))
	repo := db.NewRepo(dbConn)

	return &Env{
		Repo: repo,
		Auth: auth,
	}
}
