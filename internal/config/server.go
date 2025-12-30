package config

import (
	"fmt"
	"os"
	"strconv"
)

type serverConfigStruct struct {
	Port   int
	APIKey string
}

func LoadServerEnvironment() (serverConfigStruct, error) {
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if port == 0 {
		port = 8080
	}

	key := os.Getenv("API_KEY")
	if key == "" {
		return serverConfigStruct{}, fmt.Errorf("API_KEY is required")
	}

	return serverConfigStruct{
		Port:   port,
		APIKey: key,
	}, nil
}
