package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerLogsFields(t *testing.T) {
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))

	handler := requestLogger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	log := buf.String()
	for _, field := range []string{"method", "path", "status", "duration_ms", "remote_ip", "bytes"} {
		if !strings.Contains(log, field) {
			t.Errorf("log missing field %q: %s", field, log)
		}
	}
}

func TestResponseWriterCapturesStatus(t *testing.T) {
	w := httptest.NewRecorder()
	ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
	ww.WriteHeader(http.StatusNotFound)

	if ww.status != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", ww.status)
	}
}

func TestResponseWriterCapturesBytes(t *testing.T) {
	w := httptest.NewRecorder()
	ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
	ww.Write([]byte("hello"))

	if ww.bytes != 5 {
		t.Fatalf("expected 5 bytes, got %d", ww.bytes)
	}
}
