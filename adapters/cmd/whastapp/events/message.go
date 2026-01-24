package events

import (
	"context"

	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
	"github.com/LeonardoSola/pulse-api/internal/logger"
	"github.com/LeonardoSola/pulse-api/internal/ports"
	"go.mau.fi/whatsmeow/types/events"
)

var messageService ports.MessageService

// SetMessageService sets the message service for event handlers
func SetMessageService(service ports.MessageService) {
	messageService = service
}

func MessageHandler(deviceJID string, evt *events.Message) {
	if messageService == nil {
		logger.Log.Warn("MessageHandler: message service not initialized")
		return
	}

	message := message_entity.NewMessage(
		"",
		evt.Message.GetConversation(),
		evt.Info.Chat.String(),
		message_entity.RecipientTypeUser,
		evt.Info.Sender.ToNonAD().String(),
	)
	message.DeviceJID = deviceJID

	err := messageService.ReceiveMessage(context.Background(), message)
	if err != nil {
		logger.Log.Error("MessageHandler", "error", err)
		return
	}
}
