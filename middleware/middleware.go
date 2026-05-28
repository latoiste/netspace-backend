package middleware

import "github.com/latoiste/netspace/auth"

type Middleware struct {
	blacklist *auth.Blacklist
}

func NewMiddleware(blacklist *auth.Blacklist) *Middleware {
	return &Middleware{
		blacklist: blacklist,
	}
}
