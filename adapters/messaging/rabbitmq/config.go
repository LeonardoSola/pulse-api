package rabbitmq

import (
	"fmt"

	"github.com/LeonardoSola/pulse-api/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Messaging  *messagingAdapter
	connection *amqp.Connection
)

// InitRabbitMQ connects to RabbitMQ and declares the configured topic
// exchange. It is a no-op when RabbitMQ is disabled.
func InitRabbitMQ() error {
	if !config.Config.RabbitMQ.Enabled {
		return nil
	}

	conn, err := amqp.Dial(config.Config.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	exchange := config.Config.RabbitMQ.Exchange
	if err := ch.ExchangeDeclare(exchange, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		conn.Close()
		return fmt.Errorf("failed to declare exchange %q: %w", exchange, err)
	}

	connection = conn
	Messaging = newMessaging(ch, exchange)

	return nil
}

// Close releases the channel and connection, if any.
func Close() {
	if Messaging != nil && Messaging.channel != nil {
		Messaging.channel.Close()
	}

	if connection != nil {
		connection.Close()
	}
}
