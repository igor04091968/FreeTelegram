package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"freetelegram/internal/config"
	"freetelegram/internal/queue"
	"freetelegram/internal/router"
	"freetelegram/internal/telemetry"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config JSON")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}

	_ = queue.New(queue.Options{MaxSize: cfg.Telegram.QueueMax})

	stats := telemetry.New()
	updater := router.New(router.Options{
		Domains:            cfg.Routing.Domains,
		DNSServers:         cfg.Routing.DNSServerList,
		AggregateThreshold: cfg.Routing.AggregateThreshold,
		ApplyRoutes:        cfg.Routing.ApplyRoutes,
		Interface:          cfg.Routing.Interface,
	})
	interval := time.Duration(cfg.Routing.UpdateIntervalSec) * time.Second
	worker := router.NewWorker(updater, stats, interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.Run(ctx)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(map[string]any{
		"status":        "ok",
		"configVersion": cfg.ConfigVersion,
		"telemetry":     stats.Snapshot(),
	})
}
