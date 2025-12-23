package postgres

import (
	devices_repository "github.com/LeonardoSola/wagon-api/adapters/database/postgres/repositories/devices"
	"github.com/LeonardoSola/wagon-api/domain/device"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbAdapter struct {
	DevicesRepository device.DevicesRepository
}

func newDB(conn *pgxpool.Pool) *dbAdapter {
	return &dbAdapter{
		DevicesRepository: devices_repository.NewDevicesRepository(conn),
	}
}
