package whatsapp

import (
	"context"
	"fmt"

	"github.com/LeonardoSola/pulse-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var container *sqlstore.Container

func InitContainer() error {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	var err error

	container, err = sqlstore.New(context.Background(), "pgx", config.Config.Database.ConnectionString, dbLog)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := container.Upgrade(context.Background()); err != nil {
		return fmt.Errorf("failed to upgrade database: %w", err)
	}

	return nil
}
