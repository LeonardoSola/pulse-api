package rabbitmq

import (
	"fmt"

	"github.com/LeonardoSola/wagon-api/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Messaging *messagingAdapter
)

func InitRabbitMQ() error {
	if !config.Config.RabbitMQ.Enabled {
		return nil
	}

	connection, err := amqp.Dial(config.Config.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := connection.Channel()
	if err != nil {
		connection.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	Messaging = newMessaging(ch, config.Config.RabbitMQ.Exchange)

	return nil
}

func Close() {
	if Messaging.channel != nil {
		Messaging.channel.Close()
	}
}
