package message

import (
	"context"

	"github.com/LeonardoSola/wagon-api/internal/ports"
)

type MessageService struct {
	messageRepository MessageRepository
	whatsappService   ports.WhatsappService
}

func NewMessageService(messageRepository MessageRepository, whatsappService ports.WhatsappService) *MessageService {
	return &MessageService{
		messageRepository,
		whatsappService,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, message Message) error {
	// err := s.messageRepository.Create(ctx, message)
	// if err != nil {
	// 	return err
	// }

	return s.whatsappService.SendTextMessage(ctx, message)
}
