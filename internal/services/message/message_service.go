package message_service

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/ports"
)

type MessageService struct {
	messageRepository ports.MessageRepository
	whatsappService   ports.WhatsappService
}

func NewMessageService(messageRepository ports.MessageRepository, whatsappService ports.WhatsappService) *MessageService {
	return &MessageService{
		messageRepository,
		whatsappService,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, message message_entity.Message) error {
	// err := s.messageRepository.Create(ctx, message)
	// if err != nil {
	// 	return err
	// }

	return s.whatsappService.SendTextMessage(ctx, message)
}
