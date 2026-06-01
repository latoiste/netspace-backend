package handler

import (
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
