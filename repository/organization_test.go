package repository

import (
	"context"
	"os"
	"testing"

	"github.com/flockiot/flock-api/database"
)

func testPool(t *testing.T) *database.TestDB {
	t.Helper()
	dsn := os.Getenv("FLOCK_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://flock:flock@localhost:5432/flock?sslmode=disable"
	}
	db, err := database.NewTestDB(context.Background(), dsn)
	if err != nil {
		t.Skip("skipping: database not available:", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestOrganizationCreate(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	org, err := repo.Create(context.Background(), "acmecorp")
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}
	if org.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if org.Name != "acmecorp" {
		t.Fatalf("expected name 'acmecorp', got %q", org.Name)
	}
}

func TestOrganizationGetByID(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	created, err := repo.Create(context.Background(), "getbyid-org")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	found, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("failed to get by id: %v", err)
	}
	if found == nil {
		t.Fatal("expected organization, got nil")
	}
	if found.ID != created.ID {
		t.Fatalf("expected id %q, got %q", created.ID, found.ID)
	}
}

func TestOrganizationGetByIDNotFound(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	found, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil for nonexistent organization")
	}
}

func TestOrganizationGetByName(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	_, err := repo.Create(context.Background(), "byname-org")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	found, err := repo.GetByName(context.Background(), "byname-org")
	if err != nil {
		t.Fatalf("failed to get by name: %v", err)
	}
	if found == nil {
		t.Fatal("expected organization, got nil")
	}
	if found.Name != "byname-org" {
		t.Fatalf("expected name 'byname-org', got %q", found.Name)
	}
}

func TestOrganizationList(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	for i := range 3 {
		name := "list-org-" + string(rune('a'+i))
		if _, err := repo.Create(context.Background(), name); err != nil {
			t.Fatalf("failed to create: %v", err)
		}
	}

	orgs, err := repo.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(orgs) < 3 {
		t.Fatalf("expected at least 3 organizations, got %d", len(orgs))
	}
}

func TestOrganizationUpdate(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	created, err := repo.Create(context.Background(), "oldname")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	updated, err := repo.Update(context.Background(), created.ID, "newname")
	if err != nil {
		t.Fatalf("failed to update: %v", err)
	}
	if updated.Name != "newname" {
		t.Fatalf("expected name 'newname', got %q", updated.Name)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Fatal("expected updated_at to advance")
	}
}

func TestOrganizationDelete(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	created, err := repo.Create(context.Background(), "deleteme")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	if err := repo.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	found, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil after soft delete")
	}
}

func TestOrganizationDeleteNotFound(t *testing.T) {
	db := testPool(t)
	repo := NewOrganizationRepository(db.Pool)

	err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error deleting nonexistent organization")
	}
}
