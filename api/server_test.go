package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flockiot/flock-api/config"
	"github.com/flockiot/flock-api/version"
)

func testRouter(pool *pgxpool.Pool) http.Handler {
	return NewRouter(pool)
}

func TestLivez(t *testing.T) {
	r := testRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Fatalf("expected 'ok', got %q", w.Body.String())
	}
}

func TestReadyzWithoutDB(t *testing.T) {
	version.Value = "1.2.3"
	r := testRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when pool is nil, got %d", w.Code)
	}

	var resp readyzResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Version != "1.2.3" {
		t.Fatalf("expected version '1.2.3', got %q", resp.Version)
	}
	if resp.Status != "database not configured" {
		t.Fatalf("expected status 'database not configured', got %q", resp.Status)
	}
}

func TestNotFoundReturns404(t *testing.T) {
	r := testRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	r := testRouter(nil)
	req := httptest.NewRequest(http.MethodPost, "/livez", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestStartListensAndShutdown(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 0,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- Start(ctx, cfg, nil)
	}()

	cancel()

	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
