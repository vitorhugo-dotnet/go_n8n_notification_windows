package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vitorhugo-java/go_n8n_notification_windows/internal/alert"
	"github.com/vitorhugo-java/go_n8n_notification_windows/internal/config"
	"github.com/vitorhugo-java/go_n8n_notification_windows/internal/notifier"
)

func main() {
	cfg := config.Load()

	log.Printf("starting n8n alert notifier")
	log.Printf("endpoint: %s", cfg.EndpointURL)
	log.Printf("poll interval: %s", cfg.PollInterval)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fetcher := alert.NewFetcher(cfg.EndpointURL)

	var lastEnabled bool

	tick := time.NewTicker(cfg.PollInterval)
	defer tick.Stop()

	// Run immediately on startup, then on each tick.
	poll(ctx, fetcher, cfg, &lastEnabled)

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-tick.C:
			poll(ctx, fetcher, cfg, &lastEnabled)
		}
	}
}

func poll(ctx context.Context, fetcher *alert.Fetcher, cfg *config.Config, lastEnabled *bool) {
	state, err := fetcher.Fetch(ctx)
	if err != nil {
		log.Printf("fetch error: %v", err)
		return
	}

	log.Printf("alert state: enabled=%v priority=%s title=%q", state.Enabled, state.Priority, state.Title)

	// Only notify when alert transitions from disabled → enabled, or while it is enabled.
	if state.Enabled && !*lastEnabled {
		if err := notifier.Send(cfg.AppID, state); err != nil {
			log.Printf("notification error: %v", err)
		}
	}

	*lastEnabled = state.Enabled
}
