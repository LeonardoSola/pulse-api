package config

import (
	"fmt"
	"os"
	"strconv"
)

type serverConfigStruct struct {
	Port string
}

func LoadServerEnvironment() (serverConfigStruct, error) {
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if port == 0 {
		port = 8080
	}

	return serverConfigStruct{
		Port: fmt.Sprintf(":%d", port),
	}, nil
}
