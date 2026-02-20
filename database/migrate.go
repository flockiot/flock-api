package database

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(dsn string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %w", err)
	}

	migrateDSN := strings.Replace(dsn, "postgres://", "pgx5://", 1)
	m, err := migrate.NewWithSourceInstance("iofs", source, migrateDSN)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}
