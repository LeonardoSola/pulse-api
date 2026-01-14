package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

func (m *messagingAdapter) Publish(ctx context.Context, queue string, message any) error {
	if m.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = m.channel.PublishWithContext(
		publishCtx,
		m.exchange, // exchange
		queue,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         messageJSON,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
