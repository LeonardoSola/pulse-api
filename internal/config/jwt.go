package config

import "os"

type jwtConfigStruct struct {
	Secret string
}

func LoadJWTEnvironment() (jwtConfigStruct, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	return jwtConfigStruct{
		Secret: secret,
	}, nil
}
