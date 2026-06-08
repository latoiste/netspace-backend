package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/model"
)

func (m *Middleware) Auth(a *auth.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := a.ExtractTokenFromHeader(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if m.blacklist.IsBlacklisted(tokenString) {
				log.Println("Token is invalid")
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

			sessionData := model.SessionData{
				UserId:      claim.UserId,
				TokenString: tokenString,
			}

			ctx := context.WithValue(r.Context(), "SessionData", sessionData)

			next.ServeHTTP(w, r.WithContext(ctx))
			// next.ServeHTTP(w, r)
		})
	}
}
