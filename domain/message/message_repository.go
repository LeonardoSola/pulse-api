package message

import "context"

type MessageRepository interface {
	Create(ctx context.Context, message Message) error
	Get(ctx context.Context, messageID string) (Message, error)
	GetAll(ctx context.Context) ([]Message, error)
	Update(ctx context.Context, message Message) error
	Delete(ctx context.Context, messageID string) error
}
