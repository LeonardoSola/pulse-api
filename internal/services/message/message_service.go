package message_service

import (
	"context"

	device_entity "github.com/LeonardoSola/wagon-api/internal/entities/device"
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/logger"
	"github.com/LeonardoSola/wagon-api/internal/ports"
)

type MessageService struct {
	messageRepository ports.MessageRepository
	deviceService     ports.DeviceRepository
	whatsappService   ports.WhatsappService
	messagingService  ports.MessagingService
}

func NewMessageService(messageRepository ports.MessageRepository, deviceService ports.DeviceRepository, whatsappService ports.WhatsappService, messagingService ports.MessagingService) *MessageService {
	return &MessageService{
		messageRepository,
		deviceService,
		whatsappService,
		messagingService,
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

	message.DeviceJID = *device.JID

	if err = message.Validate(); err != nil {
		return err
	}

	var messageID string
	if message.ImageURL != "" {
		messageID, err = s.whatsappService.SendImageMessage(ctx, message)
	} else {
		messageID, err = s.whatsappService.SendTextMessage(ctx, message)
	}
	if err != nil {
		return err
	}

	message.ExternalID = messageID

	err = s.messageRepository.Create(ctx, message)
	if err != nil {
		logger.Log.Error("MessageService.CreateMessage", "error", err)
		return err
	}

	return nil
}
