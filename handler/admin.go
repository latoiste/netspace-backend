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

// wib is the app's reporting timezone (Asia/Jakarta, UTC+7, no DST). Analytics
// buckets (per-hour, per-day) are computed in WIB so they match what the café
// admin actually sees — independent of the server's timezone (UTC in prod).
// A fixed offset avoids depending on tzdata being present in the image.
var wib = time.FixedZone("WIB", 7*60*60)

func (h *Handler) handleAnalyticsMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		analytics := make([]api.AnalyticsDTO, 0, 3)

		now := time.Now().In(wib)

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

		// "User Aktif Sekarang" = live WebSocket connections right now, not the
		// DB isactive flag. The flag drifts: a process killed before its
		// disconnect handlers run leaves ghosts marked active, and the previous
		// query also only counted users created *today*. Counting live sockets
		// reflects who is actually connected this instant.
		activeNow := h.manager.TotalOnline()
		activeUsersAnalytics := api.ConstructAnalyticsDTO(0, activeNow, "User Aktif Sekarang", "live")

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
		messagesAnalytics := api.ConstructAnalyticsDTO(messagesYest, messagesToday, "Total Konversasi Hari Ini", messagesDeltaType)

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

		now := time.Now().In(wib)

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
		w.Header().Set("Content-Type", "application/json")

		writeLoginError := func(status int, message string) {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(api.AdminLoginResponse{
				Success: false,
				Error:   message,
			})
		}

		var reqBody api.AdminLoginRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			writeLoginError(http.StatusBadRequest, "Username dan password harus diisi")
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			log.Println(err)
			writeLoginError(http.StatusBadRequest, "Username dan password harus diisi")
			return
		}

		username := reqBody.Username
		password := reqBody.Password

		if username == "" || password == "" {
			writeLoginError(http.StatusBadRequest, "Username dan password harus diisi")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		admin, err := h.repo.AdminByUsername(username, ctx)
		if err != nil {
			log.Println(err)
			writeLoginError(http.StatusUnauthorized, "Username atau password salah")
			return
		}

		// TODO: ganti pake hash + insert new admin kalo ada waktu
		if admin.Password != password {
			writeLoginError(http.StatusUnauthorized, "Username atau password salah")
			return
		}

		token, err := h.auth.GenerateJWT(admin.Id)
		if err != nil {
			log.Println(err)
			writeLoginError(http.StatusInternalServerError, "Gagal membuat sesi login")
			return
		}

		adminDto := api.ConstructAdminDTO(*admin)
		response := api.AdminLoginResponse{
			Success: true,
			Admin:   &adminDto,
			Token:   token,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			Reason string `json:"reason"`
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
