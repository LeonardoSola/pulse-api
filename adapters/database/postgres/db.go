package postgres

import (
	device_repository "github.com/LeonardoSola/wagon-api/adapters/database/postgres/repositories/devices"
	"github.com/LeonardoSola/wagon-api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbAdapter struct {
	DeviceRepository ports.DeviceRepository
}

func newDB(conn *pgxpool.Pool) *dbAdapter {
	return &dbAdapter{
		DeviceRepository: device_repository.NewDeviceRepository(conn),
	}
}
