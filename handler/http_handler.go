package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/latoiste/netspace/db"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleLocation(env *db.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		enableCors(&w)

		locationSlug := r.PathValue("slug")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		location, err := env.LocationBySlug(locationSlug, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(location); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
