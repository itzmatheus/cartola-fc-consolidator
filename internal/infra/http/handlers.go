package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/itzmatheus/cartola-fc-consolidator/internal/infra/db"
)

func ListPlayerHandler(ctx context.Context, queries db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		players, err := queries.FindAllPlayers(ctx)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(players)
	}
}
