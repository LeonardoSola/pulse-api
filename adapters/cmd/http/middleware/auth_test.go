package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	auth_service "github.com/LeonardoSola/pulse-api/internal/services/auth"
)

func TestAuthAPIKey(t *testing.T) {
	handler := AuthAPIKey(auth_service.NewAPIKeyService("s3cret"))(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	tests := []struct {
		name   string
		key    string
		status int
	}{
		{name: "missing key", key: "", status: http.StatusUnauthorized},
		{name: "wrong key", key: "nope", status: http.StatusUnauthorized},
		{name: "valid key", key: "s3cret", status: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/devices", nil)
			if tt.key != "" {
				req.Header.Set(APIKeyHeader, tt.key)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("status = %d, want %d", rec.Code, tt.status)
			}
		})
	}
}
