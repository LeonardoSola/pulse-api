package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/utils"
	"github.com/LeonardoSola/wagon-api/internal/custom_error"
)

func AuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.ErrorResponse(custom_error.ErrUnauthorized, w)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.ErrorResponse(custom_error.ErrUnauthorized, w)
			return
		}

		token := parts[1]

		userID, err := bootstrap.Services.AuthService.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(err, w)
			return
		}

		// Set user id
		ctx := context.WithValue(r.Context(), authUserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
