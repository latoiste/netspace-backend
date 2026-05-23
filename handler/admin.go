package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/latoiste/netspace/model"
)

func (h *Handler) handleGetActiveSessions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		sessions, err := h.repo.GetActiveSessions(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if sessions == nil {
			sessions = []model.ActiveSession{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(sessions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) handleForceLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userIDStr := r.PathValue("userId")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "invalid userId", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := h.repo.ForceLogoutUser(ctx, userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user successfully logged out",
		})
	}
}

func (h *Handler) handleGetAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		if from == "" || to == "" {
			http.Error(w, "from and to query params are required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		data, err := h.repo.GetAnalytics(ctx, from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if data == nil {
			data = []model.AnalyticsData{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
