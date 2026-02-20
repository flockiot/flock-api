package database

import (
	"testing"
)

func TestMigrateInvalidDSN(t *testing.T) {
	err := Migrate("pgx5://invalid:invalid@localhost:1/nonexistent?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Fatal("expected error migrating with invalid DSN")
	}
}

func TestMigrationsEmbed(t *testing.T) {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		t.Fatalf("failed to read embedded migrations: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one migration file")
	}
}
