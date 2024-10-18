package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func DialHealth(ctx context.Context, addr string, dialer http.Client, logger *slog.Logger) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return err
	}

	response, err := dialer.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	io.Copy(io.Discard, response.Body)

	return nil
}

func RepeatDialHealth(ctx context.Context, addr string, dialer http.Client, logger *slog.Logger, delay time.Duration) error {
	ticker := time.NewTicker(delay)

	for {
		select {
		case <-ticker.C:
			if err := DialHealth(ctx, addr, dialer, logger); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
