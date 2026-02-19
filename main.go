package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/flockiot/flock-api/target"
)

func main() {
	targetFlag := flag.String("target", "", "comma-separated list of targets to run, or 'all'")
	flag.Parse()

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.Info("flock-api starting", "targets", targetValue)

	var wg sync.WaitGroup
	for name, fn := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				slog.Error("target failed", "target", name, "error", err)
			}
		}()
	}

	slog.Info("all targets started", "count", len(targets))

	<-ctx.Done()
	wg.Wait()
}
