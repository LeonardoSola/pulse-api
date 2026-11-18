package auth_service

import (
	"errors"
	"testing"

	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

func TestValidateAPIKey(t *testing.T) {
	svc := NewAPIKeyService("s3cret")

	if err := svc.ValidateAPIKey("s3cret"); err != nil {
		t.Errorf("valid key rejected: %v", err)
	}

	for _, key := range []string{"", "wrong", "s3cret ", "S3CRET"} {
		if err := svc.ValidateAPIKey(key); !errors.Is(err, custom_error.ErrUnauthorized) {
			t.Errorf("ValidateAPIKey(%q) = %v, want ErrUnauthorized", key, err)
		}
	}
}

func TestValidateAPIKeyRejectsEverythingWhenUnconfigured(t *testing.T) {
	if err := NewAPIKeyService("").ValidateAPIKey(""); err == nil {
		t.Error("an empty configured key must not accept an empty key")
	}
}
