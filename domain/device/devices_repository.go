package device

import (
	"context"
)

type DevicesRepository interface {
	GetAll(ctx context.Context) ([]Device, error)
	Get(ctx context.Context, deviceID string) (Device, error)
	Create(ctx context.Context, device Device) error
	Update(ctx context.Context, device Device) error
	Delete(ctx context.Context, deviceID string) error
	//
	UpdateJID(ctx context.Context, deviceID string, jid *string) error
}
