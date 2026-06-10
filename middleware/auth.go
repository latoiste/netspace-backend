package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/model"
)

func (m *Middleware) authForActor(a *auth.Auth, requiredActor string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := a.ExtractTokenFromHeader(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if m.blacklist.IsBlacklisted(tokenString) {
				log.Println("Token is blacklisted")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			token, err := a.VerifyToken(tokenString)
			if err != nil || token == nil {
				log.Println(err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claim, ok := token.Claims.(*auth.Claim)
			if !ok {
				log.Println("Invalid token fields")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			actorType := claim.ActorType
			// Tokens issued before actorType was introduced are user sessions.
			if actorType == "" {
				actorType = "user"
			}
			if actorType != requiredActor {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			sessionData := model.SessionData{
				UserId:       claim.UserId,
				TokenString:  tokenString,
				ActorType:    actorType,
				LocationSlug: claim.LocationSlug,
			}

			ctx := context.WithValue(r.Context(), "SessionData", sessionData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *Middleware) UserAuth(a *auth.Auth) func(next http.Handler) http.Handler {
	return m.authForActor(a, "user")
}

func (m *Middleware) AdminAuth(a *auth.Auth) func(next http.Handler) http.Handler {
	return m.authForActor(a, "admin")
}
