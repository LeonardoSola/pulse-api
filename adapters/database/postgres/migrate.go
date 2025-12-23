package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func runMigrations(conn *pgxpool.Pool) error {
	fmt.Println("⌛ MIGRATIONS")

	// Set goose dialect
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetTableName("schema_migrations")

	db := stdlib.OpenDBFromPool(conn)

	// Run migrations
	if err := goose.Up(db, "adapters/database/postgres/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return db.Close()
}
