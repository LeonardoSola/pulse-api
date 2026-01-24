package ports

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message message_entity.Message) error
	ReceiveMessage(ctx context.Context, message message_entity.Message) error
}
