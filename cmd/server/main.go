package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	run(log, 8000)
}

func run(logger *slog.Logger, port int) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes(logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	go shutdown(ctx, &wg, logger, &server)

	logger.InfoContext(ctx, "server listening", "port", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("error stopping server", "err", err)
		return
	}
}

func shutdown(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger, server *http.Server) {
	defer wg.Done()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("error shutting down server", "err", err)
		return
	}
}
