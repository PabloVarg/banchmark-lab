package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	run(logger, 200)
}

func run(logger *slog.Logger, workers int) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

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
}
