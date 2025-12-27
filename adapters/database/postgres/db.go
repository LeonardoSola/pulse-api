package postgres

import (
	device_repository "github.com/LeonardoSola/wagon-api/adapters/database/postgres/repositories/devices"
	"github.com/LeonardoSola/wagon-api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbAdapter struct {
	DevicesRepository ports.DevicesRepository
}

func newDB(conn *pgxpool.Pool) *dbAdapter {
	return &dbAdapter{
		DevicesRepository: device_repository.NewDevicesRepository(conn),
	}
}
