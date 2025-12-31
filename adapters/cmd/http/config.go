package http

import (
	"bytes"
	"io"
	"net/http"

	adapters_handler "github.com/LeonardoSola/wagon-api/adapters/cmd/http/handlers/adapters"
	app_handler "github.com/LeonardoSola/wagon-api/adapters/cmd/http/handlers/app"
	device_handler "github.com/LeonardoSola/wagon-api/adapters/cmd/http/handlers/devices"
	message_handler "github.com/LeonardoSola/wagon-api/adapters/cmd/http/handlers/messages"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/middleware"
	"github.com/LeonardoSola/wagon-api/internal/logger"
	"github.com/go-chi/chi/v5"
)

func Init() *chi.Mux {
	r := chi.NewRouter()

	// r.Use(logRequest)

	r.Use(middleware.AuthAPIKey)

	app_handler.InitApp(r)
	device_handler.InitDevices(r)
	message_handler.InitMessages(r)
	adapters_handler.InitAdapters(r)

	return r
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("Request", "method", r.Method, "path", r.URL.Path)
		logger.Log.Debug("Request", "headers", r.Header)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Error("Error reading request body", "error", err)
		}
		bodyToBytes := bytes.NewBuffer(body).String()
		logger.Log.Debug("Request", "body", bodyToBytes)

		logger.Log.Debug("Request", "params", r.URL.Query())
		next.ServeHTTP(w, r)
	})
}
