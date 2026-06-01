package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Handler) handleChatList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		sessionData, ok := r.Context().Value("SessionData").(model.SessionData)
		if !ok {
			log.Println("Invalid session data value")
			http.Error(w, "Invalid session data kontol", http.StatusUnauthorized)
			return
		}

		userId := sessionData.UserId
		if userId == "" {
			log.Println("No UserId in token")
			http.Error(w, "No UserId in token", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		privateMsgs, err := h.repo.AllPrivateMessages(userId, ctx)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel2()

		groupMsgs, err := h.repo.AllGroupMessages(userId, ctx2)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := api.ConstructGetChatListResponse(privateMsgs, groupMsgs)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
