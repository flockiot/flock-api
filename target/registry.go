package target

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

type StartFunc func(ctx context.Context) error

type Registry struct {
	targets map[string]StartFunc
	order   []string
}

func New() *Registry {
	return &Registry{
		targets: make(map[string]StartFunc),
	}
}

func (r *Registry) Register(name string, fn StartFunc) {
	if _, exists := r.targets[name]; exists {
		panic(fmt.Sprintf("target %q already registered", name))
	}
	r.targets[name] = fn
	r.order = append(r.order, name)
}

func (r *Registry) Get(name string) (StartFunc, bool) {
	fn, ok := r.targets[name]
	return fn, ok
}

func (r *Registry) Names() []string {
	return slices.Clone(r.order)
}

func (r *Registry) Resolve(raw string) (map[string]StartFunc, error) {
	if raw == "all" {
		result := make(map[string]StartFunc, len(r.targets))
		for name, fn := range r.targets {
			result[name] = fn
		}
		return result, nil
	}

	names := strings.Split(raw, ",")
	result := make(map[string]StartFunc, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		fn, ok := r.targets[name]
		if !ok {
			return nil, fmt.Errorf("unknown target %q (valid targets: %s)", name, strings.Join(r.order, ", "))
		}
		result[name] = fn
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no targets specified (valid targets: %s)", strings.Join(r.order, ", "))
	}
	return result, nil
}
