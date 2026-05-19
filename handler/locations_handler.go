package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/api"
	"github.com/latoiste/netspace/db"
)

func handleLocation(env *db.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		locationSlug := r.PathValue("slug")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		location, err := env.LocationBySlug(locationSlug, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(location); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Location query success")
	}
}

func handleLocationUsers(env *db.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		locationSlug := r.PathValue("slug")

		ctx1, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		id, err := env.LocationIdBySlug(locationSlug, ctx1)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		users, err := env.UsersInLocation(id, ctx2)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := api.ConstructGetUsersResponse(users)

		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
