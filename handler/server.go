package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/chat"
	"github.com/latoiste/netspace/db"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

// nanti pake middleware
func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func StartServer(env *db.Env) {
	mux := http.NewServeMux()

	hub := chat.NewHub()

	go hub.Run()

	mux.HandleFunc("/ws", handleWs(hub))

	mux.HandleFunc("/api/locations/{slug}", handleLocation(env))
	mux.HandleFunc("/api/locations/{slug}/users", handleLocationUsers(env))

	mux.HandleFunc("/api/sessions/check-in", handleCheckin(env))

	server := http.Server{
		Handler: mux,
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
