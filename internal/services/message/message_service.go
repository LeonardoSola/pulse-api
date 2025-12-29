package message_service

import (
	"context"

	device_entity "github.com/LeonardoSola/wagon-api/internal/entities/device"
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/ports"
)

type MessageService struct {
	messageRepository ports.MessageRepository
	deviceService     ports.DeviceRepository
	whatsappService   ports.WhatsappService
}

func NewMessageService(messageRepository ports.MessageRepository, deviceService ports.DeviceRepository, whatsappService ports.WhatsappService) *MessageService {
	return &MessageService{
		messageRepository,
		deviceService,
		whatsappService,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, message message_entity.Message) error {
	device, err := s.deviceService.Get(ctx, message.DeviceID)
	if err != nil {
		return err
	}

	if device.JID == nil || *device.JID == "" {
		return device_entity.ErrDeviceNotLoggedIn
	}

	message.DeviceID = *device.JID

	if err = message.Validate(); err != nil {
		return err
	}

	// err := s.messageRepository.Create(ctx, message)
	// if err != nil {
	// 	return err
	// }

	messageID, err := s.whatsappService.SendTextMessage(ctx, message)
	if err != nil {
		return err
	}

	message.ID = messageID

	return s.messageRepository.Create(ctx, message)
}
