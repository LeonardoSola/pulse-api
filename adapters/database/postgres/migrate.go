package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/LeonardoSola/wagon-api/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func runMigrations(conn *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(conn)

	gooseProvider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS("adapters/database/postgres/migrations"),
		goose.WithSlog(logger.Log.With(slog.String("component", "migrations"))),
		goose.WithTableName("schema_migrations"),
	)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	// Run migrations
	migrationsRunned, err := gooseProvider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if len(migrationsRunned) > 0 {
		for _, migration := range migrationsRunned {
			logger.Log.Info("Migrations applied", "migration", migration)
		}
	}

	current, _, _ := gooseProvider.GetVersions(context.Background())
	logger.Log.Info("✅ Migrations", "current", current)

	return db.Close()
}
