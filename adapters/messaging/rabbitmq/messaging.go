package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/LeonardoSola/wagon-api/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type messagingAdapter struct {
	channel  *amqp.Channel
	exchange string
}

func newMessaging(ch *amqp.Channel, exchange string) *messagingAdapter {
	return &messagingAdapter{
		channel:  ch,
		exchange: exchange,
	}
}

// Publish publishes a message to a queue
func (m *messagingAdapter) Publish(ctx context.Context, queue string, message []byte) error {
	if m.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	// Create context with timeout
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Publish message
	err := m.channel.PublishWithContext(
		publishCtx,
		m.exchange, // exchange
		queue,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Log.Debug("Message published", "queue", queue, "size", len(message))
	return nil
}
