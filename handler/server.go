package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	mw "github.com/latoiste/netspace/middleware"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func (h *Handler) StartServer() {
	// Clear stale "active" flags left by a previous run that died before its
	// disconnect handlers fired. On a fresh boot nobody holds a live socket yet,
	// so any isactive=true row is a ghost that would otherwise inflate the admin
	// roster and analytics until it happened to reconnect or time out.
	resetCtx, resetCancel := context.WithTimeout(context.Background(), time.Second*5)
	if err := h.repo.DeactivateAllUsers(resetCtx); err != nil {
		log.Println("failed to reset active users on startup:", err)
	}
	resetCancel()

	// Keep the shared public timeline bounded: sweep out messages older than the
	// retention window now and hourly thereafter.
	go h.runPublicMessageRetention()

	mw := mw.NewMiddleware(h.blacklist)

	r := chi.NewRouter()

	r.Use(mw.Cors)

	r.Get("/ws", h.handleWs(h.manager))

	r.Route("/api", func(r chi.Router) {
		r.Post("/sessions/check-in", h.handleCheckin())
		r.Get("/locations/{slug}", h.handleLocation())

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth(h.auth))
			r.Get("/locations/{slug}/users", h.handleLocationUsers())
			r.Get("/sessions/logout", h.handleLogout())

			r.Get("/notifications", h.handleNotification())

			r.Get("/chats", h.handleChatList())
			r.Get("/chats/{userId}/messages", h.handleDMHistory())
			r.Get("/groups/{groupId}/messages", h.handleGroupHistory())
			r.Get("/locations/{slug}/public-messages", h.handlePublicHistory())
		})

		r.Post("/admin/login", h.handleAdminLogin())

		r.Route("/admin", func(r chi.Router) {
			r.Use(mw.Auth(h.auth))
			r.Get("/dashboard/stats", h.handleAnalyticsMetrics())
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/interests", h.handleTopInterests())
				r.Get("/hourly", h.handleHourlyCheckIn())
			})

			r.Route("/locations", func(r chi.Router) {
				r.Get("/{slug}", h.handleLocationDetail())
				r.Put("/{slug}", h.handleToggleLocationStatus())
				r.Get("/{slug}/users", h.handleActiveUsers())
			})

			r.Post("/users/{userId}/kick", h.handleForceLogout())
		})
	})

	server := http.Server{
		Handler: r,
		Addr:    ":" + "8080",
	}

	log.Println("Server is listening on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// publicMessageRetention is how long a public-room message lives before the
// hourly sweep removes it. Matches the window the history endpoint serves.
const publicMessageRetention = 24 * time.Hour

// runPublicMessageRetention deletes public messages older than the retention
// window on boot and once an hour after, so the venue's shared timeline stays
// bounded and relevant to the current crowd.
func (h *Handler) runPublicMessageRetention() {
	sweep := func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		removed, err := h.repo.DeletePublicMessagesBefore(time.Now().Add(-publicMessageRetention), ctx)
		if err != nil {
			log.Println("public retention sweep:", err)
			return
		}
		if removed > 0 {
			log.Printf("public retention: removed %d old message(s)", removed)
		}
	}

	sweep()
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		sweep()
	}
}
