package handler

import (
	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/chat"
	"github.com/latoiste/netspace/db"
)

type Handler struct {
	repo      *db.Repository
	auth      *auth.Auth
	blacklist *auth.Blacklist
	manager   *chat.Manager
}

func NewHandler(repo *db.Repository, auth *auth.Auth, blacklist *auth.Blacklist, manager *chat.Manager) *Handler {
	return &Handler{
		repo:      repo,
		auth:      auth,
		blacklist: blacklist,
		manager:   manager,
	}
}
