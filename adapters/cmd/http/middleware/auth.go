package middleware

import (
	"net/http"

	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/utils"
	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

// APIKeyHeader is the request header that carries the API key.
const APIKeyHeader = "X-API-KEY"

type apiKeyValidator interface {
	ValidateAPIKey(apiKey string) error
}

// AuthAPIKey rejects requests that do not carry a valid API key.
func AuthAPIKey(validator apiKeyValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(APIKeyHeader)
			if apiKey == "" {
				utils.ErrorResponse(custom_error.ErrUnauthorized, w)
				return
			}

			if err := validator.ValidateAPIKey(apiKey); err != nil {
				utils.ErrorResponse(err, w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
