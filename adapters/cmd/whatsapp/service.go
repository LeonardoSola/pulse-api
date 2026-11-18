package whatsapp

import (
	"github.com/LeonardoSola/pulse-api/internal/ports"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type WhatsappService struct {
	container *sqlstore.Container
}

func NewWhatsappService() ports.WhatsappService {
	return &WhatsappService{
		container: container,
	}
}
