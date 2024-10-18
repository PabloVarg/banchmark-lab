package main

import (
	"log/slog"
	"net/http"
)

func routes(logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /health", healthHandler(logger))

	return mux
}
