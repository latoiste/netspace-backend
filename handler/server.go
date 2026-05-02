package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/latoiste/netspace/model"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func StartServer() {
	mux := http.NewServeMux()

	hub := model.NewHub()

	go hub.Run()

	mux.HandleFunc("/ws", handleWs(hub))

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

func handleWs(hub *model.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		var req model.Request
		err = json.Unmarshal(msg, &req)
		if err != nil {
			log.Println("Error parsing json", err)
			conn.Close()
			return
		}

		if req.Type != "register" {
			log.Printf("Expected a \"register\" request type, got \"%v\" instead\n", req.Type)
			conn.Close()
			return
		}

		var body model.RegisterBody
		err = json.Unmarshal(req.Body, &body)

		client := model.NewClient(
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
