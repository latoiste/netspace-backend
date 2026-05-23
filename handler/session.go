package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/model"
)

func (h *Handler) handleCheckin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		var reqBody api.CreateUserRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			log.Println(err)
			return
		}

		ctx1, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		location, err := h.repo.LocationBySlug(reqBody.LocationSlug, ctx1)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userId := model.GenerateUserId()
		token, err := h.auth.GenerateJWT(userId)
		if err != nil {
			log.Println("Error generating token", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
		locationId, err := h.repo.LocationIdBySlug(reqBody.LocationSlug, ctx2)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cancel()

		response := api.CreateUserResponse{
			UserId:       userId,
			SessionToken: token,
			LocationSlug: location.Slug,
			LocationName: location.Name,
		}

		user := model.User{
			Id:         userId,
			Name:       reqBody.Name,
			LocationId: locationId,
			Age:        reqBody.Age,
			Gender:     reqBody.Gender,
			Interests:  reqBody.Interests,
			Slug:       reqBody.Slug,
		}

		ctx3, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err = h.repo.InsertUser(user, ctx3); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
