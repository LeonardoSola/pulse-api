package bootstrap

import (
	"github.com/LeonardoSola/pulse-api/adapters/cmd/whatsapp"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/whatsapp/events"
	database "github.com/LeonardoSola/pulse-api/adapters/database/postgres"
	dummy_messaging "github.com/LeonardoSola/pulse-api/adapters/messaging/dummy"
	"github.com/LeonardoSola/pulse-api/adapters/messaging/rabbitmq"
	"github.com/LeonardoSola/pulse-api/internal/config"
	"github.com/LeonardoSola/pulse-api/internal/ports"
	auth_service "github.com/LeonardoSola/pulse-api/internal/services/auth"
	device_service "github.com/LeonardoSola/pulse-api/internal/services/device"
	message_service "github.com/LeonardoSola/pulse-api/internal/services/message"
)

type services struct {
	AuthService    *auth_service.JWTService
	DeviceService  *device_service.DeviceService
	MessageService *message_service.MessageService
}

var Services *services

func Init() {
	db := database.DB
	whats := whatsapp.NewWhatsappService()
	var messaging ports.MessagingService = rabbitmq.Messaging
	if !config.Config.RabbitMQ.Enabled {
		messaging = dummy_messaging.NewDummyMessaging()
	}

	Services = &services{
		AuthService:    auth_service.NewJWTService(),
		DeviceService:  device_service.NewDeviceService(db.DeviceRepository, whats),
		MessageService: message_service.NewMessageService(db.MessageRepository, db.DeviceRepository, whats, messaging),
	}

	// Set message service for event handlers (breaks import cycle)
	events.SetMessageService(Services.MessageService)
}
