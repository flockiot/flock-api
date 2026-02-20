package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	Log      LogConfig      `envPrefix:"LOG_"`
}

type ServerConfig struct {
	Host     string `env:"HOST"      envDefault:"0.0.0.0"`
	Port     int    `env:"PORT"      envDefault:"8080"`
	GRPCPort int    `env:"GRPC_PORT" envDefault:"9090"`
}

type PostgresConfig struct {
	DSN string `env:"DSN" envDefault:"postgres://flock:flock@localhost:5432/flock?sslmode=disable"`
}

type LogConfig struct {
	Level  string `env:"LEVEL"  envDefault:"info"`
	Format string `env:"FORMAT" envDefault:"json"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Prefix: "FLOCK_",
	})
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
