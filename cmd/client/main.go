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

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	run(logger, 200)
}

func run(logger *slog.Logger, workers int) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go exposeMetrics(ctx, logger, &wg)

	var group errgroup.Group
	group.SetLimit(workers)

	dialer := http.Client{
		Timeout: time.Minute,
	}

	for range workers {
		group.Go(func() error {
			if err := RepeatDialHealth(ctx, "http://server:8000/health", dialer, logger, time.Second); err != nil {
				logger.Error("error dialing", "err", err)
				return err
			}

			return nil
		})
	}
	logger.Info("launched workers", "count", workers)

	if err := group.Wait(); err != nil {
		logger.Error("a dialer failed", "err", err)
		return
	}
	logger.Info("workers done", "workers", workers)

	logger.Info("closing resources")
	wg.Wait()
}

func exposeMetrics(ctx context.Context, logger *slog.Logger, wg *sync.WaitGroup) {
	port := 2112

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	go func(ctx context.Context) {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down server", "err", err)
			return
		}
	}(ctx)

	logger.InfoContext(ctx, "metrics listening", "port", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("metrics exited unexpectedly", "err", err)
		return
	}
}
