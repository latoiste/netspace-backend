package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/latoiste/netspace/db"
)

func handleLocation(env *db.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// defer r.Body.Close()
		enableCors(&w, r)

		locationSlug := r.PathValue("slug")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		location, err := env.LocationBySlug(locationSlug, ctx)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(location); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Location query success")
	}
}
