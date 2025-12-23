package bootstrap

import (
	"github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp"
	database "github.com/LeonardoSola/wagon-api/adapters/database/postgres"
	"github.com/LeonardoSola/wagon-api/domain/auth"
	"github.com/LeonardoSola/wagon-api/domain/device"
)

type services struct {
	AuthService   *auth.JWTService
	DeviceService *device.DeviceService
}

var Services *services

func Init() {
	db := database.DB
	whats := whastapp.NewWhatsappService()

	Services = &services{
		AuthService:   auth.NewJWTService(),
		DeviceService: device.NewDeviceService(db.DevicesRepository, whats),
	}
}
