package events

import (
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/logger"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(evt *events.Message) {
	message := message_entity.Message{
		Text:      evt.Message.GetConversation(),
		From:      evt.Info.Sender.ToNonAD().String(),
		Recipient: evt.Info.Chat.String(),
	}

	logger.Log.Debug("MessageHandler", "from", message.From, "recipient", message.Recipient, "text", message.Text)

	// platform := whastapp.NewWhatsappPlatform()
	// platform.SetEvent(*evt)

	// handlers.MessageHandler(platform, message)
}
