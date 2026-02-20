package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/flockiot/flock-api/config"
	"github.com/flockiot/flock-api/database"
	"github.com/flockiot/flock-api/logging"
	"github.com/flockiot/flock-api/target"
	"github.com/flockiot/flock-api/version"
)

//go:embed VERSION
var embeddedVersion string

func main() {
	version.Value = strings.TrimSpace(embeddedVersion)

	targetFlag := flag.String("target", "", "comma-separated list of targets to run, or 'all'")
	migrateOnly := flag.Bool("migrate-only", false, "run database migrations and exit")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := logging.Setup(cfg.Log.Level, cfg.Log.Format, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error setting up logging: %v\n", err)
		os.Exit(1)
	}

	slog.Info("flock-api configured",
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
		"log_level", cfg.Log.Level,
		"log_format", cfg.Log.Format,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := database.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected")

	slog.Info("running database migrations")
	if err := database.Migrate(cfg.Postgres.DSN); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("database migrations complete")

	if *migrateOnly {
		return
	}

	targetValue := *targetFlag
	if targetValue == "" {
		targetValue = os.Getenv("FLOCK_TARGET")
	}
	if targetValue == "" {
		fmt.Fprintf(os.Stderr, "error: --target flag or FLOCK_TARGET env var is required\n")
		fmt.Fprintf(os.Stderr, "valid targets: %v\n", target.DefaultRegistry().Names())
		os.Exit(1)
	}

	registry := target.DefaultRegistry()
	targets, err := registry.Resolve(targetValue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	deps := &target.Deps{
		Config: cfg,
		DB:     pool,
	}

	slog.Info("flock-api starting", "targets", targetValue)

	var wg sync.WaitGroup
	for name, fn := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx, deps); err != nil {
				slog.Error("target failed", "target", name, "error", err)
			}
		}()
	}

	slog.Info("all targets started", "count", len(targets))

	<-ctx.Done()
	wg.Wait()
}
