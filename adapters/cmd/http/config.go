package http

import (
	adapters_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/adapters"
	app_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/app"
	device_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/devices"
	message_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/messages"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/middleware"
	"github.com/go-chi/chi/v5"
)

func Init() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.AuthAPIKey)

	app_handler.InitApp(r)
	device_handler.InitDevices(r)
	message_handler.InitMessages(r)
	adapters_handler.InitAdapters(r)

	return r
}
