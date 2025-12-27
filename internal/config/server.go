package config

import (
	"os"
	"strconv"
)

type serverConfigStruct struct {
	Port int
}

func LoadServerEnvironment() (serverConfigStruct, error) {
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if port == 0 {
		port = 8080
	}

	return serverConfigStruct{
		Port: port,
	}, nil
}
