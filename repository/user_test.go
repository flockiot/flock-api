package repository

import (
	"context"
	"testing"
)

func TestUserCreate(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	user, err := repo.Create(context.Background(), "zitadel-123", "alice@example.com", "alice")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if user.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if user.ExternalID != "zitadel-123" {
		t.Fatalf("expected external_id 'zitadel-123', got %q", user.ExternalID)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected email 'alice@example.com', got %q", user.Email)
	}
	if user.Name != "alice" {
		t.Fatalf("expected name 'alice', got %q", user.Name)
	}
}

func TestUserGetByID(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	created, err := repo.Create(context.Background(), "zitadel-getbyid", "bob@example.com", "bob")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	found, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("failed to get by id: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.ID != created.ID {
		t.Fatalf("expected id %q, got %q", created.ID, found.ID)
	}
}

func TestUserGetByIDNotFound(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	found, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil for nonexistent user")
	}
}

func TestUserGetByExternalID(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	_, err := repo.Create(context.Background(), "zitadel-ext", "ext@example.com", "ext")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	found, err := repo.GetByExternalID(context.Background(), "zitadel-ext")
	if err != nil {
		t.Fatalf("failed to get by external id: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.ExternalID != "zitadel-ext" {
		t.Fatalf("expected external_id 'zitadel-ext', got %q", found.ExternalID)
	}
}

func TestUserGetByExternalIDNotFound(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	found, err := repo.GetByExternalID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil for nonexistent external id")
	}
}

func TestUserList(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	for i := range 3 {
		ext := "zitadel-list-" + string(rune('a'+i))
		email := ext + "@example.com"
		if _, err := repo.Create(context.Background(), ext, email, ext); err != nil {
			t.Fatalf("failed to create: %v", err)
		}
	}

	users, err := repo.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(users) < 3 {
		t.Fatalf("expected at least 3 users, got %d", len(users))
	}
}

func TestUserUpdate(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	created, err := repo.Create(context.Background(), "zitadel-update", "old@example.com", "oldname")
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	updated, err := repo.Update(context.Background(), created.ID, "new@example.com", "newname")
	if err != nil {
		t.Fatalf("failed to update: %v", err)
	}
	if updated.Email != "new@example.com" {
		t.Fatalf("expected email 'new@example.com', got %q", updated.Email)
	}
	if updated.Name != "newname" {
		t.Fatalf("expected name 'newname', got %q", updated.Name)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Fatal("expected updated_at to advance")
	}
}

func TestUserDelete(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	created, err := repo.Create(context.Background(), "zitadel-delete", "del@example.com", "delme")
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

func TestUserDeleteNotFound(t *testing.T) {
	db := testPool(t)
	repo := NewUserRepository(db.Pool)

	err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error deleting nonexistent user")
	}
}
