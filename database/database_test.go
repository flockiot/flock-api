package database

import (
	"context"
	"testing"
)

func TestConnectInvalidDSN(t *testing.T) {
	_, err := Connect(context.Background(), "postgres://invalid:invalid@localhost:1/nonexistent?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatal("expected error connecting to invalid database")
	}
}

func TestConnectMalformedDSN(t *testing.T) {
	_, err := Connect(context.Background(), "not-a-valid-dsn")
	if err == nil {
		t.Fatal("expected error for malformed DSN")
	}
}
