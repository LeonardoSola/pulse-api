package auth_service

import (
	"github.com/LeonardoSola/pulse-api/internal/config"
	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

type JWTService struct {
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) ValidateAPIKey(apiKey string) error {
	if apiKey != config.Config.Server.APIKey {
		return custom_error.ErrUnauthorized
	}

	return nil
}
