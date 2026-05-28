package middleware

import (
	"net/http"

	"github.com/latoiste/netspace/auth"
)

func (m *Middleware) Auth(a *auth.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// tokenString, err := a.ExtractTokenFromHeader(r)
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }

			// if m.blacklist.IsBlacklisted(tokenString) {
			// 	log.Println("Token is invalid")
			// 	return
			// }

			// token, err := a.VerifyToken(tokenString)
			// if err != nil || token == nil {
			// 	log.Println(err)
			// 	return
			// }

			// claim, ok := token.Claims.(*auth.Claim)
			// if !ok {
			// 	log.Println("Invalid token fields")
			// 	return
			// }

			// sessionData := model.SessionData{
			// 	UserId:      claim.UserId,
			// 	TokenString: tokenString,
			// }

			// ctx := context.WithValue(r.Context(), "SessionData", sessionData)

			// next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
		})
	}
}
