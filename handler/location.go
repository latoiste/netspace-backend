package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Handler) handleLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		locationSlug := r.PathValue("slug")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		location, err := h.repo.LocationBySlug(locationSlug, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(location); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Location query success")
	}
}

func (h *Handler) handleLocationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		locationSlug := r.PathValue("slug")

		ctx1, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		id, err := h.repo.LocationIdBySlug(locationSlug, ctx1)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		users, err := h.repo.UsersInLocation(id, ctx2)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := api.ConstructGetUsersResponse(users)

		// Filter the requesting user out of the "who's here" list when a session
		// is present. The auth middleware may be disabled in dev (no SessionData),
		// in which case we simply return the full list instead of failing.
		if sessionData, ok := r.Context().Value("SessionData").(model.SessionData); ok && sessionData.UserId != "" {
			userId := sessionData.UserId
			response.Users = slices.DeleteFunc(response.Users, func(u api.UserDTO) bool {
				return u.Id == userId
			})
			response.OnlineCount = len(response.Users)
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
