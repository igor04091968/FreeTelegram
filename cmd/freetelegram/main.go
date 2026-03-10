package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"freetelegram/internal/config"
	"freetelegram/internal/queue"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config JSON")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}

	q := queue.New(queue.Options{
		MaxSize: cfg.Telegram.QueueMax,
	})

	// Placeholder: init modules, router, telemetry
	_ = q

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(map[string]any{
		"status":        "ok",
		"configVersion": cfg.ConfigVersion,
	})
}
