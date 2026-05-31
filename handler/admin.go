package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// func (h *Handler) handleGetAnalyticsMetrics() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()

// 		analytics := make([]model.AnalyticsData, 0, 4)

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
// 		defer cancel()

// 		checkInsToday, err := h.repo.TotalCheckInRange(
// 			time.Now().Truncate(time.Hour*24),
// 			time.Now(),
// 			ctx,
// 		)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, err, http.StatusInternalServerError)
// 			return
// 		}

// 		// ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
// 		// defer cancel()

// 		// checkInsYesterday, err := h.repo.TotalCheckInRange(

// 		// )

// 		_ = model.AnalyticsData{
// 			Label:     "Check-in Hari Ini",
// 			Value:     string(checkInsToday),
// 			Delta:     string(4),
// 			DeltaType: "",
// 		}

// 		if err := json.NewEncoder(w).Encode(analytics); err != nil {
// 			log.Println(err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

func (h *Handler) handleTopInterests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		response, err := h.repo.GetTopInterests(ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// func (h *Handler) handleGetActiveSessions() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 		defer cancel()

// 		sessions, err := h.repo.GetActiveSessions(ctx)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		if sessions == nil {
// 			sessions = []model.ActiveSession{}
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		if err = json.NewEncoder(w).Encode(sessions); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }

// func (h *Handler) handleForceLogout() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()

// 		userIDStr := r.PathValue("userId")
// 		userID, err := strconv.Atoi(userIDStr)
// 		if err != nil {
// 			http.Error(w, "invalid userId", http.StatusBadRequest)
// 			return
// 		}

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 		defer cancel()

// 		if err := h.repo.ForceLogoutUser(ctx, userID); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"message": "user successfully logged out",
// 		})
// 	}
// }
