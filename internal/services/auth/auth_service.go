package auth_service

import (
	"crypto/subtle"

	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

// APIKeyService validates the static API key sent by clients.
type APIKeyService struct {
	apiKey []byte
}

func NewAPIKeyService(apiKey string) *APIKeyService {
	return &APIKeyService{apiKey: []byte(apiKey)}
}

// ValidateAPIKey compares the given key with the configured one in constant time.
func (s *APIKeyService) ValidateAPIKey(apiKey string) error {
	if len(s.apiKey) == 0 || subtle.ConstantTimeCompare([]byte(apiKey), s.apiKey) != 1 {
		return custom_error.ErrUnauthorized
	}

	return nil
}
