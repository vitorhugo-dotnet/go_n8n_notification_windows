package config

import (
	"flag"
	"os"
	"time"
)

// Config holds all runtime configuration for the notifier.
type Config struct {
	EndpointURL  string
	PollInterval time.Duration
	AppID        string
}

// Load parses flags and environment variables, flags take precedence.
func Load() *Config {
	defaultURL := getEnv("N8N_ALERT_ENDPOINT", "https://n8n.hugojava.dev/webhook/alert-state")
	defaultInterval := getEnvDuration("POLL_INTERVAL", 30*time.Second)
	defaultAppID := getEnv("APP_ID", "go_n8n_notification_windows")

	url := flag.String("endpoint", defaultURL, "n8n alert-state endpoint URL")
	interval := flag.Duration("interval", defaultInterval, "polling interval (e.g. 30s, 1m)")
	appID := flag.String("appid", defaultAppID, "Windows App User Model ID for toast notifications")
	flag.Parse()

	return &Config{
		EndpointURL:  *url,
		PollInterval: *interval,
		AppID:        *appID,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
