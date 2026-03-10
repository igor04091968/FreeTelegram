package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const CurrentConfigVersion = 1

type Config struct {
	ConfigVersion int `json:"config_version"`
	Telegram      TelegramConfig `json:"telegram"`
	Routing       RoutingConfig `json:"routing"`
	Modules       ModulesConfig `json:"modules"`
}

type TelegramConfig struct {
	APIKey   string  `json:"api_key"`
	AdminIDs []int64 `json:"admin_ids"`
	QueueMax int     `json:"queue_max"`
}

type RoutingConfig struct {
	RouteFor           []string `json:"route_for"`
	ExcludeRouteFor    []string `json:"exclude_route_for"`
	DNSServerList      []string `json:"dns_server_list"`
	Domains            []string `json:"domains"`
	AggregateThreshold int      `json:"aggregate_threshold"`
	ApplyRoutes        bool     `json:"apply_routes"`
	Interface          string   `json:"interface"`
	UpdateIntervalSec  int      `json:"update_interval_sec"`
}

type ModulesConfig struct {
	Enabled []string `json:"enabled"`
	Default string   `json:"default"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.ConfigVersion == 0 {
		return nil, fmt.Errorf("config_version is required")
	}
	if cfg.ConfigVersion != CurrentConfigVersion {
		return nil, fmt.Errorf("unsupported config_version: %d", cfg.ConfigVersion)
	}
	return &cfg, nil
}
