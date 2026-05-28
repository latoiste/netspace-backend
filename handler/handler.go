package handler

import (
	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/db"
)

type Handler struct {
	repo      *db.Repository
	auth      *auth.Auth
	blacklist *auth.Blacklist
}

func NewHandler(repo *db.Repository, auth *auth.Auth, blacklist *auth.Blacklist) *Handler {
	return &Handler{
		repo:      repo,
		auth:      auth,
		blacklist: blacklist,
	}
}
