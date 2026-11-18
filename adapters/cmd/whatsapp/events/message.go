package events

import (
	"context"
	"log/slog"

	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
	"github.com/LeonardoSola/pulse-api/internal/ports"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

var messageService ports.MessageService

// SetMessageService sets the message service for event handlers
func SetMessageService(service ports.MessageService) {
	messageService = service
}

func MessageHandler(deviceJID string, evt *events.Message) {
	if messageService == nil {
		slog.Warn("MessageHandler: message service not initialized")
		return
	}

	message := message_entity.NewMessage(
		"",
		messageText(evt.Message),
		evt.Info.Chat.String(),
		message_entity.RecipientTypeUser,
		evt.Info.Sender.ToNonAD().String(),
	)
	message.DeviceJID = deviceJID
	message.ExternalID = evt.Info.ID
	if !evt.Info.Timestamp.IsZero() {
		message.CreatedAt = evt.Info.Timestamp
	}

	// Store the chat as bare ID + type (user or group), like outgoing messages.
	if err := message.SetRecipient(evt.Info.Chat.ToNonAD().String()); err != nil {
		slog.Warn("MessageHandler: unexpected chat JID", "chat", evt.Info.Chat.String())
	}

	if err := messageService.ReceiveMessage(context.Background(), message); err != nil {
		slog.Error("MessageHandler", "error", err)
	}
}

// messageText extracts the human-readable text of a message: plain text,
// extended text (replies, links) or the caption of an image/video/document.
func messageText(msg *waE2E.Message) string {
	switch {
	case msg.GetConversation() != "":
		return msg.GetConversation()
	case msg.GetExtendedTextMessage().GetText() != "":
		return msg.GetExtendedTextMessage().GetText()
	case msg.GetImageMessage().GetCaption() != "":
		return msg.GetImageMessage().GetCaption()
	case msg.GetVideoMessage().GetCaption() != "":
		return msg.GetVideoMessage().GetCaption()
	case msg.GetDocumentMessage().GetCaption() != "":
		return msg.GetDocumentMessage().GetCaption()
	default:
		return ""
	}
}
