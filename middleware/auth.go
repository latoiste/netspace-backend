package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/latoiste/netspace/auth"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.ExtractTokenFromHeader(r)
		if err != nil {
			log.Println(err)
			return
		}

		token, err := auth.VerifyToken(tokenString)
		if err != nil || token == nil {
			log.Println(err)
			return
		}

		claim, ok := token.Claims.(*auth.Claim)
		if !ok {
			log.Println("Invalid token fields")
			return
		}

		ctx := context.WithValue(r.Context(), "UserId", claim.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
