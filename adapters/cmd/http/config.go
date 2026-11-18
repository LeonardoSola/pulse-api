package http

import (
	"github.com/LeonardoSola/pulse-api/adapters/bootstrap"
	adapters_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/adapters"
	app_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/app"
	device_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/devices"
	message_handler "github.com/LeonardoSola/pulse-api/adapters/cmd/http/handlers/messages"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/middleware"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chi_middleware.Recoverer)

	// Public routes
	app_handler.InitApp(r)

	// Routes protected by API key
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthAPIKey(bootstrap.Services.AuthService))

		device_handler.InitDevices(r)
		message_handler.InitMessages(r)
		adapters_handler.InitAdapters(r)
	})

	return r
}
