package postgres

import (
	device_repository "github.com/LeonardoSola/wagon-api/adapters/database/postgres/repositories/device"
	message_repository "github.com/LeonardoSola/wagon-api/adapters/database/postgres/repositories/message"
	"github.com/LeonardoSola/wagon-api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbAdapter struct {
	DeviceRepository  ports.DeviceRepository
	MessageRepository ports.MessageRepository
}

func newDB(conn *pgxpool.Pool) *dbAdapter {
	return &dbAdapter{
		DeviceRepository:  device_repository.NewDeviceRepository(conn),
		MessageRepository: message_repository.NewMessageRepository(conn),
	}
}
