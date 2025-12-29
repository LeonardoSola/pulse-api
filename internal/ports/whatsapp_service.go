package ports

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
)

type WhatsappService interface {
	CreateDevice(ctx context.Context) (qrChan <-chan string, jidChan <-chan string, err error)
	DeleteDevice(ctx context.Context, deviceID string) error
	CheckDevice(ctx context.Context, deviceID string) (bool, error)
	SendTextMessage(ctx context.Context, message message_entity.Message) (messageID string, err error)
}
