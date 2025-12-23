package config

import (
	"fmt"
	"os"
)

type databaseConfigStruct struct {
	ConnectionString string
	Driver           string
}

func LoadDatabaseEnvironment() (databaseConfigStruct, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL_MODE"),
		)
	}

	return databaseConfigStruct{
		ConnectionString: connStr,
		Driver:           os.Getenv("DB_DRIVER"),
	}, nil
}
