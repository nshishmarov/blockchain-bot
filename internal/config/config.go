package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/sethvargo/go-envconfig"
)

// --- Type defifnitions

type TgBotConfig struct {
	TgToken  string `env:"TG_TOKEN"`
	TgChanId int64 `env:"TG_CHAN_ID"`
	DebugMode bool `env:"DEBUG_MODE, default=false"`
}

type Config struct{}

// Load loads the configuration from the environment.
func (c *Config) Load(ctx context.Context) TgBotConfig {
	var cfg TgBotConfig
	slog.Info("started loading configuration")
	if err := envconfig.Process(ctx, &cfg); err != nil {
		slog.Error("error with loading config", "error", err.Error())
		os.Exit(1)
	}
	return cfg
}
