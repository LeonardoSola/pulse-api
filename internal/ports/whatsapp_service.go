package ports

import (
	"context"
)

type WhatsappService interface {
	CreateDevice(ctx context.Context) (qrChan <-chan string, jidChan <-chan string, err error)
	GetDevice(ctx context.Context, deviceID string) (string, error)
	DeleteDevice(ctx context.Context, deviceID string) error
}
