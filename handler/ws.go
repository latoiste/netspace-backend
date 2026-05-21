package handler

import (
	"log"
	"net/http"

	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/chat"
)

func handleWs(manager *chat.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		params := r.URL.Query()
		tokenString := params.Get("token")

		token, err := auth.VerifyToken(tokenString)
		if err != nil {
			log.Println("Faield to verify token")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claim, ok := token.Claims.(*auth.Claim)
		if !ok {
			log.Println("Invalid token fields")
			http.Error(w, "Invalid token fields", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId := claim.UserId
		locationSlug := params.Get("locationSlug")

		hub := manager.LocationHub(locationSlug)

		client := chat.NewClient(
			hub,
			conn,
			userId,
			locationSlug,
		)
		hub.Clients[client] = true

		log.Println("Added new client")

		go client.ReadPump()
		go client.WritePump()
	}
}
