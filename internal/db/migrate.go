package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-admin/internal/db/migrations"

	"github.com/golang-migrate/migrate/v4"
	postgresmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrate applies all pending database migrations.
func Migrate(ctx context.Context, sqlDB *sql.DB) error {
	if sqlDB == nil {
		return fmt.Errorf("sql db is not initialized")
	}

	driver, err := postgresmigrate.WithInstance(sqlDB, &postgresmigrate.Config{})
	if err != nil {
		return fmt.Errorf("create postgres migration driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	done := make(chan error, 1)
	go func() {
		done <- m.Up()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err == nil || errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("run migrations: %w", err)
	}
}
