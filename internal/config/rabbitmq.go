package config

import (
	"fmt"
	"os"
	"strings"
)

type rabbitmqConfigStruct struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
	//
	Exchange string
	Enabled  bool
}

func LoadRabbitMQEnvironment() (rabbitmqConfigStruct, error) {
	rabbitmqEnable := os.Getenv("RABBITMQ_ENABLED")
	rabbitmqEnable = strings.ToLower(rabbitmqEnable)
	enabled := false
	if rabbitmqEnable == "true" {
		enabled = true
	}

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		host := os.Getenv("RABBITMQ_HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("RABBITMQ_PORT")
		if port == "" {
			port = "5672"
		}

		user := os.Getenv("RABBITMQ_USER")
		if user == "" {
			user = "guest"
		}

		password := os.Getenv("RABBITMQ_PASSWORD")
		if password == "" {
			password = "guest"
		}

		vhost := os.Getenv("RABBITMQ_VHOST")
		if vhost == "" {
			vhost = "/"
		}

		url = fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, password, host, port, vhost)
	}

	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	if exchange == "" {
		exchange = "pulse-api"
	}

	return rabbitmqConfigStruct{
		URL:      url,
		Exchange: exchange,
		Enabled:  enabled,
	}, nil
}
