package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.Server.GRPCPort != 9090 {
		t.Errorf("Server.GRPCPort = %d, want %d", cfg.Server.GRPCPort, 9090)
	}
	if cfg.Postgres.DSN != "postgres://flock:flock@localhost:5432/flock?sslmode=disable" {
		t.Errorf("Postgres.DSN = %q, want default", cfg.Postgres.DSN)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "json")
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("FLOCK_SERVER_HOST", "127.0.0.1")
	t.Setenv("FLOCK_SERVER_PORT", "3000")
	t.Setenv("FLOCK_SERVER_GRPC_PORT", "3001")
	t.Setenv("FLOCK_POSTGRES_DSN", "postgres://prod:secret@db:5432/flock")
	t.Setenv("FLOCK_LOG_LEVEL", "debug")
	t.Setenv("FLOCK_LOG_FORMAT", "json")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "127.0.0.1")
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 3000)
	}
	if cfg.Server.GRPCPort != 3001 {
		t.Errorf("Server.GRPCPort = %d, want %d", cfg.Server.GRPCPort, 3001)
	}
	if cfg.Postgres.DSN != "postgres://prod:secret@db:5432/flock" {
		t.Errorf("Postgres.DSN = %q, want override", cfg.Postgres.DSN)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "debug")
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "json")
	}
}

func TestLoadPartialOverride(t *testing.T) {
	t.Setenv("FLOCK_SERVER_PORT", "4000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 4000 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 4000)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want default %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.GRPCPort != 9090 {
		t.Errorf("Server.GRPCPort = %d, want default %d", cfg.Server.GRPCPort, 9090)
	}
}
