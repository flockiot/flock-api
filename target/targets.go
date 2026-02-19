package target

import (
	"context"
	"log/slog"
)

func DefaultRegistry() *Registry {
	r := New()
	r.Register("api", placeholder("api"))
	r.Register("ingester", placeholder("ingester"))
	r.Register("scheduler", placeholder("scheduler"))
	r.Register("builder", placeholder("builder"))
	r.Register("delta", placeholder("delta"))
	r.Register("registry-proxy", placeholder("registry-proxy"))
	r.Register("tunnel", placeholder("tunnel"))
	r.Register("proxy", placeholder("proxy"))
	r.Register("events-gateway", placeholder("events-gateway"))
	return r
}

func placeholder(name string) StartFunc {
	return func(ctx context.Context) error {
		slog.Info("target started", "target", name)
		<-ctx.Done()
		return nil
	}
}
