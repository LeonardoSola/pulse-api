package message_service

import (
	"context"

	device_entity "github.com/LeonardoSola/pulse-api/internal/entities/device"
	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
	"github.com/LeonardoSola/pulse-api/internal/logger"
	"github.com/LeonardoSola/pulse-api/internal/ports"
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
	message.From = *device.JID

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

func (s *MessageService) ReceiveMessage(ctx context.Context, message message_entity.Message) error {
	device, err := s.deviceService.GetByJID(ctx, message.DeviceJID)
	if err != nil {
		logger.Log.Error("MessageService.ReceiveMessage", "error", err)
		return err
	}

	if device.JID == nil || *device.JID == "" {
		return device_entity.ErrDeviceNotLoggedIn
	}

	message.DeviceID = device.ID

	err = s.messageRepository.Create(ctx, message)
	if err != nil {
		logger.Log.Error("MessageService.ReceiveMessage", "error", err)
		return err
	}

	err = s.messagingService.Publish(ctx, "messages.received", message)
	if err != nil {
		logger.Log.Error("MessageService.ReceiveMessage: failed to publish to queue", "error", err)
		return err
	}

	return nil
}
