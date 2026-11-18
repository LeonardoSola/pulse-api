package ports

import (
	"context"

	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
)

type MessageRepository interface {
	Create(ctx context.Context, message message_entity.Message) error
	Get(ctx context.Context, messageID string) (message_entity.Message, error)
	GetAll(ctx context.Context) ([]message_entity.Message, error)
	Update(ctx context.Context, message message_entity.Message) error
	Delete(ctx context.Context, messageID string) error
}
