package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	DatabaseDSN            string `env:"DATABASE_DSN,required"`
	GRPCPort               string `env:"GRPC_PORT" envDefault:"8080"`
	CancelPeriodMin        int    `env:"CANCEL_PERIOD_MIN" envDefault:"5"`
	CancelSchedulerEnabled bool   `env:"CANCEL_SCHEDULER_ENABLED" envDefault:"true"`
	LogLevel               string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Failed to parse config:", err)
	}
	return cfg
}
