package target

import (
	"context"
	"testing"
)

func noop(_ context.Context) error { return nil }

func newTestRegistry() *Registry {
	r := New()
	r.Register("api", noop)
	r.Register("ingester", noop)
	r.Register("scheduler", noop)
	return r
}

func TestRegisterAndGet(t *testing.T) {
	r := newTestRegistry()

	fn, ok := r.Get("api")
	if !ok {
		t.Fatal("expected api target to exist")
	}
	if fn == nil {
		t.Fatal("expected non-nil start function")
	}
}

func TestGetUnknown(t *testing.T) {
	r := newTestRegistry()
	_, ok := r.Get("nonexistent")
	if ok {
		t.Fatal("expected nonexistent target to not be found")
	}
}

func TestNames(t *testing.T) {
	r := newTestRegistry()
	names := r.Names()
	expected := []string{"api", "ingester", "scheduler"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d names, got %d", len(expected), len(names))
	}
	for i, name := range expected {
		if names[i] != name {
			t.Errorf("expected names[%d] = %q, got %q", i, name, names[i])
		}
	}
}

func TestNamesIsACopy(t *testing.T) {
	r := newTestRegistry()
	names := r.Names()
	names[0] = "modified"
	if r.Names()[0] == "modified" {
		t.Fatal("Names() should return a copy")
	}
}

func TestResolveAll(t *testing.T) {
	r := newTestRegistry()
	resolved, err := r.Resolve("all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(resolved))
	}
}

func TestResolveSingle(t *testing.T) {
	r := newTestRegistry()
	resolved, err := r.Resolve("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 1 {
		t.Fatalf("expected 1 target, got %d", len(resolved))
	}
	if _, ok := resolved["api"]; !ok {
		t.Fatal("expected api target in resolved set")
	}
}

func TestResolveMultiple(t *testing.T) {
	r := newTestRegistry()
	resolved, err := r.Resolve("api,ingester")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(resolved))
	}
}

func TestResolveUnknown(t *testing.T) {
	r := newTestRegistry()
	_, err := r.Resolve("bogus")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestResolveEmpty(t *testing.T) {
	r := newTestRegistry()
	_, err := r.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty target string")
	}
}

func TestDuplicateRegisterPanics(t *testing.T) {
	r := New()
	r.Register("api", noop)
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate register")
		}
	}()
	r.Register("api", noop)
}
