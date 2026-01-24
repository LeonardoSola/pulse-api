package device_repository

import (
	"context"

	"github.com/LeonardoSola/pulse-api/internal/ports"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	device_entity "github.com/LeonardoSola/pulse-api/internal/entities/device"
)

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) ports.DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(ctx context.Context, device device_entity.Device) error {
	dml := `INSERT INTO devices (id, name) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, dml, device.ID, device.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceRepository) GetAll(ctx context.Context) ([]device_entity.Device, error) {
	dml := `SELECT * FROM devices`

	rows, err := r.db.Query(ctx, dml)
	if err != nil {
		return []device_entity.Device{}, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[device_entity.Device])
}

func (r *DeviceRepository) Get(ctx context.Context, deviceID string) (device_entity.Device, error) {
	dml := `SELECT * FROM devices WHERE id = $1`

	rows, err := r.db.Query(ctx, dml, deviceID)
	if err != nil {
		return device_entity.Device{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[device_entity.Device])
}

func (r *DeviceRepository) GetByJID(ctx context.Context, jid string) (device_entity.Device, error) {
	dml := `SELECT * FROM devices WHERE jid = $1`

	rows, err := r.db.Query(ctx, dml, jid)
	if err != nil {
		return device_entity.Device{}, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[device_entity.Device])
}

func (r *DeviceRepository) Update(ctx context.Context, device device_entity.Device) error {
	dml := `UPDATE devices SET name = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.Exec(ctx, dml, device.Name, device.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceRepository) Delete(ctx context.Context, deviceID string) error {
	dml := `DELETE FROM devices WHERE id = $1`

	_, err := r.db.Exec(ctx, dml, deviceID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceRepository) UpdateJID(ctx context.Context, deviceID string, jid *string) error {
	dml := `UPDATE devices SET jid = $1 WHERE id = $2`

	_, err := r.db.Exec(ctx, dml, jid, deviceID)
	if err != nil {
		return err
	}

	return nil
}
