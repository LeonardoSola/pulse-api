package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LeonardoSola/pulse-api/internal/logger"
	"github.com/LeonardoSola/pulse-api/migrations"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func InitMigrations() (int64, error) {
	db := stdlib.OpenDBFromPool(conn)

	gooseProvider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		migrations.FS,
		goose.WithSlog(logger.Log.With(slog.String("component", "migrations"))),
		goose.WithTableName("schema_migrations"),
	)
	if err != nil {
		db.Close()
		return 0, fmt.Errorf("failed to create goose provider: %w", err)
	}

	// Run migrations
	applied, err := gooseProvider.Up(context.Background())
	if err != nil {
		db.Close()
		return 0, fmt.Errorf("failed to run migrations: %w", err)
	}

	if len(applied) > 0 {
		for _, migration := range applied {
			logger.Log.Info("Migrations applied", "migration", migration)
		}
	}

	current, err := gooseProvider.GetDBVersion(context.Background())
	if err != nil {
		db.Close()
		return 0, fmt.Errorf("failed to read migration version: %w", err)
	}

	return current, db.Close()
}
