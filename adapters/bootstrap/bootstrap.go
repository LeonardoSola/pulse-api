package bootstrap

import (
	"github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp"
	database "github.com/LeonardoSola/wagon-api/adapters/database/postgres"
	auth_service "github.com/LeonardoSola/wagon-api/internal/services/auth"
	device_service "github.com/LeonardoSola/wagon-api/internal/services/device"
	message_service "github.com/LeonardoSola/wagon-api/internal/services/message"
)

type services struct {
	AuthService    *auth_service.JWTService
	DeviceService  *device_service.DeviceService
	MessageService *message_service.MessageService
}

var Services *services

func Init() {
	db := database.DB
	whats := whastapp.NewWhatsappService()

	Services = &services{
		AuthService:    auth_service.NewJWTService(),
		DeviceService:  device_service.NewDeviceService(db.DeviceRepository, whats),
		MessageService: message_service.NewMessageService(db.MessageRepository, db.DeviceRepository, whats),
	}
}
