package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/chat"
	"github.com/latoiste/netspace/db"
)

type Env struct {
	Repo      *db.Repository
	Auth      *auth.Auth
	Blacklist *auth.Blacklist
	Manager   *chat.Manager
}

func NewEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	dbConn := db.OpenDb(ctx)
	defer cancel()

	jwtKey := os.Getenv("JWT_SECRET_KEY")

	authCredentials := auth.NewAuth([]byte(jwtKey))
	repo := db.NewRepo(dbConn)

	manager := chat.NewManager(repo)

	return &Env{
		Repo:      repo,
		Auth:      authCredentials,
		Blacklist: auth.NewBlacklist(),
		Manager:   manager,
	}
}
