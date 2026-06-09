package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Handler) handleNotification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		sessionData, ok := r.Context().Value("SessionData").(model.SessionData)
		if !ok {
			log.Println("Invalid session data value")
			http.Error(w, "Invalid session data value", http.StatusUnauthorized)
			return
		}

		userId := sessionData.UserId

		if userId == "" {
			log.Println("No UserId in token")
			http.Error(w, "No UserId in token", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		notifs, err := h.repo.NotificationByUserId(userId, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := api.ConstructGetNotificationsResponse(notifs)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
