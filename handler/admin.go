package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/api"
)

func (h *Handler) handleAnalyticsMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		analytics := make([]api.AnalyticsDTO, 0, 3)

		now := time.Now()

		todayStart := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0, 0, 0, 0,
			now.Location(),
		)
		todayCurrent := now

		yesterdayStart := todayStart.AddDate(0, 0, -1)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		checkInToday, checkInYest, err := h.repo.CheckInAnalytics(yesterdayStart, todayStart, todayCurrent, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var checkInDeltaType string
		if checkInToday >= checkInYest {
			checkInDeltaType = "up"
		} else {
			checkInDeltaType = "down"
		}
		checkInAnalytics := api.ConstructAnalyticsDTO(checkInYest, checkInToday, "Check-in Hari Ini", checkInDeltaType)

		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel2()

		activeUsersToday, activeUsersYest, err := h.repo.ActiveUsersAnalytics(yesterdayStart, todayStart, todayCurrent, ctx2)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activeUsersAnalytics := api.ConstructAnalyticsDTO(activeUsersYest, activeUsersToday, "User Aktif Sekarang", "live")

		ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel3()

		messagesToday, messagesYest, err := h.repo.MessagesAnalytics(yesterdayStart, todayStart, todayCurrent, ctx3)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var messagesDeltaType string
		if messagesToday >= messagesYest {
			messagesDeltaType = "up"
		} else {
			messagesDeltaType = "down"
		}
		messagesAnalytics := api.ConstructAnalyticsDTO(messagesYest, messagesToday, "Total konversasi hari ini", messagesDeltaType)

		analytics = append(analytics, checkInAnalytics)
		analytics = append(analytics, activeUsersAnalytics)
		analytics = append(analytics, messagesAnalytics)

		if err := json.NewEncoder(w).Encode(analytics); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) handleHourlyCheckIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		now := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		response, err := h.repo.HourlyCheckInAnalytics(8, now, ctx)
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

func (h *Handler) handleAdminLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var reqBody api.AdminLoginRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Write([]byte(
				`{
					\"success\": false,
					\"error\": \"Username dan password harus diisi\"
				}`,
			))
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			w.Write([]byte(
				`{
					\"success\": false,
					\"error\": \"Username dan password harus diisi\"
				}`,
			))
			return
		}

		username := reqBody.Username
		password := reqBody.Password

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		admin, err := h.repo.AdminByUsername(username, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: ganti pake hash + insert new admin kalo ada waktu
		w.Header().Set("Content-Type", "application/json")
		if admin.Password == password {
			adminDto := api.ConstructAdminDTO(*admin)
			response := api.AdminLoginResponse{
				Success: true,
				Admin:   adminDto,
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.Write([]byte(
				`{
					\"success\": false,
					\"error\": \"Username dan password harus diisi\"
				}`,
			))
		}
	}
}

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

func (h *Handler) handleLocationDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
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

		response := api.GetLocationDetailResponse{
			Slug:       location.Slug,
			Name:       location.Name,
			Address:    location.Address,
			PartnerId:  location.PartnerId,
			JoinedDate: location.JoinDate.Format("02 Jan 2006"),
			Capacity:   fmt.Sprintf("~ %v user", location.Capacity),
			Timezone:   location.FormatTimezoneLabel(),
			IsActive:   location.IsActive,
			QrToken:    location.QrToken,
			QrLabel:    location.QrLabel,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) handleToggleLocationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		locationSlug := r.PathValue("slug")

		type Payload struct {
			IsActive bool `json:"isActive"`
		}

		var reqBody Payload

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		err = h.repo.UpdateLocationIsActive(locationSlug, reqBody.IsActive, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			w.Write([]byte("{\"success\": true}"))
			return
		}
		w.Write([]byte("{\"success\": true}"))
	}
}

func (h *Handler) handleActiveUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
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

		response := api.ConstructGetActiveUsersResponse(users)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) handleForceLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		userId := r.PathValue("userId")

		type Payload struct {
			Reason string `json:"isActive"`
		}

		var reqBody Payload

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := h.repo.UpdateUserIsActive(userId, false, ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.manager.ForceLogoutUser(userId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"success": true,
		})
	}
}
