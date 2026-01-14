package dummy_messaging

import (
	"context"

	"github.com/LeonardoSola/wagon-api/internal/ports"
)

type dummyMessaging struct {
}

func NewDummyMessaging() ports.MessagingService {
	return &dummyMessaging{}
}

func (d *dummyMessaging) Publish(ctx context.Context, queue string, message []byte) error {
	return nil
}
