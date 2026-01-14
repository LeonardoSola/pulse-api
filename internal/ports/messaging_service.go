package ports

import (
	"context"
)

type MessagingService interface {
	Publish(ctx context.Context, queue string, message any) error
}
