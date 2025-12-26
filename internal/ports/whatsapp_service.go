package ports

import (
	"context"

	"github.com/LeonardoSola/wagon-api/domain/message"
)

type WhatsappService interface {
	CreateDevice(ctx context.Context) (qrChan <-chan string, jidChan <-chan string, err error)
	DeleteDevice(ctx context.Context, deviceID string) error
	CheckDevice(ctx context.Context, deviceID string) (bool, error)
	SendTextMessage(ctx context.Context, message message.Message) error
}
