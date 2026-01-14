package config

import (
	"fmt"
	"os"
)

type configStruct struct {
	Env      string
	Server   serverConfigStruct
	Database databaseConfigStruct
	RabbitMQ rabbitmqConfigStruct
}

var Config *configStruct

func Load() error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	database, err := LoadDatabaseEnvironment()
	if err != nil {
		return fmt.Errorf("failed to load database environment: %w", err)
	}

	server, err := LoadServerEnvironment()
	if err != nil {
		return fmt.Errorf("failed to load server environment: %w", err)
	}

	rabbitmq, err := LoadRabbitMQEnvironment()
	if err != nil {
		return fmt.Errorf("failed to load rabbitmq environment: %w", err)
	}

	Config = &configStruct{
		Env:      env,
		Server:   server,
		Database: database,
		RabbitMQ: rabbitmq,
	}

	return nil
}
