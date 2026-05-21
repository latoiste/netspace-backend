package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/chat"
	"github.com/latoiste/netspace/db"
	mw "github.com/latoiste/netspace/middleware"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func StartServer(env *db.Env) {
	manager := chat.NewManager()

	r := chi.NewRouter()

	r.Use(mw.Cors)

	r.Get("/ws", handleWs(manager, env))

	r.Route("/api", func(r chi.Router) {
		r.Post("/sessions/check-in", handleCheckin(env))
		r.Get("/locations/{slug}", handleLocation(env))
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth)
			r.Get("/locations/{slug}/users", handleLocationUsers(env))
		})
		r.Route("/admin", func(r chi.Router) {
			r.Get("/sessions", handleGetActiveSessions(env))
			r.Post("/force-logout/{userId}", handleForceLogout(env))
			r.Get("/analytics", handleGetAnalytics(env))
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
