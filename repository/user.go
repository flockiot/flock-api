package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID         string
	ExternalID string
	Email      string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, externalID, email, name string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (external_id, email, name) VALUES ($1, $2, $3)
		 RETURNING id, external_id, email, name, created_at, updated_at`,
		externalID, email, name,
	).Scan(&u.ID, &u.ExternalID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, external_id, email, name, created_at, updated_at FROM users
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&u.ID, &u.ExternalID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByExternalID(ctx context.Context, externalID string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, external_id, email, name, created_at, updated_at FROM users
		 WHERE external_id = $1 AND deleted_at IS NULL`,
		externalID,
	).Scan(&u.ID, &u.ExternalID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user by external id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, external_id, email, name, created_at, updated_at FROM users
		 WHERE deleted_at IS NULL
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.ExternalID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, id, email, name string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET email = $2, name = $3, updated_at = now()
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id, external_id, email, name, created_at, updated_at`,
		id, email, name,
	).Scan(&u.ID, &u.ExternalID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("updating user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
