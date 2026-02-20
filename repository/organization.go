package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

func (r *OrganizationRepository) Create(ctx context.Context, name string) (*Organization, error) {
	org := &Organization{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO organizations (name) VALUES ($1)
		 RETURNING id, name, created_at, updated_at`,
		name,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating organization: %w", err)
	}
	return org, nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id string) (*Organization, error) {
	org := &Organization{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at, updated_at FROM organizations
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting organization by id: %w", err)
	}
	return org, nil
}

func (r *OrganizationRepository) GetByName(ctx context.Context, name string) (*Organization, error) {
	org := &Organization{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at, updated_at FROM organizations
		 WHERE name = $1 AND deleted_at IS NULL`,
		name,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting organization by name: %w", err)
	}
	return org, nil
}

func (r *OrganizationRepository) List(ctx context.Context, limit, offset int) ([]*Organization, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, created_at, updated_at FROM organizations
		 WHERE deleted_at IS NULL
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing organizations: %w", err)
	}
	defer rows.Close()

	var orgs []*Organization
	for rows.Next() {
		org := &Organization{}
		if err := rows.Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning organization: %w", err)
		}
		orgs = append(orgs, org)
	}
	return orgs, rows.Err()
}

func (r *OrganizationRepository) Update(ctx context.Context, id, name string) (*Organization, error) {
	org := &Organization{}
	err := r.pool.QueryRow(ctx,
		`UPDATE organizations SET name = $2, updated_at = now()
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, name, created_at, updated_at`,
		id, name,
	).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("updating organization: %w", err)
	}
	return org, nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE organizations SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting organization: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}
