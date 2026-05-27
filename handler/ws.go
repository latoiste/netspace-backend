package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/auth"
	"github.com/latoiste/netspace/chat"
)

func (h *Handler) handleWs(manager *chat.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		params := r.URL.Query()
		tokenString := params.Get("token")

		token, err := h.auth.VerifyToken(tokenString)
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		locationId, err := h.repo.LocationIdBySlug(locationSlug, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			conn.Close()
			return
		}

		hub := manager.LocationHub(locationSlug, locationId)

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		user, err := h.repo.UserById(userId, ctx2)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			conn.Close()
			return
		}

		client := chat.NewClient(
			hub,
			conn,
			userId,
			user.Name,
			"😮",
			locationSlug,
		)

		err = hub.AddClient(client, user)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			conn.Close()
			return
		}

		go client.ReadPump()
		go client.WritePump()
	}
}
