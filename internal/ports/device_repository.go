package ports

import (
	"context"

	device_entity "github.com/LeonardoSola/wagon-api/internal/entities/device"
)

type DeviceRepository interface {
	GetAll(ctx context.Context) ([]device_entity.Device, error)
	Get(ctx context.Context, deviceID string) (device_entity.Device, error)
	Create(ctx context.Context, device device_entity.Device) error
	Update(ctx context.Context, device device_entity.Device) error
	Delete(ctx context.Context, deviceID string) error
	//
	UpdateJID(ctx context.Context, deviceID string, jid *string) error
}
