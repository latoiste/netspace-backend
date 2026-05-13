package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/api"
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
	hub := chat.NewHub()
	go hub.Run()

	r := chi.NewRouter()

	r.Use(mw.Cors)

	r.Get("/ws", handleWs(hub))

	r.Route("/api", func(r chi.Router) {
    r.Post("/sessions/check-in", handleCheckin(env))
    r.Group(func(r chi.Router) {
        r.Use(mw.Auth)
        r.Get("/locations/{slug}", handleLocation(env))
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

func handleWs(hub *chat.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		var req api.WsRequest
		err = conn.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		if req.Type != "register" {
			log.Printf("Expected a \"register\" request type, got \"%v\" instead\n", req.Type)
			conn.Close()
			return
		}

		var body api.RegisterBody
		err = json.Unmarshal(req.Body, &body)

		client := chat.NewClient(
			hub,
			conn,
			body.Username,
			body.Age,
			body.Interest,
		)

		client.Hub.Register <- client

		go client.ReadPump()
		go client.WritePump()
	}
}
