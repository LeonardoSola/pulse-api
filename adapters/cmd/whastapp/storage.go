package whastapp

import (
	"context"
	"database/sql"
	"fmt"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func getStorage(sqlInstance *sql.DB) (*store.Device, error) {
	var container *sqlstore.Container
	dbLog := waLog.Stdout("Database", "ERROR", true)
	ctx := context.Background()

	container = sqlstore.NewWithDB(sqlInstance, "postgres", dbLog)

	if err := container.Upgrade(ctx); err != nil {
		return nil, fmt.Errorf("failed to upgrade database: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device store: %w", err)
	}

	return deviceStore, nil
}
