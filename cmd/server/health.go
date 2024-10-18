package main

import (
	"log/slog"
	"net/http"

	"github.com/PabloVarg/benchmark-lab/internal/helpers"
)

func healthHandler(logger *slog.Logger) http.Handler {
	type output struct {
		Alive bool `json:"alive"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := output{
			Alive: true,
		}

		if err := helpers.WriteJSON(w, http.StatusOK, response); err != nil {
			helpers.ServerError(w, logger, err)
		}
	})
}
