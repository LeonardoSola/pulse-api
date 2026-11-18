package middleware

import (
	"net/http"

	"github.com/LeonardoSola/pulse-api/adapters/bootstrap"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/utils"
	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

func AuthAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			utils.ErrorResponse(custom_error.ErrUnauthorized, w)
			return
		}

		err := bootstrap.Services.AuthService.ValidateAPIKey(apiKey)
		if err != nil {
			utils.ErrorResponse(err, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
