package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/chat"
	mw "github.com/latoiste/netspace/middleware"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func (h *Handler) StartServer() {
	manager := chat.NewManager(h.repo)
	mw := mw.NewMiddleware(h.blacklist)

	r := chi.NewRouter()

	r.Use(mw.Cors)

	r.Get("/ws", h.handleWs(manager))

	r.Route("/api", func(r chi.Router) {
		r.Post("/sessions/check-in", h.handleCheckin())
		r.Get("/locations/{slug}", h.handleLocation())

		r.Group(func(r chi.Router) {
			r.Use(mw.Auth(h.auth))
			r.Get("/locations/{slug}/users", h.handleLocationUsers())
			r.Get("/sessions/logout", h.handleLogout())

			r.Get("/notifications", h.handleNotification())
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(mw.Auth(h.auth))
			r.Get("/analytics/interests", h.handleTopInterests())
			r.Get("/locations/{slug}", h.handleLocationDetail())
			// r.Get("/sessions", h.handleGetActiveSessions())
			// r.Post("/force-logout/{userId}", h.handleForceLogout())
			// r.Get("/dashboard/stats", h.handleGetAnalyticsMetrics())
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
