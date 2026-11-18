package config

import (
	"net/url"
	"os"
	"strings"
)

type rabbitmqConfigStruct struct {
	URL      string
	Exchange string
	Enabled  bool
}

func LoadRabbitMQEnvironment() (rabbitmqConfigStruct, error) {
	enabled := strings.EqualFold(os.Getenv("RABBITMQ_ENABLED"), "true")

	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = buildRabbitMQURL(
			getEnv("RABBITMQ_HOST", "localhost"),
			getEnv("RABBITMQ_PORT", "5672"),
			getEnv("RABBITMQ_USER", "guest"),
			getEnv("RABBITMQ_PASSWORD", "guest"),
			getEnv("RABBITMQ_VHOST", "/"),
		)
	}

	return rabbitmqConfigStruct{
		URL:      amqpURL,
		Exchange: getEnv("RABBITMQ_EXCHANGE", "pulse-api"),
		Enabled:  enabled,
	}, nil
}

// buildRabbitMQURL builds an AMQP URL, escaping credentials and vhost.
func buildRabbitMQURL(host, port, user, password, vhost string) string {
	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
	}

	// "/" is the default vhost and maps to an empty path; any other vhost is
	// a single path segment and must be escaped (e.g. "/prod" -> "%2Fprod").
	if vhost != "/" && vhost != "" {
		u.RawPath = "/" + url.PathEscape(vhost)
		u.Path = "/" + vhost
	} else {
		u.Path = "/"
	}

	return u.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
