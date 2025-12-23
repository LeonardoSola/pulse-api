package devices_repository

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/LeonardoSola/wagon-api/domain/device"
)

type DevicesRepository struct {
	db *pgxpool.Pool
}

func NewDevicesRepository(db *pgxpool.Pool) device.DevicesRepository {
	return &DevicesRepository{db: db}
}

func (r *DevicesRepository) Create(ctx context.Context, device device.Device) error {
	dml := `INSERT INTO devices (id, name) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, dml, device.ID, device.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *DevicesRepository) GetAll(ctx context.Context) ([]device.Device, error) {
	dml := `SELECT * FROM devices`

	rows, err := r.db.Query(ctx, dml)
	if err != nil {
		return []device.Device{}, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[device.Device])
}

func (r *DevicesRepository) Get(ctx context.Context, deviceID string) (device.Device, error) {
	dml := `SELECT * FROM devices WHERE id = $1`

	rows, err := r.db.Query(ctx, dml, deviceID)
	if err != nil {
		return device.Device{}, err
	}

	dev, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[device.Device])
	if err != nil {
		return device.Device{}, err
	}

	return dev, nil
}

func (r *DevicesRepository) Update(ctx context.Context, device device.Device) error {
	dml := `UPDATE devices SET name = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.Exec(ctx, dml, device.Name, device.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DevicesRepository) Delete(ctx context.Context, deviceID string) error {
	dml := `DELETE FROM devices WHERE id = $1`

	_, err := r.db.Exec(ctx, dml, deviceID)
	if err != nil {
		return err
	}

	return nil
}

func (r *DevicesRepository) UpdateJID(ctx context.Context, deviceID string, jid *string) error {
	dml := `UPDATE devices SET jid = $1 WHERE id = $2`

	_, err := r.db.Exec(ctx, dml, jid, deviceID)
	if err != nil {
		return err
	}

	return nil
}
