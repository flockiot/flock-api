package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flockiot/flock-api/config"
)

func TestLivez(t *testing.T) {
	r := NewRouter()
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

func TestReadyz(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestNotFoundReturns404(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	r := NewRouter()
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
		errCh <- Start(ctx, cfg)
	}()

	cancel()

	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
