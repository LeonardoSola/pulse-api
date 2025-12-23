package http

import (
	app_handler "github.com/LeonardoSola/wagon-api/adapters/cmd/http/handlers/app"
	devices_handler "github.com/LeonardoSola/wagon-api/adapters/cmd/http/handlers/devices"
	"github.com/go-chi/chi/v5"
)

func Init() *chi.Mux {
	r := chi.NewRouter()

	app_handler.InitApp(r)
	devices_handler.InitDevices(r)

	return r
}
