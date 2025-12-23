package postgres

import (
	"context"
	"fmt"

	"github.com/LeonardoSola/wagon-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB   *dbAdapter
	conn *pgxpool.Pool
)

func NewConnection() error {
	pool, err := pgxpool.New(context.Background(), config.Config.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	err = runMigrations(pool)
	if err != nil {
		return err
	}

	DB = newDB(pool)

	return nil
}

func Close() {
	conn.Close()
}
